import type { PaginationParams } from '@/core/utils/pagination.utils'
import type { UnitDTO } from '@/modules/admin/units/models/units.model'

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

export interface OrdersFilterQuery extends PaginationParams {
	search?: string
	status?: OrderStatus
	storeId?: number
}

export interface CreateOrderDTO {
	customerId?: number
	customerName: string
	employeeId?: number
	deliveryAddressId?: number
	subOrders: CreateSubOrderDTO[]
}

export interface OrderStatusesCountDTO {
	ALL: number
	PREPARING: number
	COMPLETED: number
	IN_DELIVERY: number
	DELIVERED: number
	CANCELLED: number
}

interface CreateSubOrderDTO {
	storeProductSizeId: number
	quantity: number
	storeAdditivesIds: number[]
}

export interface OrderDTO {
	id: number
	customerId?: number
	customerName: string
	employeeId?: number
	storeId: number
	deliveryAddressId?: number
	status: OrderStatus
	createdAt: string
	total: number
	subOrdersQuantity: number
	displayNumber: number
	subOrders: SuborderDTO[]
}

export interface SuborderDTO {
	id: number
	orderId: number
	productSize: OrderProductSizeDTO
	price: number
	status: SubOrderStatus
	additives: SuborderAdditiveDTO[]
	createdAt: string
	updatedAt: string
}

export interface OrderProductSizeDTO {
	id: number
	sizeName: string
	productName: string
	size: number
	unit: UnitDTO
}

export interface OrderAdditiveDTO {
	id: number
	name: string
	description: string
	size: string
}

export interface SuborderAdditiveDTO {
	id: number
	subOrderId: number
	additive: OrderAdditiveDTO
	price: number
	createdAt: string
	updatedAt: string
}

export interface OrderDetailsDTO {
	id: number
	customerName?: string
	status: string
	total: number
	suborders: SuborderDetailsDTO[]
	deliveryAddress?: OrderDeliveryAddressDTO
}

export interface SuborderDetailsDTO {
	id: number
	price: number
	status: string
	productSize: OrderProductSizeDetailsDTO
	additives: OrderAdditiveDetailsDTO[]
}

export interface OrderProductSizeDetailsDTO {
	id: number
	name: string
	unit: UnitDTO
	basePrice: number
	product: OrderProductDetailsDTO
}

export interface OrderProductDetailsDTO {
	id: number
	name: string
	description: string
	imageUrl: string
}

export interface OrderAdditiveDetailsDTO {
	id: number
	name: string
	description: string
	basePrice: number
}

export interface OrderDeliveryAddressDTO {
	id: number
	address: string
	longitude: string
	latitude: string
}

export interface OrdersExportFilterQuery {
	startDate?: string
	endDate?: string
	storeId?: number
	language?: 'kk' | 'ru' | 'en'
}
