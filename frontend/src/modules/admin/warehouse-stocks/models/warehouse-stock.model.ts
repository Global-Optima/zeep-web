import type { PaginationParams } from '@/core/utils/pagination.utils'
import type { StockMaterialsDTO } from '../../stock-materials/models/stock-materials.model'
import type { PackageMeasure } from '../../store-stock-requests/models/stock-requests.model'

export interface WarehouseStocksDTO {
	stockMaterial: StockMaterialResponse
	earliestExpirationDate: string
}

export interface StockMaterialResponse extends StockMaterialsDTO {
	packageMeasures: PackageMeasure
}

export interface WarehouseStockMaterialDetailsDTO {
	stockMaterial: StockMaterialsDTO
	packageMeasure: PackageMeasure
	earliestExpirationDate?: string
	deliveries: WarehouseStockMaterialDeliveryDTO[]
}

export interface WarehouseStockMaterialDeliveryDTO {
	supplierName: string
	quantity: number
	deliveryDate: Date
	expirationDate: Date
}

export interface GetWarehouseStockFilter extends PaginationParams {
	warehouseId?: number
	stockMaterialId?: number
	ingredientId?: number
	lowStockOnly?: boolean
	isExpiring?: boolean
	categoryId?: number
	daysToExpire?: number
	search?: string
}

export interface UpdateWarehouseStockDTO {
	quantity: number
	expirationDate?: Date
}

export interface AddMultipleWarehouseStockDTO {
	stockMaterialId: number
	quantity: number
}
