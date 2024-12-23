import type { PaginationParams } from '@/core/utils/pagination.utils'

export interface WarehouseStocks {
	stockMaterialId: number
	name: string
	description: string
	safetyStock: number
	expirationFlag: boolean
	quantity: number
	unitId: number
	category: string
	expiration: number
	package: Package
}

export interface InventoryLevel {
	stockMaterialId: number
	name: string
	quantity: number
}

export interface InventoryLevelsResponse {
	warehouseId: number
	levels: InventoryLevel[]
}

export interface Package {
	packageSize: number
	packageUnitId: number
}

export interface GetWarehouseStockFilter extends PaginationParams {
	search?: string
	warehouseId?: number
}
