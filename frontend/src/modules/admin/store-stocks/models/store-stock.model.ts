import type { IngredientsDTO } from '../../ingredients/models/ingredients.model'

export interface StoreWarehouseStockDTO {
	id: number
	name: string
	quantity: number
	lowStockAlert: boolean
	lowStockThreshold: number
	ingredient: IngredientsDTO
}

export interface GetStoreWarehouseStockFilterQuery {
	search?: string
	lowStockOnly?: boolean
	page?: number
	pageSize?: number
}

export interface AddStoreWarehouseStockDTO {
	ingredientId: number
	quantity: number
	lowStockThreshold: number
}

export interface AddMultipleStoreWarehouseStockDTO {
	ingredientStocks: AddStoreWarehouseStockDTO[]
}

export interface UpdateStoreWarehouseStockDTO {
	quantity?: number
	lowStockThreshold?: number
}
