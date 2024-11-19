package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	client *redis.Client
	ctx    context.Context
}

var (
	cacheInstance *Cache
	once          sync.Once
)

func InitCache(redisClient *redis.Client, ctx context.Context) {
	once.Do(func() {
		cacheInstance = &Cache{
			client: redisClient,
			ctx:    ctx,
		}
	})
}

func GetCacheInstance() *Cache {
	if cacheInstance == nil {
		panic("CacheUtil not initialized. Call InitCache before using it.")
	}
	return cacheInstance
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to serialize value: %w", err)
	}

	if ttl > 0 {
		err = c.client.Set(c.ctx, key, data, ttl).Err()
	} else {
		err = c.client.Set(c.ctx, key, data, 0).Err()
	}

	if err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}
	return nil
}

func (c *Cache) Get(key string, dest interface{}) error {
	data, err := c.client.Get(c.ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("key %s does not exist", key)
	} else if err != nil {
		return fmt.Errorf("failed to get key %s: %w", key, err)
	}

	if err := json.Unmarshal([]byte(data), dest); err != nil {
		return fmt.Errorf("failed to deserialize value for key %s: %w", key, err)
	}
	return nil
}

func (c *Cache) Delete(key string) error {
	err := c.client.Del(c.ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}
	return nil
}

func BuildCacheKey(module string, filters map[string]string) string {
	key := module
	for k, v := range filters {
		key = fmt.Sprintf("%s:%s:%s", key, k, v)
	}
	return key
}
