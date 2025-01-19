import type { PaginationParams } from '@/core/utils/pagination.utils'
import type { IngredientsDTO } from '../../ingredients/models/ingredients.model'
import type { UnitDTO } from '../../units/models/units.model'

export interface StockMaterialCategoryResponse {
	id: number
	name: string
	description: string
	createdAt: Date
	updatedAt: Date
}

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
	unit: UnitDTO
	category: StockMaterialCategoryResponse
	barcode: string
	ingredient: IngredientsDTO
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
