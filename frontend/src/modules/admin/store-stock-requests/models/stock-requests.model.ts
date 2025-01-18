import type { PaginationParams } from '@/core/utils/pagination.utils'

export enum StockRequestStatus {
	CREATED = 'CREATED',
	PROCESSED = 'PROCESSED',
	IN_DELIVERY = 'IN_DELIVERY',
	COMPLETED = 'COMPLETED',
	REJECTED_BY_STORE = 'REJECTED_BY_STORE',
	REJECTED_BY_WAREHOUSE = 'REJECTED_BY_WAREHOUSE',
	ACCEPTED_WITH_CHANGE = 'ACCEPTED_WITH_CHANGE',
}

export const STOCK_REQUEST_STATUS_COLOR: Record<StockRequestStatus, string> = {
	[StockRequestStatus.CREATED]: 'bg-blue-100 text-blue-800', // Informative
	[StockRequestStatus.PROCESSED]: 'bg-indigo-100 text-indigo-800', // Progressing
	[StockRequestStatus.IN_DELIVERY]: 'bg-yellow-100 text-yellow-800', // Warning/Attention
	[StockRequestStatus.COMPLETED]: 'bg-green-100 text-green-800', // Success
	[StockRequestStatus.REJECTED_BY_STORE]: 'bg-red-100 text-red-800', // Error
	[StockRequestStatus.REJECTED_BY_WAREHOUSE]: 'bg-red-100 text-red-800', // Error
	[StockRequestStatus.ACCEPTED_WITH_CHANGE]: 'bg-purple-100 text-purple-800', // Neutral/Informational
}

export const STOCK_REQUEST_STATUS_FORMATTED: Record<StockRequestStatus, string> = {
	[StockRequestStatus.CREATED]: 'Создан',
	[StockRequestStatus.PROCESSED]: 'Обработан',
	[StockRequestStatus.IN_DELIVERY]: 'В доставке',
	[StockRequestStatus.COMPLETED]: 'Завершён',
	[StockRequestStatus.REJECTED_BY_STORE]: 'Отклонён магазином',
	[StockRequestStatus.REJECTED_BY_WAREHOUSE]: 'Отклонён складом',
	[StockRequestStatus.ACCEPTED_WITH_CHANGE]: 'Принят с изменениями',
}

export const ALL_STOCK_REQUESTS_STATUSES: StockRequestStatus[] = Object.values(StockRequestStatus)

export const WAREHOUSE_STOCK_REQUEST_STATUSES: StockRequestStatus[] =
	ALL_STOCK_REQUESTS_STATUSES.filter(status => ![StockRequestStatus.CREATED].includes(status))

export const STORE_STOCK_REQUEST_STATUSES: StockRequestStatus[] = ALL_STOCK_REQUESTS_STATUSES

export interface PackageMeasure {
	quantity: number
	unitsPerPackage: number
	packageUnit: string
	totalUnitsInStock: number
}

// DTO for creating a stock request
export interface CreateStockRequestDTO {
	items: StockRequestStockMaterialDTO[]
}

// DTO for individual stock material in a stock request
export interface StockRequestStockMaterialDTO {
	stockMaterialId: number
	quantity: number
}

// DTO for rejecting a stock request
export interface RejectStockRequestStatusDTO {
	comment: string
}

// DTO for accepting a stock request with changes
export interface AcceptWithChangeRequestStatusDTO {
	comment: string
	items: StockRequestStockMaterialDTO[]
}

// DTO for updating ingredient delivery dates
export interface UpdateIngredientDates {
	deliveredDate: string // ISO Date string
	expirationDate: string // ISO Date string
}

// Stock request response object
export interface StockRequestResponse {
	requestId: number
	store: StockRequestStoreDTO
	warehouse: StockRequestWarehouseDTO
	status: StockRequestStatus
	stockMaterials: StockRequestStockMaterialResponse[]
	createdAt: string // ISO Date string
	updatedAt: string // ISO Date string
}

// Store details in stock request
export interface StockRequestStoreDTO {
	id: number
	name: string
	address: string
}

// Warehouse details in stock request
export interface StockRequestWarehouseDTO {
	id: number
	name: string
}

// Stock material response for a stock request
export interface StockRequestStockMaterialResponse {
	stockMaterialId: number
	name: string
	category: string
	packageMeasures: PackageMeasure
}

// Response for low-stock ingredients
export interface LowStockIngredientResponse {
	ingredientId: number
	name: string
	unit: string
	quantity: number
	lowStockThreshold: number
}

// Filters for fetching stock requests
export interface GetStockRequestsFilter extends PaginationParams {
	storeId?: number
	warehouseId?: number
	startDate?: string // ISO Date string
	endDate?: string // ISO Date string
	statuses?: StockRequestStatus[] // Use a string enum for `data.StockRequestStatus`
}

// DTO for stock material details
export interface StockMaterialDTO {
	stockMaterialId: number
	name: string
	category: string
	unit: string
	availableQuantity: number
}

// Availability details for stock material
export interface StockMaterialAvailabilityDTO {
	stockMaterialId: number
	name: string
	category: string
	unit: string
	availableQuantity: number
	warehouse: StockRequestWarehouseDTO
}

// Filters for fetching stock materials
export interface StockMaterialFilter {
	category?: string
	search?: string
}
