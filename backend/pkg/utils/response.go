package utils

import (
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"net/http"

	"github.com/Global-Optima/zeep-web/backend/internal/errors/handlerErrors"
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

// TODO create single func SendResponseWithStatus()
func SendMessageWithStatus(c *gin.Context, componentName string, status int) {
	c.JSON(status, gin.H{
		"message": localization.TranslateResponse(status, componentName),
	})
}

func SendErrorWithStatus(c *gin.Context, componentName string, status int) {
	c.JSON(status, gin.H{
		"error": localization.TranslateResponse(status, componentName),
	})
}

func SendInternalServerError(c *gin.Context, componentName string) {
	SendErrorWithStatus(c, componentName, http.StatusInternalServerError)
}

func SendBadRequestError(c *gin.Context, componentName string) {
	SendErrorWithStatus(c, componentName, http.StatusBadRequest)
}

func SendNotFoundError(c *gin.Context, componentName string) {
	SendErrorWithStatus(c, componentName, http.StatusNotFound)
}

func SendSuccessCreatedResponse(c *gin.Context, componentName string) {
	SendMessageWithStatus(c, componentName, http.StatusCreated)
}

func SendDetailedError(c *gin.Context, err error, details ...string) {
	if handlerErr, ok := err.(handlerErrors.HandlerErrorInterface); ok {
		extendedError := handlerErr
		if len(details) > 0 {
			extendedError = handlerErr.WithDetails(details...)
		}
		response := gin.H{
			"error": extendedError.Error(),
		}
		if len(extendedError.Details()) > 0 {
			response["details"] = extendedError.Details()
		}
		c.JSON(handlerErr.Status(), response)
	} else {
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
	}
}
