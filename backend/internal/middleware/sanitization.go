package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Global-Optima/zeep-web/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

func SanitizeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.ContentType()

		if contentType == "multipart/form-data" {
			c.Next()
			return
		}

		if contentType != "application/json" {
			c.Next()
			return
		}

		var requestData map[string]interface{}
		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		for key, value := range requestData {
			if strVal, ok := value.(string); ok {
				sanitized, valid := utils.SanitizeString(strVal)
				if !valid {
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
						"error": "Invalid input",
						"field": key,
					})
					return
				}
				requestData[key] = sanitized
			}
		}

		sanitizedBody, _ := json.Marshal(requestData)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(sanitizedBody))

		c.Next()
	}
}
