package async

import (
	db "bank/db/sqlc"
	"context"

	"github.com/hibiken/asynq"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(context.Context, *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func (r *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(taskNameSendVerifyEmail, r.ProcessTaskSendVerifyEmail)

	return r.server.Start(mux)
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProcessor {
	return &RedisTaskProcessor{
		server: asynq.NewServer(redisOpt, asynq.Config{}),
		store:  store,
	}
}
