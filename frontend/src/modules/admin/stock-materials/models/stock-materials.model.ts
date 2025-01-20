import type { PaginationParams } from '@/core/utils/pagination.utils'
import type { SuppliersDTO } from '@/modules/admin/suppliers/models/suppliers.model'
import type { IngredientsDTO } from '../../ingredients/models/ingredients.model'
import type { StockMaterialCategoryDTO } from '../../stock-material-categories/models/stock-material-categories.model'
import type { UnitDTO } from '../../units/models/units.model'

export interface CreateStockMaterialDTO {
	name: string
	description?: string
	safetyStock: number
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
	unitId?: number
	categoryId?: number
	supplierId: number
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
	unit: UnitDTO
	supplier: SuppliersDTO
	category: StockMaterialCategoryDTO
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
	isActive?: boolean
	supplierId?: number
	ingredientId?: number
	categoryId?: number
	expirationInDays?: number
}
