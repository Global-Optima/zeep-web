import { apiClient } from '@/core/config/axios-instance.config'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type { StoresFilter } from '../models/stores-dto.model'
import type { Store } from '../models/stores.models'

class StoreService {
	async getStores(filter?: StoresFilter): Promise<Store[]> {
		try {
			const response = await apiClient.get<Store[]>('/stores', {
				params: buildRequestFilter(filter),
			})
			return response.data
		} catch (error) {
			console.error('Failed to fetch stores:', error)
			throw error
		}
	}

	async getStore(id: number): Promise<Store> {
		try {
			const response = await apiClient.get<Store>(`/stores/${id}`)
			return response.data
		} catch (error) {
			console.error(`Failed to fetch store by id ${id}:`, error)
			throw error
		}
	}

	async updateStore(id: number, dto: Partial<Store>) {
		try {
			const response = await apiClient.put<void>(`/stores/${id}`, dto)
			return response.data
		} catch (error) {
			console.error(`Failed to fetch store by id ${id}:`, error)
			throw error
		}
	}

	async createStore(dto: Store) {
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
