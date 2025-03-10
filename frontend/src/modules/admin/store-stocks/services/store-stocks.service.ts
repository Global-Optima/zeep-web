import { apiClient } from '@/core/config/axios-instance.config'
import type { PaginatedResponse } from '@/core/utils/pagination.utils'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type {
	AddMultipleStoreWarehouseStockDTO,
	AddStoreWarehouseStockDTO,
	GetStoreWarehouseStockFilterQuery,
	StoreWarehouseStockDTO,
	UpdateStoreWarehouseStockDTO,
} from '../models/store-stock.model'
import type { IngredientFilter, IngredientsDTO } from '@/modules/admin/ingredients/models/ingredients.model'

class StoreStockService {
	/**
	 * Get a paginated list of store stocks with optional filters.
	 */
	async getStoreWarehouseStockList(filter?: GetStoreWarehouseStockFilterQuery) {
		try {
			const response = await apiClient.get<PaginatedResponse<StoreWarehouseStockDTO[]>>(
				'/store-stocks',
				{
					params: buildRequestFilter(filter),
				},
			)
			return response.data
		} catch (error) {
			console.error('Failed to fetch store warehouse stock list: ', error)
			throw error
		}
	}

	async getAvailableIngredients(filter?: IngredientFilter) {
		const response = await apiClient.get<PaginatedResponse<IngredientsDTO[]>>(
			'/store-stocks/available-to-add',
			{
				params: buildRequestFilter(filter),
			},
		)
		return response.data
	}

	/**
	 * Get a single store warehouse stock by ID.
	 */
	async getStoreWarehouseStockById(id: number) {
		try {
			const response = await apiClient.get<StoreWarehouseStockDTO>(`/store-stocks/${id}`)
			return response.data
		} catch (error) {
			console.error(`Failed to fetch store warehouse stock with ID ${id}: `, error)
			throw error
		}
	}

	/**
	 * Add a new store warehouse stock item.
	 */
	async addStoreWarehouseStock(data: AddStoreWarehouseStockDTO) {
		try {
			const response = await apiClient.post<void>('/store-stocks', data)
			return response.data
		} catch (error) {
			console.error('Failed to add store warehouse stock: ', error)
			throw error
		}
	}

	/**
	 * Add multiple store warehouse stock items.
	 */
	async addMultipleStoreWarehouseStock(data: AddMultipleStoreWarehouseStockDTO) {
		try {
			const response = await apiClient.post<void>('/store-stocks/multiple', data)
			return response.data
		} catch (error) {
			console.error('Failed to add multiple store warehouse stock items: ', error)
			throw error
		}
	}

	/**
	 * Update a store warehouse stock item by ID.
	 */
	async updateStoreWarehouseStockById(id: number, data: UpdateStoreWarehouseStockDTO) {
		try {
			const response = await apiClient.put<void>(`/store-stocks/${id}`, data)
			return response.data
		} catch (error) {
			console.error(`Failed to update store warehouse stock with ID ${id}: `, error)
			throw error
		}
	}

	/**
	 * Delete a store warehouse stock item by ID.
	 */
	async deleteStoreWarehouseStockById(id: number) {
		try {
			const response = await apiClient.delete<void>(`/store-stocks/${id}`)
			return response.data
		} catch (error) {
			console.error(`Failed to delete store warehouse stock with ID ${id}: `, error)
			throw error
		}
	}
}

export const storeStocksService = new StoreStockService()
