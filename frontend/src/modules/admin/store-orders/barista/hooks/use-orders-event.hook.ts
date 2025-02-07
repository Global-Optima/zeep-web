import type { OrderDTO, OrderStatus } from '@/modules/admin/store-orders/models/orders.models'
import { useQueryClient } from '@tanstack/vue-query'
import { useWebSocket } from '@vueuse/core'
import { computed, reactive, ref, watchEffect } from 'vue'

interface OrderFilterOptions {
	status?: OrderStatus
}

const state = reactive<{
	allOrders: OrderDTO[]
}>({
	allOrders: [],
})

const wsUrl = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/api/v1'

export function useOrderEventsService(filter: OrderFilterOptions = {}) {
	const queryClient = useQueryClient()
	const url = `${wsUrl}/orders/ws`

	const localFilter = ref({ ...filter })

	const {
		status: socketStatus,
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
		if (!localFilter.value.status) return state.allOrders
		return state.allOrders.filter(order => order.status === localFilter.value.status)
	})

	function setFilter(newFilter: OrderFilterOptions) {
		localFilter.value = { ...localFilter.value, ...newFilter }
	}

	function invalidateOrderStatuses() {
		queryClient.invalidateQueries({ queryKey: ['order-statuses'] })
	}

	function handleInitialData(initialOrders: OrderDTO[]) {
		state.allOrders = sortOrders(initialOrders)
		invalidateOrderStatuses()
	}

	function handleOrderCreated(newOrder: OrderDTO) {
		state.allOrders.unshift(newOrder)
		state.allOrders = sortOrders([...state.allOrders])
		invalidateOrderStatuses()
	}

	function handleOrderUpdated(updated: OrderDTO) {
		const idx = state.allOrders.findIndex(o => o.id === updated.id)
		if (idx !== -1) {
			state.allOrders[idx] = { ...state.allOrders[idx], ...updated }
		} else {
			state.allOrders.unshift(updated)
		}
		state.allOrders = sortOrders([...state.allOrders])
		invalidateOrderStatuses()
	}

	function handleOrderDeleted(orderId: number) {
		state.allOrders = state.allOrders.filter(o => o.id !== orderId)
		invalidateOrderStatuses()
	}

	function sortOrders(orderList: OrderDTO[]): OrderDTO[] {
		return orderList.sort(
			(a, b) => new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime(),
		)
	}

	return {
		status: socketStatus,
		filteredOrders,
		setFilter,
		send,
		open,
		close,
	}
}
