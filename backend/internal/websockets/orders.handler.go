package websockets

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetStoreChannel(storeID uint) string {
	return fmt.Sprintf("orders:%d", storeID)
}

func RegisterOrderWebsocketRoutes(router *gin.RouterGroup, handler gin.HandlerFunc) {
	router.GET("/orders/:storeId", handler)
}

func OrdersWebSocketHandler(hub *WebSocketHub) gin.HandlerFunc {
	return func(c *gin.Context) {
		storeIDStr := c.Param("storeId")
		if storeIDStr == "" {
			storeIDStr = c.Query("storeId")
		}
		storeID, err := strconv.ParseUint(storeIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing storeId"})
			return
		}

		channel := GetStoreChannel(uint(storeID))
		WebSocketHandler(hub, channel)(c)
	}
}
