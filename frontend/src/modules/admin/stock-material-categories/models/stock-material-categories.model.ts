import type { PaginationParams } from '@/core/utils/pagination.utils'

export interface CreateStockMaterialCategoryDTO {
	name: string
	description?: string
}

export interface UpdateStockMaterialCategoryDTO {
	name?: string
	description?: string
}

export interface StockMaterialCategoryDTO {
	id: number
	name: string
	description: string
	createdAt?: string
	updatedAt?: string
}

export interface StockMaterialCategoryFilterDTO extends PaginationParams {
	search?: string
}
