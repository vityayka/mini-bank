package async

import (
	db "bank/db/sqlc"
	"bank/utils"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	taskNameSendVerifyEmail = "task:send_verify_email"
	codeExpiresIn           = 15 * time.Minute
)

type PayloadSendVerifyEmail struct {
	UserID int64 `json:"user_id"`
}

// DistributeTaskVerifyEmail implements TaskDistributor.
func (r *RedisTaskDistributor) DistributeTaskVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail, opt ...asynq.Option) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("couldn't marshal task payload: %w", err)
	}

	task := asynq.NewTask(taskNameSendVerifyEmail, payloadBytes, opt...)
	info, err := r.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("couldn't enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Dur("timeout", info.Retention).Str("queue", info.Queue).Bytes("payload", task.Payload()).
		Int("max_retry", info.MaxRetry).Msg("enqueued task")

	return nil
}

func (r *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	user, err := r.store.GetUser(ctx, payload.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user not found: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("store.GetUser err: %w", err)
	}

	verifyEmail, err := r.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		UserID:    user.ID,
		Email:     user.Email,
		Code:      utils.RandomString(32),
		IsUsed:    false,
		ExpiredAt: time.Now().Add(codeExpiresIn),
	})

	if err != nil {
		return fmt.Errorf("failed to create verify_email: %w", err)
	}

	link := fmt.Sprintf("http://localhost:8080/v1/verify_email?id=%d&code=%s", verifyEmail.ID, verifyEmail.Code)
	err = r.mailSender.Send(
		"Verify email",
		fmt.Sprintf("Thanks for signing up! Please follow the <a href=\"%s\">link</a> to verify", link),
		[]string{user.Email}, nil, nil, nil,
	)

	if err != nil {
		return fmt.Errorf("failed to send a verification email to %s: %w", user.Email, err)
	}

	log.Info().Str("type", task.Type()).Str("email", user.Email).Bytes("payload", task.Payload()).
		Msg("processed task")

	return nil
}
