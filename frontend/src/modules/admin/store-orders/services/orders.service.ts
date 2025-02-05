import { apiClient } from '@/core/config/axios-instance.config'
import type { PaginatedResponse } from '@/core/utils/pagination.utils'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import { saveAs } from 'file-saver'
import type {
	CreateOrderDTO,
	OrderDTO,
	OrdersExportFilterQuery,
	OrdersFilterQuery,
	OrderStatusesCountDTO,
	SuborderDTO,
} from '../models/orders.models'

class OrderService {
	async getAllOrders(filter?: OrdersFilterQuery) {
		try {
			const response = await apiClient.get<PaginatedResponse<OrderDTO[]>>('/orders', {
				params: buildRequestFilter(filter),
			})
			return response.data
		} catch (error) {
			console.error('Failed to fetch orders:', error)
			throw error
		}
	}

	async createOrder(orderDTO: CreateOrderDTO) {
		try {
			return apiClient.post<{ orderId: number }>('/orders', orderDTO).then(res => res.data)
		} catch (error) {
			console.error('Failed to create order:', error)
			throw error
		}
	}

	async getSuborderBarcode(suborderId: number) {
		try {
			const response = await apiClient.get<Blob>(`/orders/suborders/${suborderId}/barcode`, {
				responseType: 'blob',
			})
			return response.data
		} catch (error) {
			console.error(`Failed to fetch barcode for suborder ID ${suborderId}:`, error)
			throw error
		}
	}

	async getSuborderBarcodeFile(suborderId: number) {
		try {
			const response = await apiClient.get<Blob>(`/orders/suborders/${suborderId}/barcode`, {
				responseType: 'blob',
			})

			return response.data
		} catch (error) {
			console.error(`Failed to fetch barcode for stock material ID ${suborderId}:`, error)
			throw error
		}
	}

	async completeSubOrder(orderId: number, subOrderId: number): Promise<void> {
		try {
			await apiClient.put(`/orders/${orderId}/suborders/${subOrderId}/complete`, {})
		} catch (error) {
			console.error('Failed to complete sub-order:', error)
			throw error
		}
	}

	async completeSubOrderById(subOrderId: number) {
		try {
			const response = await apiClient.put<SuborderDTO>(
				`/orders/suborders/${subOrderId}/complete`,
				{},
			)
			return response.data
		} catch (error) {
			console.error('Failed to complete sub-order:', error)
			throw error
		}
	}

	async revertSubOrderById(subOrderId: number): Promise<void> {
		try {
			await apiClient.put(`/orders/suborders/${subOrderId}/revert`, {})
		} catch (error) {
			console.error('Failed to revert sub-order:', error)
			throw error
		}
	}

	async generatePDFReceipt(orderId: number): Promise<Blob> {
		try {
			const response = await apiClient.get(`/orders/${orderId}/receipt`, {
				responseType: 'blob',
			})
			return response.data
		} catch (error) {
			console.error('Failed to generate PDF receipt:', error)
			throw error
		}
	}

	async getStatusesCount(): Promise<OrderStatusesCountDTO> {
		try {
			const response = await apiClient.get<OrderStatusesCountDTO>('/orders/statuses/count', {})
			return response.data
		} catch (error) {
			console.error('Failed to fetch statuses count:', error)
			throw error
		}
	}

	async getSubOrderCount(orderId: number): Promise<number> {
		try {
			const response = await apiClient.get<{ count: number }>('/orders/suborders/count', {
				params: { order_id: orderId },
			})
			return response.data.count
		} catch (error) {
			console.error('Failed to fetch sub-order count:', error)
			throw error
		}
	}

	async exportOrders(filter?: OrdersExportFilterQuery) {
		try {
			const response = await apiClient.get('/orders/export', {
				params: buildRequestFilter(filter),
				responseType: 'blob',
			})

			saveAs(response.data)
		} catch (error) {
			console.error('Failed to export orders:', error)
			throw error
		}
	}
}

export const ordersService = new OrderService()
