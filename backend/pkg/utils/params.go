package utils

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseParam(c *gin.Context, paramName string) (uint, error) {
	idParam := c.Param(paramName)
	if idParam == "" {
		return 0, fmt.Errorf("parameter %s is required", paramName)
	}

	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid %s: must be a positive integer", paramName)
	}

	return uint(id), nil
}
