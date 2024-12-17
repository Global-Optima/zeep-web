import { apiClient } from '@/core/config/axios-instance.config'
import type { PaginatedResponse } from '@/core/utils/pagination.utils'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type {
	CreateOrderDTO,
	OrderDTO,
	OrdersFilterQuery,
	OrderStatusesCountDTO,
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

	async completeSubOrder(storeId: number, orderId: number, subOrderId: number): Promise<void> {
		try {
			await apiClient.put(
				`/orders/${orderId}/suborders/${subOrderId}/complete`,
				{},
				{ params: { storeId } },
			)
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

	async getStatusesCount(storeId: number): Promise<OrderStatusesCountDTO> {
		try {
			const response = await apiClient.get<OrderStatusesCountDTO>('/orders/statuses/count', {
				params: { storeId },
			})
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
}

export const orderService = new OrderService()
