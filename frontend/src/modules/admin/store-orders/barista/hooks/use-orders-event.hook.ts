import { useToast } from '@/core/components/ui/toast'
import {
	generateSubOrderQR,
	getSavedBaristaQRSettings,
	useQRPrinter,
} from '@/core/hooks/use-qr-print.hook'
import type {
	OrderDTO,
	OrderStatus,
	SuborderDTO,
} from '@/modules/admin/store-orders/models/orders.models'

import { useWebSocket } from '@vueuse/core'
import { computed, reactive, ref } from 'vue'

interface OrderFilterOptions {
	status?: OrderStatus
	timeGapMinutes?: number
}

interface UseOrderEventsServiceOptions {
	printOnCreate?: boolean
}

const state = reactive<{ allOrders: OrderDTO[] }>({
	allOrders: [],
})

// --------------------------------------
// HELPERS
// --------------------------------------
const wsUrl = import.meta.env.VITE_WS_URL || 'ws://example:8080/api/v1'

function convertOrderToReactive(order: OrderDTO): OrderDTO {
	return reactive({
		...order,
		subOrders: order.subOrders?.map(subOrder => reactive(subOrder)) ?? [],
	}) as OrderDTO
}

function mergeOrder(existingOrder: OrderDTO, updatedOrder: OrderDTO) {
	// Merge top-level fields
	Object.assign(existingOrder, updatedOrder)

	// Replace subOrders with fresh reactive copies
	existingOrder.subOrders = updatedOrder.subOrders.map((sub: SuborderDTO) => reactive(sub))
}

function sortOrders(orderList: OrderDTO[]): OrderDTO[] {
	return orderList.sort((a, b) => new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime())
}

/**
 * The main hook to manage live orders via WebSocket.
 *
 * @param initialFilter - Initial filter for orders (e.g., { status: 'PENDING' })
 * @param options - Additional config, including whether to print on create.
 */
export function useOrderEvents(
	initialFilter: OrderFilterOptions = {},
	options: UseOrderEventsServiceOptions = { printOnCreate: false },
) {
	const { printOnCreate = true } = options

	// ----------------------------------
	// Setup: Toast + QR Printer
	// ----------------------------------
	const { toast } = useToast()
	const { printQR } = useQRPrinter()

	// ----------------------------------
	// Reactive Filter
	// ----------------------------------
	const localFilter = ref<OrderFilterOptions>({ ...initialFilter })

	// ----------------------------------
	// WebSocket Setup
	// ----------------------------------
	const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone
	const url = `${wsUrl}/orders/ws?timezone=${timezone}&timeGapMinutes=${localFilter.value.timeGapMinutes ?? 60}&includeYesterdayOrders=true`

	const {
		status: socketStatus,
		send,
		open,
		close,
	} = useWebSocket(url, {
		immediate: true,
		autoReconnect: {
			retries: 10,
			delay: 2000,
		},
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

	// ----------------------------------
	// Computed
	// ----------------------------------
	const filteredOrders = computed(() => {
		if (!localFilter.value.status) return state.allOrders
		return state.allOrders.filter(o => o.status === localFilter.value.status)
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
		state.allOrders.forEach(order => {
			counts[order.status] += 1
		})
		return counts
	})

	// ----------------------------------
	// Public Methods
	// ----------------------------------
	function setFilter(newFilter: OrderFilterOptions) {
		localFilter.value = { ...localFilter.value, ...newFilter }
	}

	// ----------------------------------
	// Event Handlers
	// ----------------------------------
	function handleInitialData(initialOrders: OrderDTO[]) {
		state.allOrders = sortOrders(initialOrders.map(convertOrderToReactive))
	}

	async function handleOrderCreated(newOrder: OrderDTO) {
		const [reactiveOrder, isNew] = upsertOrder(newOrder)
		if (!isNew || !printOnCreate) return

		const { width, height } = getSavedBaristaQRSettings()

		try {
			const suborderQRs = reactiveOrder.subOrders.map(s => generateSubOrderQR(s))
			printQR(suborderQRs, { labelHeightMm: height, labelWidthMm: width, desktopOnly: true })
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
		state.allOrders = state.allOrders.filter(o => o.id !== orderId)
	}

	// ----------------------------------
	// Upsert Helper
	// ----------------------------------
	function upsertOrder(orderData: OrderDTO): [OrderDTO, boolean] {
		const idx = state.allOrders.findIndex(o => o.id === orderData.id)
		if (idx !== -1) {
			// Merge updates into the existing reactive order
			mergeOrder(state.allOrders[idx], orderData)
			state.allOrders = sortOrders(state.allOrders)
			return [state.allOrders[idx], false]
		} else {
			// Insert as a new reactive order
			const reactiveOrder = convertOrderToReactive(orderData)
			state.allOrders.unshift(reactiveOrder)
			state.allOrders = sortOrders(state.allOrders)
			return [reactiveOrder, true]
		}
	}

	// ----------------------------------
	// Return API
	// ----------------------------------
	return {
		// State
		status: socketStatus,
		filteredOrders,
		orderCountsByStatus,

		// Methods
		setFilter,
		send,
		open,
		close,
	}
}
