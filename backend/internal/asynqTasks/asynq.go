package asynqTasks

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"sync"
	"time"
)

const (
	ASYNQ_CONCURRENCY = 10
	ASYNQ_RETRY_DELAY = 5 * time.Second
)

var once sync.Once

type AsynqManager struct {
	client    *asynq.Client
	server    *asynq.Server
	inspector *asynq.Inspector
	mux       *asynq.ServeMux
	logger    *zap.SugaredLogger
}

type AsynqManagerTask struct {
}

type MyRedisConnOpt struct {
	Rdb *redis.Client
}

func (opt *MyRedisConnOpt) MakeRedisClient() interface{} {
	return opt.Rdb
}

func initAsynq(redisClient *redis.Client, logger *zap.SugaredLogger) (*AsynqManager, error) {
	var manger *AsynqManager
	var initErr error

	once.Do(func() {
		redisConn := &MyRedisConnOpt{Rdb: redisClient}

		server := asynq.NewServer(
			redisConn,
			asynq.Config{
				Concurrency: ASYNQ_CONCURRENCY,
				BaseContext: func() context.Context {
					return context.WithValue(context.Background(), "service", "asynq-task-queue")
				},
				RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
					return ASYNQ_RETRY_DELAY
				},
				ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
					logger.Errorf("❌ Task %s failed: %v", task.Type(), err)
				}),
				Logger: logger,
			},
		)

		manger = &AsynqManager{
			client:    asynq.NewClient(redisConn),
			server:    server,
			inspector: asynq.NewInspector(redisConn),
			mux:       asynq.NewServeMux(),
			logger:    logger,
		}
		initErr = nil
	})

	return manger, initErr
}

func (m *AsynqManager) startAsynqServer() {
	go func() {
		if err := m.server.Run(m.mux); err != nil {
			m.logger.Fatalf("🔥 Failed to start Asynq server: %v", err)
		}
	}()

	m.logger.Info("✅ Asynq Server Started Successfully")
}

func NewAsyncManager(redisClient *redis.Client, logger *zap.SugaredLogger) (*AsynqManager, error) {
	manager, err := initAsynq(redisClient, logger)
	if err != nil {
		return nil, err
	}

	manager.startAsynqServer()

	return manager, nil
}

func (m *AsynqManager) RegisterTask(pattern string, task func(context.Context, *asynq.Task) error) {
	m.mux.HandleFunc(pattern, task)
}

func (m *AsynqManager) EnqueueTask(taskType string, payload []byte, delay time.Duration) error {
	_, err := m.client.Enqueue(asynq.NewTask(taskType, payload), asynq.ProcessIn(delay))
	if err != nil {
		m.logger.Errorf("❌ Failed to enqueue task %s: %v", taskType, err)
		return err
	}

	m.logger.Infof("✅ Task %s enqueued successfully with delay %v", taskType, delay)
	return nil
}
