import { apiClient } from '@/core/config/axios-instance.config'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type {
	CreateIngredientCategoryDTO,
	CreateIngredientDTO,
	IngredientCategoryDTO,
	IngredientCategoryFilter,
	IngredientFilter,
	IngredientsDTO,
	UpdateIngredientCategoryDTO,
	UpdateIngredientDTO,
} from '../models/ingredients.model'
import type { PaginatedResponse } from './../../../../core/utils/pagination.utils'

class IngredientsService {
	async getIngredients(filter?: IngredientFilter) {
		const response = await apiClient.get<PaginatedResponse<IngredientsDTO[]>>('/ingredients', {
			params: buildRequestFilter(filter),
		})
		return response.data
	}

	async getIngredientById(id: number) {
		const response = await apiClient.get<IngredientsDTO>(`/ingredients/${id}`)
		return response.data
	}

	async createIngredient(dto: CreateIngredientDTO) {
		const response = await apiClient.post<void>('/ingredients', dto)
		return response.data
	}

	async updateIngredient(id: number, dto: UpdateIngredientDTO) {
		const response = await apiClient.put<void>(`/ingredients/${id}`, dto)
		return response.data
	}

	async deleteIngredient(id: number) {
		await apiClient.delete<void>(`/ingredients/${id}`)
	}

	async getIngredientCategories(filter?: IngredientCategoryFilter) {
		const response = await apiClient.get<PaginatedResponse<IngredientCategoryDTO[]>>(
			'/ingredient-categories',
			{ params: buildRequestFilter(filter) },
		)
		return response.data
	}

	async getIngredientCategoryById(id: number) {
		const response = await apiClient.get<IngredientCategoryDTO>(`/ingredient-categories/${id}`)
		return response.data
	}

	async createIngredientCategory(dto: CreateIngredientCategoryDTO) {
		const response = await apiClient.post<void>('/ingredient-categories', dto)
		return response.data
	}

	async updateIngredientCategory(id: number, dto: UpdateIngredientCategoryDTO) {
		const response = await apiClient.put<void>(`/ingredient-categories/${id}`, dto)
		return response.data
	}

	async deleteIngredientCategory(id: number) {
		await apiClient.delete<void>(`/ingredient-categories/${id}`)
	}
}

export const ingredientsService = new IngredientsService()
