package async

import (
	db "bank/db/sqlc"
	"bank/mail"
	"context"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(context.Context, *asynq.Task) error
}

type RedisTaskProcessor struct {
	server     *asynq.Server
	store      db.Store
	mailSender mail.EmailSender
}

func (r *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(taskNameSendVerifyEmail, r.ProcessTaskSendVerifyEmail)

	return r.server.Start(mux)
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, mailSender mail.EmailSender) TaskProcessor {
	return &RedisTaskProcessor{
		server: asynq.NewServer(redisOpt, asynq.Config{
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Err(err).
					Str("task_type", task.Type()).
					Bytes("payload", task.Payload()).
					Msg("error_processing_task")
			}),
			Logger: &Logger{},
		}),
		store:      store,
		mailSender: mailSender,
	}
}
