import type { PaginationParams } from '@/core/utils/pagination.utils'
import type { MachineCategory } from '@/modules/admin/product-categories/utils/category-options'

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
  machineCategory: MachineCategory
}

export interface StockMaterialCategoryFilterDTO extends PaginationParams {
	search?: string
}
