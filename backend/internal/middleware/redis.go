package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RedisMiddleware(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("redis", rdb)
		c.Next()
	}
}
