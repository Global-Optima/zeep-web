import { apiClient } from '@/core/config/axios-instance.config'
import type { PaginatedResponse } from '@/core/utils/pagination.utils'
import type {
	AcceptWithChangeRequestStatusDTO,
	CreateStockRequestDTO,
	GetStockRequestsFilter,
	RejectStockRequestStatusDTO,
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

	async addStockMaterialToLatestCart(item: StockRequestStockMaterialDTO) {
		const response = await apiClient.post(`${this.baseUrl}/add-material-to-latest-cart`, item)
		return response.data
	}

	async deleteStockRequest(requestId: number) {
		const response = await apiClient.delete(`${this.baseUrl}/${requestId}`)
		return response.data
	}

	async acceptWithChanges(requestId: number, data: AcceptWithChangeRequestStatusDTO) {
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

	async setProcessedStatus(requestId: number) {
		const response = await apiClient.patch(`${this.baseUrl}/status/${requestId}/processed`)
		return response.data
	}
}

export const stockRequestsService = new StockRequestsService()
