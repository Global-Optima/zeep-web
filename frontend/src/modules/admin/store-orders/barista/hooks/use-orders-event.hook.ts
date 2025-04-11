import { useToast } from '@/core/components/ui/toast'
import { getSavedBaristaQRSettings } from '@/core/hooks/use-qr-print.hook'
import type { OrderDTO, OrderStatus } from '@/modules/admin/store-orders/models/orders.models'
import { useWebSocket } from '@vueuse/core'
import { computed, ref } from 'vue'
import { useOrderQRPrinter } from '../../hooks/use-order-qr-print.hook'

interface OrderFilterOptions {
	status?: OrderStatus
	timeGapMinutes?: number
}

interface UseOrderEventsServiceOptions {
	printOnCreate?: boolean
}

// Use a ref to store orders instead of a reactive wrapper
const orders = ref<OrderDTO[]>([])

function mergeOrder(existingOrder: OrderDTO, updatedOrder: OrderDTO) {
	// Merge top-level properties
	Object.assign(existingOrder, updatedOrder)
}

function sortOrders(orderList: OrderDTO[]): OrderDTO[] {
	// Return a new sorted copy – ensures we trigger reactivity only on actual order changes.
	return orderList.slice().sort((a, b) => {
		return new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime()
	})
}

// -------------------------------------------------
// Main Composable: useOrderEvents
// -------------------------------------------------
export function useOrderEvents(
	initialFilter: OrderFilterOptions = {},
	options: UseOrderEventsServiceOptions = { printOnCreate: false },
) {
	// Destructure with a default option; you can control whether to print on create.
	const { printOnCreate = false } = options

	// Toast and QR Printer setup
	const { toast } = useToast()
	const { printSubOrderQR } = useOrderQRPrinter()

	// Reactive filter – only a shallow reactive object (a ref) is needed here.
	const localFilter = ref<OrderFilterOptions>({ ...initialFilter })

	// WebSocket URL built once; note that we only use localFilter.value.timeGapMinutes once here.
	const wsUrl = import.meta.env.VITE_WS_URL || 'ws://example:8080/api/v1'
	const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone
	const url = `${wsUrl}/orders/ws?timezone=${timezone}&timeGapMinutes=${
		localFilter.value.timeGapMinutes ?? 60
	}&includeYesterdayOrders=true`

	const {
		status: socketStatus,
		send,
		open,
		close,
	} = useWebSocket(url, {
		immediate: true,
		autoReconnect: { retries: 10, delay: 2000 },
		onConnected() {
			console.log('[WS] Connected!')
		},
		onDisconnected() {
			toast({
				title: 'Вы отключились',
				description: 'Соединение с заказами было прервано.',
			})
		},
		onError(err) {
			console.error('[WS] Error:', err)
			toast({
				title: 'Ошибка соединения',
				description: 'Не удалось подключиться к серверу.',
				variant: 'destructive',
			})
		},
		onMessage(_, event) {
			if (!event.data) return
			try {
				const msg = JSON.parse(event.data)
				switch (msg.type) {
					case 'initial_data':
						handleInitialData(msg.payload)
						break
					case 'order_succeeded':
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
				console.error('Не удалось разобрать сообщение WebSocket:', err)
				toast({
					title: 'Ошибка обработки данных',
					description: 'Не удалось обработать данные от сервера.',
					variant: 'destructive',
				})
			}
		},
	})

	// Computed properties to filter and count orders.
	const filteredOrders = computed(() => {
		if (!localFilter.value.status) return orders.value
		return orders.value.filter(o => o.status === localFilter.value.status)
	})

	const orderCountsByStatus = computed<Record<OrderStatus, number>>(() => {
		const counts: Record<OrderStatus, number> = {
			WAITING_FOR_PAYMENT: 0,
			PENDING: 0,
			PREPARING: 0,
			COMPLETED: 0,
			IN_DELIVERY: 0,
			DELIVERED: 0,
			CANCELLED: 0,
		}
		orders.value.forEach(order => {
			counts[order.status] += 1
		})
		return counts
	})

	// Public method to update the filter.
	function setFilter(newFilter: OrderFilterOptions) {
		localFilter.value = { ...localFilter.value, ...newFilter }
	}

	// -------------------------------------------------
	// Event Handlers
	// -------------------------------------------------
	function handleInitialData(initialOrders: OrderDTO[]) {
		orders.value = sortOrders(initialOrders)
	}

	async function handleOrderCreated(newOrder: OrderDTO) {
		const [reactiveOrder, isNew] = upsertOrder(newOrder)
		if (!isNew || !printOnCreate) return

		const { width, height } = getSavedBaristaQRSettings()
		try {
			await printSubOrderQR(reactiveOrder, reactiveOrder.subOrders, {
				labelWidthMm: width,
				labelHeightMm: height,
				desktopOnly: false,
			})
		} catch (err) {
			console.error('Ошибка при печати QR-кода:', err)
			toast({
				title: 'Ошибка печати',
				description: 'Не удалось напечатать один или несколько QR-кодов.',
				variant: 'destructive',
			})
		}
	}

	function handleOrderUpdated(updatedOrder: OrderDTO) {
		upsertOrder(updatedOrder)
	}

	function handleOrderDeleted(orderId: number) {
		orders.value = orders.value.filter(o => o.id !== orderId)
	}

	// -------------------------------------------------
	// Upsert Helper: Inserts or updates an order.
	// -------------------------------------------------
	function upsertOrder(orderData: OrderDTO): [OrderDTO, boolean] {
		const idx = orders.value.findIndex(o => o.id === orderData.id)
		if (idx !== -1) {
			mergeOrder(orders.value[idx], orderData)
			orders.value = sortOrders(orders.value)
			return [orders.value[idx], false]
		} else {
			orders.value.unshift(orderData)
			orders.value = sortOrders(orders.value)
			return [orderData, true]
		}
	}

	// -------------------------------------------------
	// Return a lean public API.
	// -------------------------------------------------
	return {
		status: socketStatus,
		filteredOrders,
		allOrders: orders,
		orderCountsByStatus,
		setFilter,
		send,
		open,
		close,
	}
}
