import { apiClient } from '@/core/config/axios-instance.config'
import type { Warehouse } from '@/modules/warehouse/models/warehouse.model'

class WarehouseService {
	async getWarehouses(): Promise<Warehouse[]> {
		try {
			const response = await apiClient.get<Warehouse[]>('/warehouses')
			return response.data
		} catch (error) {
			console.error('Failed to fetch warehouses:', error)
			throw error
		}
	}
}

export const warehouseService = new WarehouseService()
