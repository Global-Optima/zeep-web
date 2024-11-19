package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func RedisMiddleware(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("redis", rdb)
		c.Next()
	}
}
