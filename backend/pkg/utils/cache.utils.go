package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
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

func BuildCacheKeyFromFilter(module string, filter FilterProvider) string {
	if module == "" || filter == nil {
		return ""
	}

	filters := make(map[string]string)

	// Extract pagination fields
	if pagination := filter.GetPagination(); pagination != nil {
		filters["page"] = strconv.Itoa(pagination.Page)
		filters["pageSize"] = strconv.Itoa(pagination.PageSize)
		filters["totalCount"] = strconv.Itoa(pagination.TotalCount)
		filters["totalPages"] = strconv.Itoa(pagination.TotalPages)
	}

	// Extract sort fields
	if sort := filter.GetSort(); sort != nil {
		filters["sortField"] = sort.Field
		filters["sortDirection"] = sort.Direction
	}

	// Extract other fields using reflection
	extractFields(filters, reflect.ValueOf(filter))

	// Use BuildCacheKey to generate the key
	return BuildCacheKey(module, filters)
}

func extractFields(filters map[string]string, value reflect.Value) {
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return // Ignore nil pointers
		}
		value = value.Elem() // Dereference pointer
	}

	if value.Kind() != reflect.Struct {
		return // Only process structs
	}

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := value.Type().Field(i)
		tag := fieldType.Tag.Get("form")

		// Skip untagged or empty fields
		if tag == "" || !field.IsValid() {
			continue
		}

		switch field.Kind() {
		case reflect.Ptr:
			if !field.IsNil() {
				filters[tag] = fmt.Sprintf("%v", field.Elem())
			}
		case reflect.Struct:
			// Handle nested structs
			extractFields(filters, field)
		case reflect.Slice, reflect.Array:
			// Handle slices/arrays by joining elements
			var elements []string
			for j := 0; j < field.Len(); j++ {
				elements = append(elements, fmt.Sprintf("%v", field.Index(j).Interface()))
			}
			filters[tag] = strings.Join(elements, ",")
		case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64, reflect.String:
			filters[tag] = fmt.Sprintf("%v", field.Interface())
		default:
			if fieldType.Type.Kind() == reflect.String {
				filters[tag] = fmt.Sprintf("%v", field.Interface())
			} else {
				filters[tag] = fmt.Sprintf("%v", field.Interface())
			}
		}
	}
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
