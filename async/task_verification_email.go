package async

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const taskNameSendVerifyEmail = "task:send_verify_email"

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
		return fmt.Errorf("stire.GetUser err: %w", err)
	}

	log.Info().Str("type", task.Type()).Str("email", user.Email).Bytes("payload", task.Payload()).
		Msg("processed task")

	return nil
}
