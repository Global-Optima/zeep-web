import { apiClient } from '@/core/config/axios-instance.config'
import type { PaginatedResponse } from '@/core/utils/pagination.utils'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type { GetWarehouseStockFilter, WarehouseStocks } from '../models/warehouse-stock.model'

class WarehouseStocksService {
	async getWarehouseStocks(filter?: GetWarehouseStockFilter) {
		try {
			const response = await apiClient.get<PaginatedResponse<WarehouseStocks[]>>(
				`/warehouse-stocks`,
				{ params: buildRequestFilter(filter) },
			)
			return response.data
		} catch (error) {
			console.error('Failed to fetch store stocks:', error)
			throw error
		}
	}
}

export const warehouseStocksService = new WarehouseStocksService()
