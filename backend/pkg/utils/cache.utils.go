package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
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
	if IsEmpty(value) {
		return fmt.Errorf("value for key %s is empty and cannot be cached", key)
	}

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to serialize value: %w", err)
	}

	err = c.client.Set(c.ctx, key, data, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}
	return nil
}

func (c *Cache) Get(key string, dest interface{}) error {
	data, err := c.client.Get(c.ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to get key %s: %w", key, err)
	}

	if len(data) == 0 {
		return fmt.Errorf("empty data")
	}

	if err := json.Unmarshal([]byte(data), dest); err != nil {
		return fmt.Errorf("failed to deserialize value for key %s: %w", key, err)
	}

	return nil
}

// Delete a key from the cache
func (c *Cache) Delete(key string) error {
	err := c.client.Del(c.ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}
	return nil
}

// Build a cache key dynamically based on module and filters
func BuildCacheKey(module string, filters map[string]string) string {
	if module == "" {
		return "" // Return an empty string if the module is empty
	}

	key := module
	for k, v := range filters {
		key = fmt.Sprintf("%s:%s:%s", key, k, v)
	}
	return key
}

func IsEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	default:
		return false
	}
}
