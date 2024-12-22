import { apiClient } from '@/core/config/axios-instance.config'
import type { PaginatedResponse } from '@/core/utils/pagination.utils'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type {
	CreateStoreStockRequestDTO,
	CreateStoreStockRequestItemDTO,
	GetStoreStockRequestsFilter,
	StoreStockRequestResponse,
	UpdateStoreStockRequestStatusDTO,
} from '../models/store-stock-request.model'

class StoreStockRequestService {
	async getStockRequests(filter?: GetStoreStockRequestsFilter) {
		try {
			const response = await apiClient.get<PaginatedResponse<StoreStockRequestResponse[]>>(
				'/stock-requests',
				{
					params: buildRequestFilter(filter),
				},
			)
			return response.data
		} catch (error) {
			console.error('Failed to fetch stock requests:', error)
			throw error
		}
	}

	async getStockRequestById(id: number) {
		try {
			const response = await apiClient.get<StoreStockRequestResponse>(`/stock-requests/${id}`)
			return response.data
		} catch (error) {
			console.error('Failed to fetch stock request by id:', error)
			throw error
		}
	}

	async getLowStockIngredients() {
		try {
			const response = await apiClient.get('/stock-requests/low-stock')
			return response.data
		} catch (error) {
			console.error('Failed to fetch low-stock ingredients:', error)
			throw error
		}
	}

	async getMarketplaceProducts() {
		try {
			const response = await apiClient.get('/stock-requests/marketplace-products')
			return response.data
		} catch (error) {
			console.error('Failed to fetch marketplace products:', error)
			throw error
		}
	}

	async createStockRequest(data: CreateStoreStockRequestDTO) {
		try {
			const response = await apiClient.post('/stock-requests', data)
			return response.data
		} catch (error) {
			console.error('Failed to create stock request:', error)
			throw error
		}
	}

	async updateStockRequestStatus(requestId: number, data: UpdateStoreStockRequestStatusDTO) {
		try {
			const response = await apiClient.put(`/stock-requests/${requestId}/status`, data)
			return response.data
		} catch (error) {
			console.error('Failed to update stock request status:', error)
			throw error
		}
	}

	async addStockRequestIngredient(requestId: number, item: CreateStoreStockRequestItemDTO) {
		try {
			const response = await apiClient.post(`/stock-requests/${requestId}/ingredients`, item)
			return response.data
		} catch (error) {
			console.error('Failed to add ingredient to stock request:', error)
			throw error
		}
	}

	async deleteStockRequestIngredient(ingredientId: number) {
		try {
			const response = await apiClient.delete(`/stock-requests/ingredients/${ingredientId}`)
			return response.data
		} catch (error) {
			console.error('Failed to delete stock request ingredient:', error)
			throw error
		}
	}
}

export const storeStockRequestService = new StoreStockRequestService()
