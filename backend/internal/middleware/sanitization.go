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
		if c.Request.Method != http.MethodPost && c.Request.Method != http.MethodPut && c.Request.Method != http.MethodPatch {
			c.Next()
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		var requestData map[string]interface{}
		if err := json.Unmarshal(body, &requestData); err != nil {
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

		sanitizedBody, err := json.Marshal(requestData)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(sanitizedBody))

		c.Next()
	}
}
