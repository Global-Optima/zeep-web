import { useToast } from '@/core/components/ui/toast'
import { useBarcodePrinter } from '@/core/hooks/use-barcode-print.hook'
import type { OrderDTO, OrderStatus } from '@/modules/admin/store-orders/models/orders.models'
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
	const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone

	const { printBarcode } = useBarcodePrinter()
	const { toast } = useToast()
	const url = `${wsUrl}/orders/ws?timezone=${timezone}`

	const localFilter = ref({ ...filter })

	// WebSocket connection setup
	const {
		status: socketStatus,
		data: messageEvent,
		send,
		open,
		close,
	} = useWebSocket(url, {
		immediate: true,
		autoReconnect: true,
		onError: err => {
			console.error('Ошибка WebSocket:', err)
			toast({
				title: 'Ошибка соединения',
				description: 'Не удалось подключиться к серверу.',
				variant: 'destructive',
			})
		},
		onDisconnected: () => {
			toast({
				title: 'Соединение прервано.',
				description: 'Соединение с сервером прервано.',
			})
		},
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
					console.warn('Неизвестный тип события:', msg.type)
					toast({
						title: 'Неизвестное событие',
						description: `Получен неизвестный тип события: ${msg.type}`,
						variant: 'destructive',
					})
			}
		} catch (err) {
			console.error('Не удалось разобрать входящее сообщение:', err)
			toast({
				title: 'Ошибка обработки данных',
				description: 'Не удалось обработать данные от сервера.',
				variant: 'destructive',
			})
		}
	})

	// Computed property for filtered orders
	const filteredOrders = computed(() => {
		if (!localFilter.value.status) return state.allOrders
		return state.allOrders.filter(order => order.status === localFilter.value.status)
	})

	// Reactive computation of order counts by status
	const orderCountsByStatus = computed<Record<OrderStatus, number>>(() => {
		const counts: Record<OrderStatus, number> = {
			PENDING: 0,
			PREPARING: 0,
			COMPLETED: 0,
			IN_DELIVERY: 0,
			DELIVERED: 0,
			CANCELLED: 0,
		}

		state.allOrders.forEach(order => {
			if (counts[order.status] !== undefined) {
				counts[order.status]++
			}
		})

		return counts
	})

	// Set filter options
	function setFilter(newFilter: OrderFilterOptions) {
		localFilter.value = { ...localFilter.value, ...newFilter }
	}

	// Handle initial data from WebSocket
	function handleInitialData(initialOrders: OrderDTO[]) {
		state.allOrders = sortOrders(initialOrders)
	}

	// Handle newly created orders
	async function handleOrderCreated(newOrder: OrderDTO) {
		state.allOrders.unshift(newOrder)
		state.allOrders = sortOrders([...state.allOrders])

		try {
			await Promise.all(
				newOrder.subOrders.map(async subOrder => {
					const productName = `${subOrder.productSize.productName} ${subOrder.productSize.sizeName}`
					const barcode = `suborder-${subOrder.id}`
					await printBarcode(productName, barcode, { showModal: false })

					toast({
						title: 'Штрих-код напечатан',
						description: `Штрих-код для "${productName}" успешно напечатан.`,
					})
				}),
			)
		} catch (error) {
			console.error('Ошибка при печати штрих-кода:', error)
			toast({
				title: 'Ошибка печати',
				description: 'Не удалось напечатать один или несколько штрих-кодов.',
				variant: 'destructive',
			})
		}
	}

	// Handle updated orders
	function handleOrderUpdated(updated: OrderDTO) {
		const idx = state.allOrders.findIndex(o => o.id === updated.id)
		if (idx !== -1) {
			state.allOrders[idx] = { ...state.allOrders[idx], ...updated }
		} else {
			state.allOrders.unshift(updated)
		}
		state.allOrders = sortOrders([...state.allOrders])
		toast({
			title: 'Заказ обновлен',
			description: `Заказ с ${updated.id} успешно обновлен.`,
		})
	}

	// Handle deleted orders
	function handleOrderDeleted(orderId: number) {
		state.allOrders = state.allOrders.filter(o => o.id !== orderId)
		toast({
			title: 'Заказ удален',
			description: `Заказ с ${orderId} успешно удален.`,
		})
	}

	// Sort orders by creation time
	function sortOrders(orderList: OrderDTO[]): OrderDTO[] {
		return orderList.sort(
			(a, b) => new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime(),
		)
	}

	return {
		status: socketStatus,
		filteredOrders,
		orderCountsByStatus,
		setFilter,
		send,
		open,
		close,
	}
}
