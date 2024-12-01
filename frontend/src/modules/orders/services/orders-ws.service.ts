// websocketService.ts

import { useWebSocket, type UseWebSocketOptions } from '@vueuse/core'
import { onMounted, onUnmounted, ref } from 'vue'
import type { OrderDTO } from '../models/orders.models'

interface OrderDTO {}

export function useOrderWebSocket() {
	const orders = ref<OrderDTO[]>([])
	const wsUrl = 'ws://localhost:8080/ws'

	const options: UseWebSocketOptions = {
		autoReconnect: {
			retries: 3,
			delay: 3000,
			onFailed() {
				console.error('WebSocket reconnection failed')
			},
		},
		heartbeat: {
			message: JSON.stringify({ type: 'ping' }),
			interval: 10000,
		},
	}

	const { status, data, send, open, close } = useWebSocket(wsUrl, options)

	// Handle incoming messages
	const handleMessage = (event: MessageEvent) => {
		try {
			const message = JSON.parse(event.data)

			if (message.type === 'NEW_ORDER') {
				orders.value.unshift(message.order)
			} else if (message.type === 'UPDATE_ORDER') {
				const index = orders.value.findIndex(o => o.id === message.order.id)
				if (index !== -1) {
					orders.value[index] = message.order
				}
			} else if (message.type === 'REMOVE_ORDER') {
				orders.value = orders.value.filter(o => o.id !== message.order.id)
			}
		} catch (error) {
			console.error('Failed to parse WebSocket message:', error)
		}
	}

	onMounted(() => {
		// Listen for incoming messages
		if (data.value) {
			data.value.addEventListener('message', handleMessage)
		}
	})

	onUnmounted(() => {
		// Clean up
		if (data.value) {
			data.value.removeEventListener('message', handleMessage)
		}
		close()
	})

	return {
		orders,
		status,
		send,
		open,
		close,
	}
}
