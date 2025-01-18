import type { OrderDTO } from '@/modules/orders/models/orders.models'
import { useQueryClient } from '@tanstack/vue-query'
import { useWebSocket } from '@vueuse/core'
import { ref, watchEffect } from 'vue'

const orders = ref<OrderDTO[]>([])

const wsUrl = `${import.meta.env.VITE_WS_URL || 'ws://localhost:8080/api/v1'}`

export function useOrderEventsService() {
	const url = `${wsUrl}/orders/ws`
	const queryClient = useQueryClient()

	const {
		status,
		data: messageEvent,
		send,
		open,
		close,
	} = useWebSocket(url, {
		immediate: true,
		autoReconnect: true,
		onError: err => console.error('WebSocket error:', err),
	})

	// Watch for WebSocket events
	watchEffect(() => {
		if (!messageEvent.value) return

		try {
			const msg = JSON.parse(messageEvent.value)
			console.log('EVVEEENT', msg.payload)

			switch (msg.type) {
				case 'initial_data':
					handleInitialData(msg.payload)
					break
				case 'order_created':
					handleOrderCreated(msg.payload)
					break
				case 'order_updated':
					handleOrderUpdated(msg.payload)
					break
				case 'order_deleted':
					handleOrderDeleted(msg.payload)
					break
				default:
					console.warn('Unknown event type:', msg.type)
			}
		} catch (err) {
			console.error('Failed to parse incoming message:', err)
		}
	})

	function invalidateOrderStatuses() {
		queryClient.invalidateQueries({ queryKey: ['statuses'] })
	}

	// Handle initial data
	function handleInitialData(initialOrders: OrderDTO[]) {
		orders.value = sortOrders(initialOrders)
		invalidateOrderStatuses()
	}

	// Handle order created event
	function handleOrderCreated(order: OrderDTO) {
		orders.value.unshift(order)
		orders.value = sortOrders(orders.value)
		invalidateOrderStatuses()
	}

	// Handle order updated event
	function handleOrderUpdated(order: OrderDTO) {
		const idx = orders.value.findIndex(o => o.id === order.id)
		if (idx !== -1) {
			orders.value[idx] = order
		} else {
			orders.value.unshift(order)
		}
		orders.value = sortOrders(orders.value)
		invalidateOrderStatuses()
	}

	// Handle order deleted event
	function handleOrderDeleted(orderId: number) {
		orders.value = orders.value.filter(o => o.id !== orderId)
		invalidateOrderStatuses()
	}

	// Sort orders by timestamp
	function sortOrders(orderList: OrderDTO[]): OrderDTO[] {
		return orderList.sort(
			(a, b) => new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime(),
		)
	}

	// Expose reactive orders reference
	function getOrdersRef() {
		return orders
	}

	return {
		status,
		getOrdersRef,
		close,
		open,
		send,
	}
}
