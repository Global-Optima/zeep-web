import type { PaginationParams } from '@/core/utils/pagination.utils'

export enum StoreStockRequestStatus {
	CREATED = 'CREATED',
	PROCESSED = 'PROCESSED',
	IN_DELIVERY = 'IN_DELIVERY',
	COMPLETED = 'COMPLETED',
	REJECTED = 'REJECTED',
}

export interface GetStoreStockRequestsFilter extends PaginationParams {
	storeId?: number
	warehouseId?: number
	status?: StoreStockRequestStatus
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
