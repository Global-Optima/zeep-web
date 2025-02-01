import { apiClient } from '@/core/config/axios-instance.config'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type {
	AllFranchiseeFilter,
	FranchiseeDTO,
} from '@/modules/admin/franchisee/models/franchisee.model'

class FranchiseesService {
	async getAllFranchisees(filter?: AllFranchiseeFilter) {
		try {
			const response = await apiClient.get<FranchiseeDTO[]>('/franchisees/all', {
				params: buildRequestFilter(filter),
			})
			return response.data
		} catch (error) {
			console.error('Failed to fetch stores:', error)
			throw error
		}
	}
}

export const franchiseeService = new FranchiseesService()
