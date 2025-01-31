import { apiClient } from '@/core/config/axios-instance.config'
import type { PaginatedResponse } from '@/core/utils/pagination.utils'
import type { WarehouseDTO } from '@/modules/warehouse/models/warehouse.model'

class WarehouseService {
	async getWarehouses() {
		try {
			const response = await apiClient.get<PaginatedResponse<WarehouseDTO[]>>('/warehouses/all')
			return response.data
		} catch (error) {
			console.error('Failed to fetch warehouses:', error)
			throw error
		}
	}
}

export const warehouseService = new WarehouseService()
