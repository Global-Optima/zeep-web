package types

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusCompleted OrderStatus = "COMPLETED"
	OrderStatusConfirmed OrderStatus = "CONFIRMED"
	OrderStatusShipped   OrderStatus = "SHIPPED"
	OrderStatusDelivered OrderStatus = "DELIVERED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)
