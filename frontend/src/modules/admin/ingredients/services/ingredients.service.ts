import { apiClient } from '@/core/config/axios-instance.config'
import type { PaginatedResponse } from '@/core/utils/pagination.utils'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type { StockMaterialsDTO } from '@/modules/admin/stock-materials/models/stock-materials.model'
import type { StoreStocksFilter } from '@/modules/admin/store-stocks/models/store-stock.model'

class IngredientsService {
	async getIngredients(filter?: StoreStocksFilter) {
		try {
			const response = await apiClient.get<PaginatedResponse<StockMaterialsDTO[]>>(`/ingredients`, {
				params: buildRequestFilter(filter),
			})
			return response.data
		} catch (error) {
			console.error('Failed to fetch store stocks:', error)
			throw error
		}
	}
}

export const ingredientsService = new IngredientsService()
