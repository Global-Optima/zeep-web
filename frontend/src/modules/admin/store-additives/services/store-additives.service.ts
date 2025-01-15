import { apiClient } from '@/core/config/axios-instance.config'
import type { PaginatedResponse } from '@/core/utils/pagination.utils'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type { AdditiveFilterQuery } from '../../additives/models/additives.model'
import type {
	CreateStoreAdditiveDTO,
	StoreAdditiveCategoryDTO,
	StoreAdditiveDTO,
	UpdateStoreAdditiveDTO,
} from '../models/store-additves.model'

class StoreAdditivesService {
	/**
	 * Fetch all store additives with optional filters
	 */
	async getStoreAdditives(filter?: AdditiveFilterQuery) {
		try {
			const response = await apiClient.get<PaginatedResponse<StoreAdditiveDTO[]>>(
				'/store-additives',
				{
					params: buildRequestFilter(filter),
				},
			)
			return response.data
		} catch (error) {
			console.error('Failed to fetch store additives:', error)
			throw error
		}
	}

	/**
	 * Fetch details of a specific store additive by ID
	 */
	async getStoreAdditiveById(id: number) {
		try {
			const response = await apiClient.get<StoreAdditiveDTO>(`/store-additives/${id}`)
			return response.data
		} catch (error) {
			console.error(`Failed to fetch store additive details for ID ${id}:`, error)
			throw error
		}
	}

	/**
	 * Fetch all store additive categories
	 */
	async getStoreAdditiveCategories(productSizeId: number) {
		try {
			const response = await apiClient.get<StoreAdditiveCategoryDTO[]>(
				`/store-additives/categories/${productSizeId}`,
			)
			return response.data
		} catch (error) {
			console.error('Failed to fetch store additive categories:', error)
			throw error
		}
	}

	/**
	 * Create a new store additive
	 */
	async createStoreAdditive(data: CreateStoreAdditiveDTO[]) {
		try {
			const response = await apiClient.post<void>('/store-additives', data)
			return response.data
		} catch (error) {
			console.error('Failed to create store additive:', error)
			throw error
		}
	}

	/**
	 * Update an existing store additive by ID
	 */
	async updateStoreAdditive(id: number, data: UpdateStoreAdditiveDTO) {
		try {
			const response = await apiClient.put<void>(`/store-additives/${id}`, data)
			return response.data
		} catch (error) {
			console.error(`Failed to update store additive with ID ${id}:`, error)
			throw error
		}
	}

	/**
	 * Delete a store additive by ID
	 */
	async deleteStoreAdditive(id: number) {
		try {
			await apiClient.delete<void>(`/store-additives/${id}`)
		} catch (error) {
			console.error(`Failed to delete store additive with ID ${id}:`, error)
			throw error
		}
	}
}

export const storeAdditivesService = new StoreAdditivesService()
