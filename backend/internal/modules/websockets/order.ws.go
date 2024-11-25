package websockets

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type OrderEvent struct {
	Type    string      `json:"type"` // e.g., "order_added", "order_deleted", "status_changed", "current_orders"
	Payload interface{} `json:"payload"`
}

type WebSocketManager struct {
	clients   map[*websocket.Conn]bool
	broadcast chan OrderEvent
	lock      sync.Mutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var Manager = NewWebSocketManager()

func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan OrderEvent),
	}
}

func (wm *WebSocketManager) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}

	wm.lock.Lock()
	wm.clients[conn] = true
	wm.lock.Unlock()

	log.Println("New WebSocket connection established.")

	defer func() {
		wm.lock.Lock()
		delete(wm.clients, conn)
		wm.lock.Unlock()
		conn.Close()
		log.Println("WebSocket connection closed.")
	}()
}

func (wm *WebSocketManager) HandleBroadcasts() {
	for {
		event := <-wm.broadcast

		wm.lock.Lock()
		for client := range wm.clients {
			err := client.WriteJSON(event)
			if err != nil {
				log.Println("WebSocket broadcast error:", err)
				client.Close()
				delete(wm.clients, client)
			}
		}
		wm.lock.Unlock()
	}
}

func (wm *WebSocketManager) SendEvent(event OrderEvent) {
	wm.broadcast <- event
}
