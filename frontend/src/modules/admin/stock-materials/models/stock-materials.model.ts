import type { PaginationParams } from '@/core/utils/pagination.utils'

export interface CreateStockMaterialDTO {
	name: string
	description?: string
	safetyStock: number
	expirationFlag: boolean
	unitId: number
	supplierId: number
	categoryId: number
	ingredientId: number
	barcode?: string
	expirationPeriodInDays: number
}

export interface UpdateStockMaterialDTO {
	name?: string
	description?: string
	safetyStock?: number
	expirationFlag?: boolean
	unitId?: number
	categoryId?: number
	ingredientId?: number
	barcode?: string
	expirationPeriodInDays?: number
	isActive?: boolean
}

export interface StockMaterialsDTO {
	id: number
	name: string
	description: string
	safetyStock: number
	expirationFlag: boolean
	unitId: number
	unitName?: string
	category: string
	barcode: string
	ingredient: string
	expirationPeriodInDays: number
	isActive: boolean
	createdAt: string
	updatedAt: string
}

export interface StockMaterialsFilter extends PaginationParams {
	search?: string
	lowStock?: boolean
	expirationFlag?: boolean
	isActive?: boolean
	supplierId?: number
	ingredientId?: number
	categoryId?: number
	expirationInDays?: number
}
