import { apiClient } from '@/core/config/axios-instance.config'
import type { PaginatedResponse } from '@/core/utils/pagination.utils'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type {
	CreateStockMaterialDTO,
	GeneratedStockMaterialBarcode,
	StockMaterialsDTO,
	StockMaterialsFilter,
	UpdateStockMaterialDTO,
} from '../models/stock-materials.model'

class StockMaterialService {
	private readonly baseUrl: string = '/stock-materials'

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

	async getBarcodeFile(stockMaterialId: number) {
		try {
			const response = await apiClient.get<Blob>(`/stock-materials/${stockMaterialId}/barcode`, {
				responseType: 'blob', // Ensure the response is treated as a Blob
			})
			return response.data
		} catch (error) {
			console.error(`Failed to fetch barcode for stock material ID ${stockMaterialId}:`, error)
			throw error
		}
	}

	async generateBarcode() {
		try {
			const response = await apiClient.get<GeneratedStockMaterialBarcode>(`${this.baseUrl}/barcode`)
			return response.data
		} catch (error) {
			console.error(`Failed to create barcode for stock material:`, error)
			throw error
		}
	}

	async getStockMaterialById(id: number) {
		try {
			const response = await apiClient.get<StockMaterialsDTO>(`${this.baseUrl}/${id}`)
			return response.data
		} catch (error) {
			console.error(`Failed to fetch stock material with ID ${id}:`, error)
			throw error
		}
	}

	async createStockMaterial(data: CreateStockMaterialDTO) {
		try {
			const response = await apiClient.post(this.baseUrl, data)
			return response.data
		} catch (error) {
			console.error('Failed to create stock material:', error)
			throw error
		}
	}

	async updateStockMaterial(id: number, data: UpdateStockMaterialDTO) {
		try {
			const response = await apiClient.put(`${this.baseUrl}/${id}`, data)
			return response.data
		} catch (error) {
			console.error(`Failed to update stock material with ID ${id}:`, error)
			throw error
		}
	}

	async deleteStockMaterial(id: number) {
		try {
			const response = await apiClient.delete(`${this.baseUrl}/${id}`)
			return response.data
		} catch (error) {
			console.error(`Failed to delete stock material with ID ${id}:`, error)
			throw error
		}
	}

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

export const stockMaterialsService = new StockMaterialService()
