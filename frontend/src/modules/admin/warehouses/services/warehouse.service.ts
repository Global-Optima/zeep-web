import { apiClient } from '@/core/config/axios-instance.config'
import { type PaginatedResponse } from '@/core/utils/pagination.utils'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type {
	WarehouseDTO,
	WarehouseFilter,
} from '@/modules/admin/warehouses/models/warehouse.model'

class WarehouseService {
	async getAll() {
		try {
			const response = await apiClient.get<WarehouseDTO[]>('/warehouses/all')
			return response.data
		} catch (error) {
			console.error('Failed to fetch warehouses:', error)
			throw error
		}
	}

	async getPaginated(filter?: WarehouseFilter) {
		try {
			const response = await apiClient.get<PaginatedResponse<WarehouseDTO[]>>('/warehouses', {
				params: buildRequestFilter(filter),
			})
			return response.data
		} catch (error) {
			console.error('Failed to fetch warehouses:', error)
			throw error
		}
	}
}

export const warehouseService = new WarehouseService()
