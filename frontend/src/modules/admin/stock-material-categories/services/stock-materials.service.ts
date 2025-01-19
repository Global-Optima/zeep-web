import { apiClient } from '@/core/config/axios-instance.config'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type {
	CreateStockMaterialCategoryDTO,
	StockMaterialCategoryDTO,
	StockMaterialCategoryFilterDTO,
	UpdateStockMaterialCategoryDTO,
} from '@/modules/admin/stock-material-categories/models/stock-material-categories.model'
import { type PaginatedResponse } from './../../../../core/utils/pagination.utils'

class StockMaterialCategoryService {
	private readonly baseUrl: string = '/stock-material-categories'

	async getAll(filter?: StockMaterialCategoryFilterDTO) {
		try {
			const response = await apiClient.get<PaginatedResponse<StockMaterialCategoryDTO[]>>(
				this.baseUrl,
				{ params: buildRequestFilter(filter) },
			)
			return response.data
		} catch (error) {
			console.error('Failed to fetch stock material categories:', error)
			throw error
		}
	}

	async getById(id: number): Promise<StockMaterialCategoryDTO> {
		try {
			const response = await apiClient.get<StockMaterialCategoryDTO>(`${this.baseUrl}/${id}`)
			return response.data
		} catch (error) {
			console.error(`Failed to fetch stock material category with ID ${id}:`, error)
			throw error
		}
	}

	async create(data: CreateStockMaterialCategoryDTO): Promise<StockMaterialCategoryDTO> {
		try {
			const response = await apiClient.post<StockMaterialCategoryDTO>(this.baseUrl, data)
			return response.data
		} catch (error) {
			console.error('Failed to create stock material category:', error)
			throw error
		}
	}

	async update(
		id: number,
		data: UpdateStockMaterialCategoryDTO,
	): Promise<StockMaterialCategoryDTO> {
		try {
			const response = await apiClient.put<StockMaterialCategoryDTO>(`${this.baseUrl}/${id}`, data)
			return response.data
		} catch (error) {
			console.error(`Failed to update stock material category with ID ${id}:`, error)
			throw error
		}
	}

	async delete(id: number): Promise<void> {
		try {
			await apiClient.delete(`${this.baseUrl}/${id}`)
		} catch (error) {
			console.error(`Failed to delete stock material category with ID ${id}:`, error)
			throw error
		}
	}
}

export const stockMaterialCategoryService = new StockMaterialCategoryService()
