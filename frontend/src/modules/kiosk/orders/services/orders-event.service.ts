import type { OrderDTO, OrderStatus } from '@/modules/orders/models/orders.models'
import { useQueryClient } from '@tanstack/vue-query'
import { useWebSocket } from '@vueuse/core'
import { computed, ref, watchEffect } from 'vue'

interface OrderFilterOptions {
	status?: OrderStatus
}

const allOrders = ref<OrderDTO[]>([])

const wsUrl = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/api/v1'

export function useOrderEventsService(filter: OrderFilterOptions = {}) {
	const queryClient = useQueryClient()
	const url = `${wsUrl}/orders/ws`
	const localFilter = ref({ ...filter })

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

	watchEffect(() => {
		if (!messageEvent.value) return

		try {
			const msg = JSON.parse(messageEvent.value)
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

	const filteredOrders = computed(() => {
		if (!localFilter.value.status) return allOrders.value
		return allOrders.value.filter(order => order.status === localFilter.value.status)
	})

	function setFilter(filter: OrderFilterOptions) {
		localFilter.value = { ...localFilter.value, ...filter }
	}

	function invalidateOrderStatuses() {
		queryClient.invalidateQueries({ queryKey: ['order-statuses'] })
	}

	function handleInitialData(initialOrders: OrderDTO[]) {
		allOrders.value = sortOrders(initialOrders)
		invalidateOrderStatuses()
	}

	function handleOrderCreated(order: OrderDTO) {
		allOrders.value.unshift(order)
		allOrders.value = sortOrders(allOrders.value)
		invalidateOrderStatuses()
	}

	function handleOrderUpdated(order: OrderDTO) {
		const idx = allOrders.value.findIndex(o => o.id === order.id)
		if (idx !== -1) {
			allOrders.value[idx] = order
		} else {
			allOrders.value.unshift(order)
		}
		allOrders.value = sortOrders(allOrders.value)
		invalidateOrderStatuses()
	}

	function handleOrderDeleted(orderId: number) {
		allOrders.value = allOrders.value.filter(o => o.id !== orderId)
		invalidateOrderStatuses()
	}

	function sortOrders(orderList: OrderDTO[]): OrderDTO[] {
		return orderList.sort(
			(a, b) => new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime(),
		)
	}

	return {
		status,
		filteredOrders,
		setFilter,
		send,
		open,
		close,
	}
}
