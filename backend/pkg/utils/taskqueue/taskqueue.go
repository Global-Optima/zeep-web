package taskqueue

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
)

type TaskQueue interface {
	RegisterTask(pattern string, task func(context.Context, *asynq.Task) error)
	EnqueueTask(taskType string, payload []byte, delay time.Duration) error
}
