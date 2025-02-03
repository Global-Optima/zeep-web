import { apiClient } from '@/core/config/axios-instance.config'
import type { PaginatedResponse } from '@/core/utils/pagination.utils'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type {
	AllFranchiseeFilter,
	CreateFranchiseeDTO,
	FranchiseeDTO,
	FranchiseeFilterDTO,
	UpdateFranchiseeDTO,
} from '@/modules/admin/franchisees/models/franchisee.model'

class FranchiseesService {
	private readonly baseUrl: string = '/franchisees'

	async getAll(filter?: AllFranchiseeFilter) {
		try {
			const response = await apiClient.get<FranchiseeDTO[]>(`${this.baseUrl}/all`, {
				params: buildRequestFilter(filter),
			})
			return response.data
		} catch (error) {
			console.error('Failed to fetch stores:', error)
			throw error
		}
	}

	async getPaginated(filter?: FranchiseeFilterDTO) {
		try {
			const response = await apiClient.get<PaginatedResponse<FranchiseeDTO[]>>(this.baseUrl, {
				params: buildRequestFilter(filter),
			})
			return response.data
		} catch (error) {
			console.error('Failed to fetch franchisees:', error)
			throw error
		}
	}

	async getById(id: number): Promise<FranchiseeDTO> {
		try {
			const response = await apiClient.get<FranchiseeDTO>(`${this.baseUrl}/${id}`)
			return response.data
		} catch (error) {
			console.error(`Failed to fetch franchisee with ID ${id}:`, error)
			throw error
		}
	}

	async getMyFranchisee(): Promise<FranchiseeDTO> {
		try {
			const response = await apiClient.get<FranchiseeDTO>(`${this.baseUrl}/my`)
			return response.data
		} catch (error) {
			console.error('Failed to fetch current franchisee:', error)
			throw error
		}
	}

	async create(data: CreateFranchiseeDTO): Promise<FranchiseeDTO> {
		try {
			const response = await apiClient.post<FranchiseeDTO>(this.baseUrl, data)
			return response.data
		} catch (error) {
			console.error('Failed to create franchisee:', error)
			throw error
		}
	}

	async update(id: number, data: UpdateFranchiseeDTO): Promise<FranchiseeDTO> {
		try {
			const response = await apiClient.put<FranchiseeDTO>(`${this.baseUrl}/${id}`, data)
			return response.data
		} catch (error) {
			console.error(`Failed to update franchisee with ID ${id}:`, error)
			throw error
		}
	}

	async delete(id: number): Promise<void> {
		try {
			await apiClient.delete(`${this.baseUrl}/${id}`)
		} catch (error) {
			console.error(`Failed to delete franchisee with ID ${id}:`, error)
			throw error
		}
	}
}

export const franchiseeService = new FranchiseesService()
