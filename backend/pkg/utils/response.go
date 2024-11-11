package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func SuccessPaginatedResponse(c *gin.Context, data interface{}, meta interface{}) {
	if meta != nil {
		c.JSON(http.StatusOK, gin.H{
			"data": data,
			"meta": meta,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	}
}

func SendMessageWithStatus(c *gin.Context, message string, status int) {
	c.JSON(status, gin.H{
		"message": message,
	})
}

func SendErrorWithStatus(c *gin.Context, message string, status int) {
	c.JSON(status, gin.H{
		"error": message,
	})
}

func SendInternalError(c *gin.Context, message string) {
	SendErrorWithStatus(c, message, http.StatusInternalServerError)
}

func SendBadRequestError(c *gin.Context, message string) {
	SendErrorWithStatus(c, message, http.StatusBadRequest)
}

func SendNotFoundError(c *gin.Context, message string) {
	SendErrorWithStatus(c, message, http.StatusNotFound)
}
