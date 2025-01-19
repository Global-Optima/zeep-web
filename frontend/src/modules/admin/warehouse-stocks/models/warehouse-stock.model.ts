import type { PaginationParams } from '@/core/utils/pagination.utils'
import type { StockMaterialsDTO } from '../../stock-materials/models/stock-materials.model'

export interface WarehouseStocksDTO {
	stockMaterial: StockMaterialsDTO
	totalQuantity: number
	earliestExpirationDate: string
}

export interface WarehouseStockDetailsDTO extends StockMaterialsDTO {
	totalQuantity: number
	earliestExpirationDate: string
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
	lowStockOnly?: boolean
	expirationDays?: number
	search?: string
}

export interface UpdateWarehouseStockDTO {
	quantity?: number
	expirationDate: Date
}
