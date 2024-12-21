import { apiClient } from '@/core/config/axios-instance.config'
import type { PaginatedResponse } from '@/core/utils/pagination.utils'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type {
	CreateMultipleStoreStock,
	StoreStocks,
	StoreStocksFilter,
	UpdateStoreStock,
} from '@/modules/admin/store-warehouse/models/store-stock.model'

class StoreStocksService {
	async getStoreStocks(storeId: number, filter?: StoreStocksFilter) {
		try {
			const response = await apiClient.get<PaginatedResponse<StoreStocks[]>>(
				`/store-warehouse-stock/${storeId}`,
				{
					params: buildRequestFilter(filter),
				},
			)
			return response.data
		} catch (error) {
			console.error('Failed to fetch store stocks:', error)
			throw error
		}
	}

	async getStoreStock(storeId: number, id: number) {
		try {
			const response = await apiClient.get<StoreStocks>(`/store-warehouse-stock/${storeId}/${id}`)
			return response.data
		} catch (error) {
			console.error(`Failed to fetch store stock by id ${id}:`, error)
			throw error
		}
	}

	async updateStoreStock(storeId: number, id: number, dto: UpdateStoreStock) {
		try {
			const response = await apiClient.put<void>(`/store-warehouse-stock/${storeId}/${id}`, dto)
			return response.data
		} catch (error) {
			console.error(`Failed to update store stock by id ${id}:`, error)
			throw error
		}
	}

	async createMultipleStoreStock(storeId: number, dto: CreateMultipleStoreStock) {
		try {
			const response = await apiClient.post<void>(`/store-warehouse-stock/${storeId}`, dto)
			return response.data
		} catch (error) {
			console.error(`Failed to create store stock:`, error)
			throw error
		}
	}
}

export const storeStocksService = new StoreStocksService()
