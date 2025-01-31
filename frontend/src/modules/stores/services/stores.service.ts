import { apiClient } from '@/core/config/axios-instance.config'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type { CreateStoreDTO, StoresFilter, UpdateStoreDTO } from '../models/stores-dto.model'
import type { StoreDTO } from '../models/stores.models'
import { type PaginatedResponse } from './../../../core/utils/pagination.utils'

class StoreService {
	async getStores(filter?: StoresFilter) {
		try {
			const response = await apiClient.get<PaginatedResponse<StoreDTO[]>>('/stores/all', {
				params: buildRequestFilter(filter),
			})
			return response.data
		} catch (error) {
			console.error('Failed to fetch stores:', error)
			throw error
		}
	}

	async getStore(id: number): Promise<StoreDTO> {
		try {
			const response = await apiClient.get<StoreDTO>(`/stores/${id}`)
			return response.data
		} catch (error) {
			console.error(`Failed to fetch store by id ${id}:`, error)
			throw error
		}
	}

	async updateStore(id: number, dto: UpdateStoreDTO) {
		try {
			const response = await apiClient.put<void>(`/stores/${id}`, dto)
			return response.data
		} catch (error) {
			console.error(`Failed to fetch store by id ${id}:`, error)
			throw error
		}
	}

	async createStore(dto: CreateStoreDTO) {
		try {
			const response = await apiClient.post<void>(`/stores`, dto)
			return response.data
		} catch (error) {
			console.error(`Failed to create store:`, error)
			throw error
		}
	}
}

export const storesService = new StoreService()
