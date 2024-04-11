package async

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistrubutor interface {
	DistributeTaskVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail, opt ...asynq.Option) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistrubutor {
	return &RedisTaskDistributor{
		client: asynq.NewClient(redisOpt),
	}
}
