import { apiClient } from '@/core/config/axios-instance.config'
import type {
	CreateIngredientDTO,
	IngredientFilter,
	IngredientsDTO,
	UpdateIngredientDTO,
} from '../models/ingredients.model'
import type { PaginatedResponse } from './../../../../core/utils/pagination.utils'

export class IngredientsService {
	async getIngredients(filter?: IngredientFilter) {
		const response = await apiClient.get<PaginatedResponse<IngredientsDTO[]>>('/ingredients', {
			params: filter,
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
}

export const ingredientsService = new IngredientsService()
