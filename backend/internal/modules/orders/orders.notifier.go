package orders

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders/types"
	"github.com/Global-Optima/zeep-web/backend/internal/websockets"
)

type OrdersNotifier struct {
	hub *websockets.WebSocketHub
}

type EventName string

const (
	EventOrderCreated      EventName = "ORDER_CREATED"
	EventSubOrderCompleted EventName = "SUBORDER_COMPLETED"
	EventOrderCompleted    EventName = "ORDER_COMPLETED"
)

func NewOrdersNotifier(hub *websockets.WebSocketHub) *OrdersNotifier {
	return &OrdersNotifier{
		hub: hub,
	}
}

func (n *OrdersNotifier) NotifyNewOrder(orderID uint, storeID uint, data types.OrderEvent) {
	channel := websockets.GetStoreChannel(storeID)
	event := EventOrderCreated
	message := fmt.Sprintf("New order with ID %d has been created", orderID)
	n.hub.Broadcast(channel, string(event), map[string]interface{}{
		"message": message,
		"data":    data,
	})
}

func (n *OrdersNotifier) NotifySubOrderCompleted(orderID uint, subOrderID uint, storeID uint, data *types.OrderEvent) {
	channel := websockets.GetStoreChannel(storeID)
	event := EventSubOrderCompleted
	message := fmt.Sprintf("Sub-order %d of order %d has been completed", subOrderID, orderID)
	n.hub.Broadcast(channel, string(event), map[string]interface{}{
		"message": message,
		"data":    data,
	})
}

func (n *OrdersNotifier) NotifyOrderCompleted(orderID uint, storeID uint, data *types.OrderEvent) {
	channel := websockets.GetStoreChannel(storeID)
	event := EventOrderCompleted
	message := fmt.Sprintf("Order %d has been completed", orderID)
	n.hub.Broadcast(channel, string(event), map[string]interface{}{
		"message": message,
		"data":    data,
	})
}
