import { apiClient } from '@/core/config/axios-instance.config'
import { type PaginatedResponse } from '@/core/utils/pagination.utils'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type {
	CreateWarehouseDTO,
	UpdateWarehouseDTO,
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

	async getById(id: number) {
		try {
			const response = await apiClient.get<WarehouseDTO>(`/warehouses/${id}`)
			return response.data
		} catch (error) {
			console.error('Failed to fetch warehouse by id:', error)
			throw error
		}
	}

	async update(id: number, dto: UpdateWarehouseDTO) {
		try {
			const response = await apiClient.put<void>(`/warehouses/${id}`, dto)
			return response.data
		} catch (error) {
			console.error('Failed to update warehouse by id:', error)
			throw error
		}
	}

	async create(dto: CreateWarehouseDTO) {
		try {
			const response = await apiClient.post<void>(`/warehouses`, dto)
			return response.data
		} catch (error) {
			console.error('Failed to create warehouse by id:', error)
			throw error
		}
	}
}

export const warehouseService = new WarehouseService()
