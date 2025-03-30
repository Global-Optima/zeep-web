package websockets

import (
	"sync"

	"github.com/gorilla/websocket"
)

var (
	hubInstance *WebSocketHub
	once        sync.Once
)

func GetHubInstance() *WebSocketHub {
	once.Do(func() {
		hubInstance = NewWebSocketHub()
		go hubInstance.Run() // Start the hub goroutine
	})
	return hubInstance
}

type WebSocketHub struct {
	channels   map[string]map[*websocket.Conn]bool // Channel -> Connections
	broadcast  chan BroadcastMessage               // Messages to be broadcast
	register   chan Subscription                   // New connections
	unregister chan Subscription                   // Disconnections
	mutex      sync.RWMutex                        // For thread-safe channel management
}

type Subscription struct {
	Conn    *websocket.Conn
	Channel string
}

type BroadcastMessage struct {
	Channel string      `json:"channel"`
	Event   string      `json:"event"`
	Data    interface{} `json:"data"`
}

func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		channels:   make(map[string]map[*websocket.Conn]bool),
		broadcast:  make(chan BroadcastMessage),
		register:   make(chan Subscription),
		unregister: make(chan Subscription),
	}
}

func (hub *WebSocketHub) Run() {
	for {
		select {
		case sub := <-hub.register:
			hub.mutex.Lock()
			if hub.channels[sub.Channel] == nil {
				hub.channels[sub.Channel] = make(map[*websocket.Conn]bool)
			}
			hub.channels[sub.Channel][sub.Conn] = true
			hub.mutex.Unlock()

		case sub := <-hub.unregister:
			hub.mutex.Lock()
			if connections, ok := hub.channels[sub.Channel]; ok {
				if _, exists := connections[sub.Conn]; exists {
					delete(connections, sub.Conn)
					if len(connections) == 0 {
						delete(hub.channels, sub.Channel)
					}
				}
			}
			hub.mutex.Unlock()

		case message := <-hub.broadcast:
			hub.mutex.RLock()
			if connections, ok := hub.channels[message.Channel]; ok {
				for conn := range connections {
					err := conn.WriteJSON(message)
					if err != nil {
						_ = conn.Close()
						delete(connections, conn)
					}
				}
			}
			hub.mutex.RUnlock()
		}
	}
}

func (hub *WebSocketHub) Broadcast(channel, event string, data interface{}) {
	hub.broadcast <- BroadcastMessage{
		Channel: channel,
		Event:   event,
		Data:    data,
	}
}
