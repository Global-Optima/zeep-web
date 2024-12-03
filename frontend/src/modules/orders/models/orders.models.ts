export enum OrderStatus {
	PENDING = 'PENDING',
	COMPLETED = 'COMPLETED',
	CONFIRMED = 'CONFIRMED',
	SHIPPED = 'SHIPPED',
	DELIVERED = 'DELIVERED',
	CANCELLED = 'CANCELLED',
}

export interface CreateOrderItemDTO {
	productSizeId: number
	quantity: number
	additivesIds: number[]
}

export interface CreateOrderDTO {
	customerId?: number
	customerName: string
	employeeId?: number
	storeId: number
	deliveryAddressId?: number
	orderItems: CreateOrderItemDTO[]
}

export interface OrderProductAdditiveDTO {
	id: number
	orderProductId: number
	additiveId: number
	price: number
	createdAt: string
	updatedAt: string
}

export interface OrderProductDTO {
	id: number
	orderId: number
	productSizeId: number
	quantity: number
	price: number
	additives?: OrderProductAdditiveDTO[]
	createdAt: string
	updatedAt: string
}

export interface OrderDTO {
	id: number
	customerId: number
	employeeId?: number
	storeId?: number
	deliveryAddressId?: number
	orderStatus: OrderStatus
	orderDate: string
	total: number
	orderProducts?: OrderProductDTO[]
}

export interface SubOrderEvent {
	subOrderId: number
	status: OrderStatus
}

export interface OrderEvent {
	orderId: number
	storeId?: number
	status: OrderStatus
	timestamp: string
	items: SubOrderEvent[]
}
