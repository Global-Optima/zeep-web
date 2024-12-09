export enum OrderStatus {
	PREPARING = 'PREPARING',
	COMPLETED = 'COMPLETED',
	IN_DELIVERY = 'IN_DELIVERY',
	DELIVERED = 'DELIVERED',
	CANCELLED = 'CANCELLED',
}

export enum SubOrderStatus {
	PREPARING = 'PREPARING',
	COMPLETED = 'COMPLETED',
}

export interface OrderStatusesCountDTO {
	ALL: number
	PREPARING: number
	COMPLETED: number
	IN_DELIVERY: number
	DELIVERED: number
	CANCELLED: number
}

export interface CreateSubOrderDTO {
	productSizeId: number
	quantity: number
	additivesIds: number[]
}

export interface CreateOrderDTO {
	customerName: string
	storeId: number
	deliveryAddressId?: number
	subOrders: CreateSubOrderDTO[]
}

export interface SubOrderAdditiveDTO {
	id: number
	subOrderId: number
	additive: OrderAdditiveDTO
	price: number
	createdAt: string
	updatedAt: string
}

export interface SubOrderDTO {
	id: number
	orderId: number
	productSize: OrderProductSizeDTO
	price: number
	status: SubOrderStatus
	additives: SubOrderAdditiveDTO[]
	createdAt: string
	updatedAt: string
}

export interface OrderDTO {
	id: number
	customerId?: number
	customerName?: string
	employeeId?: number
	storeId: number
	deliveryAddressId?: number
	status: OrderStatus
	createdAt: string
	total: number
	subOrdersQuantity: number
	subOrders: SubOrderDTO[]
}

export interface OrderProductSizeDTO {
	id: number
	sizeName: string
	productName: string
	size: string
}

export interface OrderAdditiveDTO {
	id: number
	name: string
	description: string
	size: string
}
