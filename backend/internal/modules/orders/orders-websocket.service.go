package orders

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders/types"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type EventType string

const (
	EventTypeInitialData    EventType = "initial_data"
	EventTypeOrderSucceeded EventType = "order_succeeded"
	EventTypeOrderUpdated   EventType = "order_updated"
	EventTypeOrderDeleted   EventType = "order_deleted"
)

type WebSocketMessage struct {
	Type    EventType   `json:"type"`
	Payload interface{} `json:"payload"`
}

type Client struct {
	Conn    *websocket.Conn
	StoreID uint
}

type Hub struct {
	mu          sync.RWMutex
	connections map[uint]map[*Client]bool // StoreID -> Set of connected clients
}

var (
	hubInstance *Hub
	once        sync.Once
	upgrader    = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Add proper origin checks for production
			return true
		},
	}
)

// GetHubInstance ensures there is only one Hub instance.
func GetHubInstance() *Hub {
	once.Do(func() {
		hubInstance = &Hub{
			connections: make(map[uint]map[*Client]bool),
		}
	})
	return hubInstance
}

// AddClient adds a WebSocket client to the hub.
func (h *Hub) AddClient(storeID uint, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.connections[storeID] == nil {
		h.connections[storeID] = make(map[*Client]bool)
	}
	h.connections[storeID][client] = true
	log.Printf("Client connected for store %d. Total clients: %d", storeID, len(h.connections[storeID]))
}

// RemoveClient removes a WebSocket client from the hub.
func (h *Hub) RemoveClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	clients, exists := h.connections[client.StoreID]
	if exists {
		delete(clients, client)
		if len(clients) == 0 {
			delete(h.connections, client.StoreID)
		}
		log.Printf("Client disconnected from store %d. Total clients: %d", client.StoreID, len(h.connections[client.StoreID]))
	}

	defer func() {
		_ = client.Conn.Close()
	}()
}

// BroadcastMessage broadcasts a message to all clients connected to a specific store.
func (h *Hub) BroadcastMessage(storeID uint, eventType EventType, payload interface{}) error {
	message := WebSocketMessage{
		Type:    eventType,
		Payload: payload,
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal WebSocket message: %w", err)
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, exists := h.connections[storeID]
	if !exists {
		return nil // No clients connected to the store
	}

	for client := range clients {
		err := client.Conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Printf("Error broadcasting message to store %d: %v", storeID, err)
			h.RemoveClient(client)
		}
	}

	return nil
}

// BroadcastOrderSucceeded broadcasts an order creation event.
func BroadcastOrderSucceeded(storeID uint, order types.OrderDTO) {
	_ = GetHubInstance().BroadcastMessage(storeID, EventTypeOrderSucceeded, order)
}

// BroadcastOrderUpdated broadcasts an order update event.
func BroadcastOrderUpdated(storeID uint, order types.OrderDTO) {
	_ = GetHubInstance().BroadcastMessage(storeID, EventTypeOrderUpdated, order)
}

// BroadcastOrderDeleted broadcasts an order deletion event.
func BroadcastOrderDeleted(storeID uint, orderID uint) {
	_ = GetHubInstance().BroadcastMessage(storeID, EventTypeOrderDeleted, map[string]uint{"orderId": orderID})
}

// HandleClient initializes and manages a WebSocket client connection.
func HandleClient(storeID uint, conn *websocket.Conn, initialData []types.OrderDTO) {
	hub := GetHubInstance()
	client := &Client{
		Conn:    conn,
		StoreID: storeID,
	}

	// Add the client to the hub
	hub.AddClient(storeID, client)
	defer hub.RemoveClient(client)

	// Send initial data
	initialMessage := WebSocketMessage{
		Type:    EventTypeInitialData,
		Payload: initialData,
	}
	data, _ := json.Marshal(initialMessage)
	_ = client.Conn.WriteMessage(websocket.TextMessage, data)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading WebSocket message for store %d: %v", storeID, err)
			break
		}
	}
}

// UpgradeConnection upgrades an HTTP connection to a WebSocket.
func UpgradeConnection(c *gin.Context) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return nil, err
	}
	return conn, nil
}
