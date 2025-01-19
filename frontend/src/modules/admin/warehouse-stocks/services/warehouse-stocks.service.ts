import { apiClient } from '@/core/config/axios-instance.config'
import type { PaginatedResponse } from '@/core/utils/pagination.utils'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type {
	GetWarehouseStockFilter,
	UpdateWarehouseStockDTO,
	WarehouseStockDetailsDTO,
	WarehouseStocksDTO,
} from '../models/warehouse-stock.model'

class WarehouseStocksService {
	async getWarehouseStocks(filter?: GetWarehouseStockFilter) {
		try {
			const response = await apiClient.get<PaginatedResponse<WarehouseStocksDTO[]>>(
				`/warehouses/stocks`,
				{ params: buildRequestFilter(filter) },
			)
			return response.data
		} catch (error) {
			console.error('Failed to fetch store stocks:', error)
			throw error
		}
	}

	async getWarehouseStockById(stockMaterialId: number) {
		try {
			const response = await apiClient.get<WarehouseStockDetailsDTO>(
				`/warehouses/stocks/${stockMaterialId}`,
			)
			return response.data
		} catch (error) {
			console.error('Failed to fetch store stock:', error)
			throw error
		}
	}

	async updateWarehouseStocksById(id: number, data: UpdateWarehouseStockDTO) {
		try {
			const response = await apiClient.put<void>(`/warehouses/stocks/${id}`, data)
			return response.data
		} catch (error) {
			console.error(`Failed to update warehouse stocks with ID ${id}: `, error)
			throw error
		}
	}
}

export const warehouseStocksService = new WarehouseStocksService()
