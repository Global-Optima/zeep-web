import { apiClient } from '@/core/config/axios-instance.config'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'

import type { PaginatedResponse } from '@/core/utils/pagination.utils'
import type {
	CreateProvisionDTO,
	ProvisionDetailsDTO,
	ProvisionDTO,
	ProvisionFilter,
	ProvisionTechnicalMap,
	UpdateProvisionDTO,
} from '../models/provision.models'

class ProvisionsService {
	async getProvisions(filter?: ProvisionFilter) {
		const response = await apiClient.get<PaginatedResponse<ProvisionDTO[]>>('/provisions', {
			params: buildRequestFilter(filter),
		})
		return response.data
	}

	async getProvisionById(id: number) {
		const response = await apiClient.get<ProvisionDetailsDTO>(`/provisions/${id}`)
		return response.data
	}

	async createProvision(dto: CreateProvisionDTO) {
		const response = await apiClient.post<void>('/provisions', dto)
		return response.data
	}

	async updateProvision(id: number, dto: UpdateProvisionDTO) {
		const response = await apiClient.put<void>(`/provisions/${id}`, dto)
		return response.data
	}

	async deleteProvision(id: number) {
		await apiClient.delete<void>(`/provisions/${id}`)
	}

	async getProvisionTechMap(id: number) {
		try {
			const response = await apiClient.get<ProvisionTechnicalMap>(`/provisions/${id}/technical-map`)
			return response.data
		} catch (error) {
			console.error(`Failed to fetch provision technical map for ID ${id}: `, error)
			throw error
		}
	}
}

export const provisionsService = new ProvisionsService()
