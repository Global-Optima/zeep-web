import type { PaginationParams } from '@/core/utils/pagination.utils'
import type { WarehouseDTO } from '@/modules/warehouse/models/warehouse.model'
import type { StockMaterialsDTO } from '../../stock-materials/models/stock-materials.model'
import type {
	PackageMeasure,
	PackageMeasureWithQuantity,
} from '../../store-stock-requests/models/stock-requests.model'
import type { SupplierDTO } from '../../suppliers/models/suppliers.model'

export interface WarehouseStocksDTO {
	stockMaterial: StockMaterialResponse
	earliestExpirationDate: string
}

export interface StockMaterialResponse extends StockMaterialsDTO {
	packageMeasures: PackageMeasureWithQuantity
}

export interface WarehouseStockMaterialDetailsDTO {
	stockMaterial: StockMaterialsDTO
	packageMeasure: PackageMeasureWithQuantity
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

export interface WarehouseDeliveriesDTO {
	id: number
	barcode: string
	quantity: number
	materials: WarehouseDeliveryStockMaterialDTO[]
	supplier: SupplierDTO
	warehouse: WarehouseDTO
	deliveryDate: Date
	expirationDate: Date
}

export interface WarehouseDeliveryStockMaterialDTO {
	stockMaterial: StockMaterialsDTO
	package: PackageMeasure
	quantity: number
}

export interface ReceiveWarehouseDelivery {
	supplierId: number
	materials: ReceiveWarehouseStockMaterial[]
}

export interface ReceiveWarehouseStockMaterial {
	stockMaterialId: number
	quantity: number
	packageId: number
}

export interface WarehouseDeliveriesFilterDTO extends PaginationParams {
	search?: string
	warehouseId?: number
	startDate?: Date
	endDate?: Date
}
