package orders

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/websockets"
)

type OrdersNotifier struct {
	hub *websockets.WebSocketHub
}

func NewOrdersNotifier(hub *websockets.WebSocketHub) *OrdersNotifier {
	return &OrdersNotifier{
		hub: hub,
	}
}

func (n *OrdersNotifier) NotifyNewOrder(orderID uint, data interface{}) {
	channel := "orders"
	event := "order_created"
	message := fmt.Sprintf("New order with ID %d has been created", orderID)
	n.hub.Broadcast(channel, event, map[string]interface{}{
		"message": message,
		"data":    data,
	})
}

func (n *OrdersNotifier) NotifySubOrderCompleted(orderID uint, subOrderID uint, data interface{}) {
	channel := "orders"
	event := "suborder_completed"
	message := fmt.Sprintf("Sub-order %d of order %d has been completed", subOrderID, orderID)
	n.hub.Broadcast(channel, event, map[string]interface{}{
		"message": message,
		"data":    data,
	})
}

func (n *OrdersNotifier) NotifyOrderCompleted(orderID uint, data interface{}) {
	channel := "orders"
	event := "order_completed"
	message := fmt.Sprintf("Order %d has been completed", orderID)
	n.hub.Broadcast(channel, event, map[string]interface{}{
		"message": message,
		"data":    data,
	})
}
