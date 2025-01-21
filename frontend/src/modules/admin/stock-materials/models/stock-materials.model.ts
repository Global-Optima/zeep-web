import type { PaginationParams } from '@/core/utils/pagination.utils'
import type { IngredientsDTO } from '../../ingredients/models/ingredients.model'
import type { StockMaterialCategoryDTO } from '../../stock-material-categories/models/stock-material-categories.model'
import type { UnitDTO } from '../../units/models/units.model'
import type { PackageMeasure } from '../../store-stock-requests/models/stock-requests.model'

export interface CreateStockMaterialDTO {
	name: string
	description?: string
	safetyStock: number
	unitId: number
	categoryId: number
	ingredientId: number
	barcode: string
	expirationPeriodInDays: number

	packages: CreateStockMaterialPackagesDTO[]
}

export interface CreateStockMaterialPackagesDTO {
	size: number
	unitId: number
}

export interface UpdateStockMaterialPackagesDTO {
	size: number
	unitId: number
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
  packages: UpdateStockMaterialPackagesDTO[]
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
	createdAt: string
	updatedAt: string
  packages: PackageMeasure[]
}

export interface StockMaterialsFilter extends PaginationParams {
	search?: string
	lowStock?: boolean
	isActive?: boolean
	ingredientId?: number
	categoryId?: number
	expirationInDays?: number
}
