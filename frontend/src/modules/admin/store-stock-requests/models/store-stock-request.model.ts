import type { PaginationParams } from '@/core/utils/pagination.utils'

export enum StoreStockRequestStatus {
	CREATED = 'CREATED',
	PROCESSED = 'PROCESSED',
	IN_DELIVERY = 'IN_DELIVERY',
	COMPLETED = 'COMPLETED',
	REJECTED = 'REJECTED',
}

export const ALL_STOCK_REQUEST_STATUSES: StoreStockRequestStatus[] = [
	StoreStockRequestStatus.CREATED,
	StoreStockRequestStatus.PROCESSED,
	StoreStockRequestStatus.IN_DELIVERY,
	StoreStockRequestStatus.COMPLETED,
	StoreStockRequestStatus.REJECTED,
]

export const WAREHOUSE_STOCK_REQUEST_STATUSES: StoreStockRequestStatus[] = [
	StoreStockRequestStatus.PROCESSED,
	StoreStockRequestStatus.IN_DELIVERY,
	StoreStockRequestStatus.COMPLETED,
	StoreStockRequestStatus.REJECTED,
]

export interface GetStoreStockRequestsFilter extends PaginationParams {
	storeId?: number
	warehouseId?: number
	statuses?: StoreStockRequestStatus[]
	startDate?: string
	endDate?: string
}

export interface CreateStoreStockRequestDTO {
	storeId: number
	items: CreateStoreStockRequestItemDTO[]
}

export interface CreateStoreStockRequestItemDTO {
	stockMaterialId: number
	quantity: number
}

export interface UpdateStoreStockRequestStatusDTO {
	status: StoreStockRequestStatus
}

export interface StoreStockRequestResponse {
	requestId: number
	storeId: number
	storeName: string
	warehouseId: number
	warehouseName: string
	status: StoreStockRequestStatus
	items: StoreStockRequestItemResponse[]
	createdAt: string
	updatedAt: string
}

export interface StoreStockRequestItemResponse {
	stockMaterialId: number
	name: string
	category: string
	unit: string
	quantity: number
}
