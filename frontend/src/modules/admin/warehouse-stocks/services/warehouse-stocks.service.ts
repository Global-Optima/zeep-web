import { apiClient } from '@/core/config/axios-instance.config'
import type { PaginatedResponse } from '@/core/utils/pagination.utils'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type { StockMaterialsDTO } from '@/modules/admin/stock-materials/models/stock-materials.model'
import type {
	AddMultipleWarehouseStockDTO,
	AvailableWarehouseStockMaterialsFilter,
	GetWarehouseStockFilter,
	ReceiveWarehouseDelivery,
	UpdateWarehouseStockDTO,
	WarehouseDeliveryDTO,
	WarehouseDeliveryFilter,
	WarehouseStockMaterialDetailsDTO,
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
			const response = await apiClient.get<WarehouseStockMaterialDetailsDTO>(
				`/warehouses/stocks/${stockMaterialId}`,
			)
			return response.data
		} catch (error) {
			console.error('Failed to fetch store stock:', error)
			throw error
		}
	}

	async getAvailableStockMaterials(filter?: AvailableWarehouseStockMaterialsFilter) {
		try {
			const response = await apiClient.get<PaginatedResponse<StockMaterialsDTO[]>>(
				`/warehouses/stocks/available-to-add`,
				{ params: buildRequestFilter(filter) },
			)
			return response.data
		} catch (error) {
			console.error('Failed to fetch available warehouse stock materials:', error)
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

	async addMultipleWarehouseStock(dto: AddMultipleWarehouseStockDTO[]) {
		try {
			const response = await apiClient.post<void>(`/warehouses/stocks/add`, dto)
			return response.data
		} catch (error) {
			console.error(`Failed to add multiple warehouse stocks: `, error)
			throw error
		}
	}

	async getWarehouseDeliveries(filter?: WarehouseDeliveryFilter) {
		try {
			const response = await apiClient.get<PaginatedResponse<WarehouseDeliveryDTO[]>>(
				`/warehouses/stocks/deliveries`,
				{ params: buildRequestFilter(filter) },
			)
			return response.data
		} catch (error) {
			console.error('Failed to fetch warehouse stock deliveries:', error)
			throw error
		}
	}

	async getWarehouseDeliveryId(deliveryId: number) {
		try {
			const response = await apiClient.get<WarehouseDeliveryDTO>(
				`/warehouses/stocks/deliveries/${deliveryId}`,
			)
			return response.data
		} catch (error) {
			console.error('Failed to fetch warehouse stock delivery by id:', error)
			throw error
		}
	}

	async receiveWarehouseDelivery(dto: ReceiveWarehouseDelivery) {
		try {
			const response = await apiClient.post<void>(`/warehouses/stocks/receive`, dto)
			return response.data
		} catch (error) {
			console.error('Failed to receive warehouse stock delivery:', error)
			throw error
		}
	}
}

export const warehouseStocksService = new WarehouseStocksService()
