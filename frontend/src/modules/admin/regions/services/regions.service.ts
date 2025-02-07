import { apiClient } from '@/core/config/axios-instance.config'
import type { PaginatedResponse } from '@/core/utils/pagination.utils'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type {
	AllRegionsFilter,
	CreateRegionDTO,
	RegionDTO,
	RegionFilterDTO,
	UpdateRegionDTO,
} from '@/modules/admin/regions/models/regions.model'

class RegionsService {
	private readonly baseUrl: string = '/regions'

	async getAll(filter?: AllRegionsFilter) {
		try {
			const response = await apiClient.get<RegionDTO[]>('/regions/all', {
				params: buildRequestFilter(filter),
			})
			return response.data
		} catch (error) {
			console.error('Failed to fetch stores:', error)
			throw error
		}
	}

	async getPaginated(filter?: RegionFilterDTO) {
		try {
			const response = await apiClient.get<PaginatedResponse<RegionDTO[]>>(this.baseUrl, {
				params: buildRequestFilter(filter),
			})
			return response.data
		} catch (error) {
			console.error('Failed to fetch regions:', error)
			throw error
		}
	}

	async getById(id: number): Promise<RegionDTO> {
		try {
			const response = await apiClient.get<RegionDTO>(`${this.baseUrl}/${id}`)
			return response.data
		} catch (error) {
			console.error(`Failed to fetch region with ID ${id}:`, error)
			throw error
		}
	}

	async create(data: CreateRegionDTO): Promise<RegionDTO> {
		try {
			const response = await apiClient.post<RegionDTO>(this.baseUrl, data)
			return response.data
		} catch (error) {
			console.error('Failed to create region:', error)
			throw error
		}
	}

	async update(id: number, data: UpdateRegionDTO): Promise<RegionDTO> {
		try {
			const response = await apiClient.put<RegionDTO>(`${this.baseUrl}/${id}`, data)
			return response.data
		} catch (error) {
			console.error(`Failed to update region with ID ${id}:`, error)
			throw error
		}
	}

	async delete(id: number): Promise<void> {
		try {
			await apiClient.delete(`${this.baseUrl}/${id}`)
		} catch (error) {
			console.error(`Failed to delete region with ID ${id}:`, error)
			throw error
		}
	}
}

export const regionsService = new RegionsService()
