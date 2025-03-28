package websockets

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins; restrict in production.
	},
}

func WebSocketHandler(hub *WebSocketHub, channel string) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade to WebSocket"})
			return
		}

		sub := Subscription{
			Conn:    conn,
			Channel: channel,
		}
		hub.register <- sub

		defer func() {
			hub.unregister <- sub
			_ = conn.Close()
		}()

		for {
			var msg BroadcastMessage
			err := conn.ReadJSON(&msg)
			if err != nil {
				break
			}
		}
	}
}
