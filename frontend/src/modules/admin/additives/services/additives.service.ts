import { apiClient } from '@/core/config/axios-instance.config'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import { buildFormData } from '@/core/utils/request-form-data-builder.utils'
import type {
	AdditiveCategoriesFilterQuery,
	AdditiveCategoryDetailsDTO,
	AdditiveDetailsDTO,
	AdditiveDTO,
	AdditiveFilterQuery,
	CreateAdditiveCategoryDTO,
	CreateAdditiveDTO,
	UpdateAdditiveCategoryDTO,
	UpdateAdditiveDTO,
} from '../models/additives.model'
import type { PaginatedResponse } from '@/core/utils/pagination.utils'

class AdditiveService {
	async getAdditives(filter?: AdditiveFilterQuery) {
		const response = await apiClient.get<PaginatedResponse<AdditiveDTO[]>>('/additives', {
			params: buildRequestFilter(filter),
		})
		return response.data
	}

	async getAdditiveById(id: number) {
		const response = await apiClient.get<AdditiveDetailsDTO>(`/additives/${id}`)
		return response.data
	}

	async createAdditive(dto: CreateAdditiveDTO) {
		const formData = buildFormData<CreateAdditiveDTO>(dto)

		const response = await apiClient.post<void>('/additives', formData, {
			headers: {
				'Content-Type': 'multipart/form-data',
			},
		})
		return response.data
	}

	async updateAdditive(id: number, dto: UpdateAdditiveDTO) {
		const formData = buildFormData<UpdateAdditiveDTO>(dto)

		const response = await apiClient.put<void>(`/additives/${id}`, formData, {
			headers: {
				'Content-Type': 'multipart/form-data',
			},
		})
		return response.data
	}

	async deleteAdditive(id: number) {
		await apiClient.delete(`/additives/${id}`)
	}

	async getAdditiveCategories(filter?: AdditiveCategoriesFilterQuery) {
		const response = await apiClient.get<PaginatedResponse<AdditiveCategoryDetailsDTO[]>>(
			'/additives/categories',
			{
				params: buildRequestFilter(filter),
			},
		)
		return response.data
	}

	// Fetch a single additive category by ID
	async getAdditiveCategoryById(id: number) {
		const response = await apiClient.get<AdditiveCategoryDetailsDTO>(`/additives/categories/${id}`)
		return response.data
	}

	// Create a new additive category
	async createAdditiveCategory(dto: CreateAdditiveCategoryDTO) {
		const response = await apiClient.post<void>('/additives/categories', dto)
		return response.data
	}

	// Update an existing additive category
	async updateAdditiveCategory(id: number, dto: UpdateAdditiveCategoryDTO) {
		const response = await apiClient.put<void>(`/additives/categories/${id}`, dto)
		return response.data
	}

	// Delete an additive category by ID
	async deleteAdditiveCategory(id: number) {
		await apiClient.delete(`/additives/categories/${id}`)
	}
}

export const additivesService = new AdditiveService()
