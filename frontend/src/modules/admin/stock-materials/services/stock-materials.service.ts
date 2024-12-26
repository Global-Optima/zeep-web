import { apiClient } from '@/core/config/axios-instance.config'
import type { PaginatedResponse } from '@/core/utils/pagination.utils'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type {
	CreateStockMaterialDTO,
	StockMaterialsDTO,
	StockMaterialsFilter,
	UpdateStockMaterialDTO,
} from '@/modules/admin/stock-materials/models/stock-materials.model'

class StockMaterialsService {
	private readonly baseUrl: string = '/stock-materials'

	// Get all stock materials with filtering and pagination
	async getAllStockMaterials(filter?: StockMaterialsFilter) {
		try {
			const response = await apiClient.get<PaginatedResponse<StockMaterialsDTO[]>>(this.baseUrl, {
				params: buildRequestFilter(filter),
			})
			return response.data
		} catch (error) {
			console.error('Failed to fetch stock materials:', error)
			throw error
		}
	}

	// Get a stock material by ID
	async getStockMaterialById(id: number) {
		try {
			const response = await apiClient.get<StockMaterialsDTO>(`${this.baseUrl}/${id}`)
			return response.data
		} catch (error) {
			console.error(`Failed to fetch stock material with ID ${id}:`, error)
			throw error
		}
	}

	// Create a new stock material
	async createStockMaterial(data: CreateStockMaterialDTO) {
		try {
			const response = await apiClient.post(this.baseUrl, data)
			return response.data
		} catch (error) {
			console.error('Failed to create stock material:', error)
			throw error
		}
	}

	// Update a stock material
	async updateStockMaterial(id: number, data: UpdateStockMaterialDTO) {
		try {
			const response = await apiClient.put(`${this.baseUrl}/${id}`, data)
			return response.data
		} catch (error) {
			console.error(`Failed to update stock material with ID ${id}:`, error)
			throw error
		}
	}

	// Delete a stock material
	async deleteStockMaterial(id: number) {
		try {
			const response = await apiClient.delete(`${this.baseUrl}/${id}`)
			return response.data
		} catch (error) {
			console.error(`Failed to delete stock material with ID ${id}:`, error)
			throw error
		}
	}

	// Deactivate a stock material
	async deactivateStockMaterial(id: number) {
		try {
			const response = await apiClient.patch(`${this.baseUrl}/${id}/deactivate`)
			return response.data
		} catch (error) {
			console.error(`Failed to deactivate stock material with ID ${id}:`, error)
			throw error
		}
	}
}

export const stockMaterialsService = new StockMaterialsService()
