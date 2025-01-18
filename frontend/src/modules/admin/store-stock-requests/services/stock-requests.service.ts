import { apiClient } from '@/core/config/axios-instance.config'
import type { PaginatedResponse } from '@/core/utils/pagination.utils'
import type {
	AcceptWithChangeRequestStatusDTO,
	CreateStockRequestDTO,
	GetStockRequestsFilter,
	LowStockIngredientResponse,
	RejectStockRequestStatusDTO,
	StockMaterialAvailabilityDTO,
	StockMaterialDTO,
	StockMaterialFilter,
	StockRequestResponse,
	StockRequestStockMaterialDTO,
} from '../models/stock-requests.model'

class StockRequestsService {
	private readonly baseUrl: string = '/stock-requests'

	async getStockRequests(filter?: GetStockRequestsFilter) {
		const response = await apiClient.get<PaginatedResponse<StockRequestResponse[]>>(this.baseUrl, {
			params: filter,
		})
		return response.data
	}

	async getStockRequestById(requestId: number) {
		const response = await apiClient.get<StockRequestResponse>(`${this.baseUrl}/${requestId}`)
		return response.data
	}

	async createStockRequest(data: CreateStockRequestDTO) {
		const response = await apiClient.post(this.baseUrl, data)
		return response.data
	}

	async getLastCreatedStockRequest() {
		const response = await apiClient.get<StockRequestResponse>(`${this.baseUrl}/current`)
		return response.data
	}

	async updateStockRequestMaterials(requestId: number, items: StockRequestStockMaterialDTO[]) {
		const response = await apiClient.put(`${this.baseUrl}/${requestId}`, items)
		return response.data
	}

	async deleteStockRequest(requestId: number) {
		const response = await apiClient.delete(`${this.baseUrl}/${requestId}`)
		return response.data
	}

	async acceptWithChange(requestId: number, data: AcceptWithChangeRequestStatusDTO) {
		const response = await apiClient.patch(
			`${this.baseUrl}/status/${requestId}/accept-with-change`,
			data,
		)
		return response.data
	}

	async rejectStore(requestId: number, data: RejectStockRequestStatusDTO) {
		const response = await apiClient.patch(`${this.baseUrl}/status/${requestId}/reject-store`, data)
		return response.data
	}

	async rejectWarehouse(requestId: number, data: RejectStockRequestStatusDTO) {
		const response = await apiClient.patch(
			`${this.baseUrl}/status/${requestId}/reject-warehouse`,
			data,
		)
		return response.data
	}

	async setInDeliveryStatus(requestId: number) {
		const response = await apiClient.patch(`${this.baseUrl}/status/${requestId}/in-delivery`)
		return response.data
	}

	async setCompletedStatus(requestId: number) {
		const response = await apiClient.patch(`${this.baseUrl}/status/${requestId}/completed`)
		return response.data
	}

	async getLowStockIngredients() {
		const response = await apiClient.get<LowStockIngredientResponse[]>(`${this.baseUrl}/low-stock`)
		return response.data
	}

	async getAllStockMaterials(filter?: StockMaterialFilter) {
		const response = await apiClient.get<PaginatedResponse<StockMaterialDTO[]>>(
			`${this.baseUrl}/materials`,
			{
				params: filter,
			},
		)
		return response.data
	}

	async getStockMaterialsByIngredient(ingredientId: number) {
		const response = await apiClient.get<StockMaterialAvailabilityDTO>(
			`${this.baseUrl}/materials/${ingredientId}`,
		)
		return response.data
	}
}

export const stockRequestsService = new StockRequestsService()
