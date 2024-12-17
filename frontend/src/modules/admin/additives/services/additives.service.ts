import { apiClient } from '@/core/config/axios-instance.config'
import type { PaginatedResponse } from '@/core/utils/pagination.utils'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type {
	AdditiveCategories,
	AdditiveCategoriesFilterQuery,
	AdditiveFilterQuery,
	Additives,
} from '../models/additives.model'

class AdditivesService {
	async getAdditiveCategories(filter?: AdditiveCategoriesFilterQuery) {
		try {
			const response = await apiClient.get<AdditiveCategories[]>('/additives/categories', {
				params: buildRequestFilter(filter),
			})
			return response.data
		} catch (error) {
			console.error('Failed to fetch additive categories:', error)
			throw error
		}
	}

	async getAdditives(filter?: AdditiveFilterQuery) {
		try {
			const response = await apiClient.get<PaginatedResponse<Additives[]>>('/additives', {
				params: buildRequestFilter(filter),
			})
			return response.data
		} catch (error) {
			console.error('Failed to fetch additives:', error)
			throw error
		}
	}
}

export const additivesService = new AdditivesService()
