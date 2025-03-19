import type { PaginationParams } from '@/core/utils/pagination.utils'
import type { IngredientsDTO } from '../../ingredients/models/ingredients.model'
import type { StockMaterialCategoryDTO } from '../../stock-material-categories/models/stock-material-categories.model'
import type { UnitDTO } from '../../units/models/units.model'

export interface CreateStockMaterialDTO {
	name: string
	description?: string
	safetyStock: number
	unitId: number
	categoryId: number
	ingredientId: number
	barcode: string
	expirationPeriodInDays: number
	size: number
	isActive: boolean
}

export interface UpdateStockMaterialDTO {
	name?: string
	description?: string
	safetyStock?: number
	unitId?: number
	categoryId?: number
	ingredientId?: number
	expirationPeriodInDays?: number
	isActive?: boolean
	size?: number
}

export interface StockMaterialsDTO {
	id: number
	name: string
	description: string
	safetyStock: number
	unit: UnitDTO
	category: StockMaterialCategoryDTO
	barcode: string
	ingredient: IngredientsDTO
	expirationPeriodInDays: number
	isActive: boolean
	size: number
	createdAt: string
	updatedAt: string
}

export interface StockMaterialsFilter extends PaginationParams {
	search?: string
	lowStock?: boolean
	isActive?: boolean
	ingredientId?: number
	categoryId?: number
	supplierId?: number
	expirationInDays?: number
}

export interface GeneratedStockMaterialBarcode {
	barcode: string
}
