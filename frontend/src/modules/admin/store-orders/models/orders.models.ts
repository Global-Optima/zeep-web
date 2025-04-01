import type { PaginationParams } from '@/core/utils/pagination.utils'
import type { UnitDTO } from '@/modules/admin/units/models/units.model'

export enum OrderStatus {
	WAITING_FOR_PAYMENT = 'WAITING_FOR_PAYMENT',
	PENDING = 'PENDING',
	PREPARING = 'PREPARING',
	COMPLETED = 'COMPLETED',
	IN_DELIVERY = 'IN_DELIVERY',
	DELIVERED = 'DELIVERED',
	CANCELLED = 'CANCELLED',
}

export enum SubOrderStatus {
	PENDING = 'PENDING',
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

export interface CheckCustomerName {
	customerName: string
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
	completedAt?: Date
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
	machineId: string
}

export interface OrderAdditiveDTO {
	id: number
	name: string
	description: string
	size: string
	machineId: string
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
	status: OrderStatus
	total: number
	suborders: SuborderDetailsDTO[]
	deliveryAddress?: OrderDeliveryAddressDTO
	completedAt?: Date
	displayNumber: number
}

export interface SuborderDetailsDTO {
	id: number
	price: number
	status: string
	storeProductSize: OrderProductSizeDetailsDTO
	storeAdditives: OrderAdditiveDetailsDTO[]
	completedAt?: Date
}

export interface OrderProductSizeDetailsDTO {
	id: number
	name: string
	unit: UnitDTO
	storePrice: number
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

export interface TransactionDTO {
	bin: string
	transactionId: string
	processId: string
	paymentMethod: string
	amount: number
	currency: string
	qrNumber?: string
	cardMask?: string
	icc?: string
}

export interface OrdersTimeZoneFilter {
	storeId?: number
	timezone?: string
	timezoneOffset?: number
	timeGapMinutes?: number
	includeYesterdayOrders?: boolean
	statuses?: OrderStatus[]
}

export interface ToggleNextSuborderStatusOptions {
	includeIfCompletedGapMinutes?: number
}
