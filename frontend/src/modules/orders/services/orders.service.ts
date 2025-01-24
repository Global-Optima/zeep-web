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
} from '../models/orders.models'
import { format } from 'date-fns'

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

	async completeSubOrder(orderId: number, subOrderId: number): Promise<void> {
		try {
			await apiClient.put(`/orders/${orderId}/suborders/${subOrderId}/complete`, {})
		} catch (error) {
			console.error('Failed to complete sub-order:', error)
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

			// Extract filename from Content-Disposition header
			const contentDisposition = response.headers['Content-Disposition']
			let filename = `orders_export_${format(new Date(), "dd_MM_yyyy")}.xlsx`

			if (contentDisposition) {
				const filenameMatch = contentDisposition.match(/filename="?(.+)"?/)
				if (filenameMatch && filenameMatch[1]) {
					filename = filenameMatch[1]
				}
			}

			const blob = new Blob([response.data], { type: response.headers['Content-Type']?.toString() })

			saveAs(blob, filename)
		} catch (error) {
			console.error('Failed to export orders:', error)
			throw error
		}
	}
}

export const ordersService = new OrderService()
