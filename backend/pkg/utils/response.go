package utils

import (
	"net/http"

	"github.com/Global-Optima/zeep-web/backend/internal/errors/handlerErrors"
	"github.com/gin-gonic/gin"
)

const (
	ERROR_MESSAGE_BINDING_JSON  = "Неверный ввод: ошибка привязки JSON данных"
	ERROR_MESSAGE_BINDING_QUERY = "Ошибка привязки параметров запроса"
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

// TODO create func SendCustomResponseWithStatus(c *gin.Context, localizedMessages *localization.LocalizedMessage, status int)

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

func SendSuccessCreatedResponse(c *gin.Context, message string) {
	SendMessageWithStatus(c, message, http.StatusCreated)
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
