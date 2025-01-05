package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ERROR_MESSAGE_BINDING_JSON  = "invalid input: failed to bind json body"
	ERROR_MESSAGE_BINDING_QUERY = "failed to bind query parameters"
)

func SendSuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func SendResponseWithStatus(c *gin.Context, data interface{}, status int) {
	c.JSON(status, data)
}

func SendSuccessResponseWithPagination(c *gin.Context, data interface{}, pagination *Pagination) {
	c.JSON(http.StatusOK, gin.H{
		"data":       data,
		"pagination": pagination,
	})
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

func SendInternalServerError(c *gin.Context, message string) {
	SendErrorWithStatus(c, message, http.StatusInternalServerError)
}

func SendBadRequestError(c *gin.Context, message string) {
	SendErrorWithStatus(c, message, http.StatusBadRequest)
}

func SendNotFoundError(c *gin.Context, message string) {
	SendErrorWithStatus(c, message, http.StatusNotFound)
}

func SuccessCreatedResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, data)
}
