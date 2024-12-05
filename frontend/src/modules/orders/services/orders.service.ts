import { apiClient } from '@/core/config/axios-instance.config'
import type {
	CreateOrderDTO,
	OrderDTO,
	OrderProductDTO,
	OrderStatus,
} from '../models/orders.models'

class OrderService {
	async getAllOrders(status?: OrderStatus): Promise<OrderDTO[]> {
		try {
			const response = await apiClient.get<OrderDTO[]>('/orders', {
				params: { status },
			})
			return response.data
		} catch (error) {
			console.error('Failed to fetch orders:', error)
			throw error
		}
	}

	async getSubOrders(orderId: number): Promise<OrderProductDTO[]> {
		try {
			const response = await apiClient.get<OrderProductDTO[]>('/orders/suborders', {
				params: { order_id: orderId },
			})
			return response.data
		} catch (error) {
			console.error('Failed to fetch sub-orders:', error)
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

	async completeSubOrder(subOrderId: number): Promise<void> {
		try {
			await apiClient.post(`/orders/suborders/${subOrderId}/complete`)
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

	async getStatusesCount(): Promise<{ [key in OrderStatus]?: number }> {
		try {
			const response =
				await apiClient.get<{ [key in OrderStatus]?: number }>('/orders/statuses/count')
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
