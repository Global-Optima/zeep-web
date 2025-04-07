import { apiClient } from '@/core/config/axios-instance.config'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'

import type { PaginatedResponse } from '@/core/utils/pagination.utils'
import type {
	CreateStoreProvisionDTO,
	StoreProvisionDetailsDTO,
	StoreProvisionDTO,
	StoreProvisionFilter,
	UpdateStoreProvisionDTO,
} from '../models/store-provision.models'

class StoreProvisionsService {
	async getStoreProvisions(filter?: StoreProvisionFilter) {
		const response = await apiClient.get<PaginatedResponse<StoreProvisionDTO[]>>(
			'/store-provisions',
			{
				params: buildRequestFilter(filter),
			},
		)
		return response.data
	}

	async getStoreProvisionById(id: number) {
		const response = await apiClient.get<StoreProvisionDetailsDTO>(`/store-provisions/${id}`)
		return response.data
	}

	async createStoreProvision(dto: CreateStoreProvisionDTO) {
		const response = await apiClient.post<void>('/store-provisions', dto)
		return response.data
	}

	async completeStoreProvision(id: number) {
		const response = await apiClient.post<void>(`/store-provisions/${id}/complete`)
		return response.data
	}

	async updateStoreProvision(id: number, dto: UpdateStoreProvisionDTO) {
		const response = await apiClient.put<void>(`/store-provisions/${id}`, dto)
		return response.data
	}

	async deleteStoreProvision(id: number) {
		await apiClient.delete<void>(`/store-provisions/${id}`)
	}
}

export const storeProvisionsService = new StoreProvisionsService()
