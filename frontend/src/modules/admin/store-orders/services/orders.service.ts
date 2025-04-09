import { apiClient } from '@/core/config/axios-instance.config'
import { encryptPayload, type EncryptedData } from '@/core/services/aes.service'
import type { PaginatedResponse } from '@/core/utils/pagination.utils'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import { saveAs } from 'file-saver'
import type {
	CheckCustomerName,
	CreateOrderDTO,
	OrderDetailsDTO,
	OrderDTO,
	OrdersExportFilterQuery,
	OrdersFilterQuery,
	OrdersTimeZoneFilter,
	SuborderDTO,
	ToggleNextSuborderStatusOptions,
	TransactionDTO,
} from '../models/orders.models'

class OrderService {
	private readonly paymentSecret = import.meta.env.VITE_PAYMENT_SECRET as string

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

	async getBaristaOrders(filter?: OrdersTimeZoneFilter) {
		try {
			const response = await apiClient.get<OrderDTO[]>('/orders/kiosk', {
				params: buildRequestFilter(filter),
			})
			return response.data
		} catch (error) {
			console.error('Failed to fetch barista orders:', error)
			throw error
		}
	}

	async getOrderById(id: number) {
		try {
			const response = await apiClient.get<OrderDetailsDTO>(`/orders/${id}`)
			return response.data
		} catch (error) {
			console.error('Failed to fetch order:', error)
			throw error
		}
	}

	async createOrder(orderDTO: CreateOrderDTO) {
		try {
			return apiClient.post<OrderDTO>('/orders', orderDTO).then(res => res.data)
		} catch (error) {
			console.error('Failed to create order:', error)
			throw error
		}
	}

	async failOrderPayment(orderId: number) {
		try {
			return apiClient
				.post<{ orderId: number }>(`/orders/${orderId}/payment/fail`)
				.then(res => res.data)
		} catch (error) {
			console.error('Failed to fail order payment:', error)
			throw error
		}
	}

	async successOrderPayment(orderId: number, dto: TransactionDTO): Promise<{ orderId: number }> {
		try {
			// Load secret key from environment variables
			if (!this.paymentSecret) {
				throw new Error('Payment is not defined in environment variables')
			}

			// Encrypt the TransactionDTO
			const encryptedData: EncryptedData = encryptPayload(dto, this.paymentSecret)

			// Send encrypted payload to the backend
			const res = await apiClient.post<{ orderId: number }>(
				`/orders/${orderId}/payment/success`,
				encryptedData,
			)

			return res.data
		} catch (error) {
			console.error('Failed to complete order payment:', error)
			throw error
		}
	}

	async checkCustomerName(orderDTO: CheckCustomerName) {
		try {
			return apiClient.post<void>('/orders/check-name', orderDTO).then(res => res.data)
		} catch (error) {
			console.error('Failed to validate customer name:', error)
			throw error
		}
	}

	async toggleNextSuborderStatus(subOrderId: number, options?: ToggleNextSuborderStatusOptions) {
		try {
			const response = await apiClient.put<SuborderDTO>(
				`/orders/suborders/${subOrderId}/status-change`,
				{},
				{ params: buildRequestFilter(options) },
			)
			return response.data
		} catch (error) {
			console.error('Failed to change sub-order status:', error)
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
