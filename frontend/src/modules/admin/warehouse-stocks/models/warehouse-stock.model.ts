import type { PaginationParams } from '@/core/utils/pagination.utils'
import type { StockMaterialsDTO } from '../../stock-materials/models/stock-materials.model'
import type { SupplierDTO } from '../../suppliers/models/suppliers.model'
import type { WarehouseDTO } from '@/modules/admin/warehouses/models/warehouse.model'

export interface WarehouseStocksDTO {
	stockMaterial: StockMaterialsDTO
	quantity: number
	earliestExpirationDate: string
}

export interface WarehouseStockMaterialDetailsDTO {
	stockMaterial: StockMaterialsDTO
	quantity: number
	earliestExpirationDate?: string
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

export interface AvailableWarehouseStockMaterialsFilter extends PaginationParams {
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

export interface WarehouseDeliveryDTO {
	id: number
	materials: WarehouseDeliveryStockMaterialDTO[]
	supplier: SupplierDTO
	warehouse: WarehouseDTO
	deliveryDate: Date
}

export interface WarehouseDeliveryStockMaterialDTO {
	stockMaterial: StockMaterialsDTO
	quantity: number
	barcode: string
	expirationDate: Date
}

export interface ReceiveWarehouseDelivery {
	supplierId: number
	materials: ReceiveWarehouseStockMaterial[]
}

export interface ReceiveWarehouseStockMaterial {
	stockMaterialId: number
	quantity: number
}

export interface WarehouseDeliveryFilter extends PaginationParams {
	search?: string
	warehouseId?: number
	startDate?: Date
	endDate?: Date
}
