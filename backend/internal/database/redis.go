package database

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

var redisInstance *RedisClient

func InitRedis(host string, port int, password string, db int) (*RedisClient, error) {
	if redisInstance != nil {
		return redisInstance, nil
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})

	ctx := context.Background()
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	redisInstance = &RedisClient{
		Client: rdb,
		Ctx:    ctx,
	}

	return redisInstance, nil
}

func GetRedisClient() (*RedisClient, error) {
	if redisInstance == nil {
		return nil, fmt.Errorf("Redis client is not initialized")
	}
	return redisInstance, nil
}
