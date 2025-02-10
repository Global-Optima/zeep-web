import { useToast } from '@/core/components/ui/toast'
import { useBarcodePrinter } from '@/core/hooks/use-barcode-print.hook'
import type {
	OrderDTO,
	OrderStatus,
	SuborderDTO,
} from '@/modules/admin/store-orders/models/orders.models'
import { useWebSocket } from '@vueuse/core'
import { computed, reactive, ref, watchEffect } from 'vue'

interface OrderFilterOptions {
	status?: OrderStatus
}

// Global reactive state for all orders
const state = reactive<{ allOrders: OrderDTO[] }>({
	allOrders: [],
})

const wsUrl = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/api/v1'

/**
 * Convert a raw OrderDTO into a deeply reactive object.
 */
function convertOrderToReactive(order: OrderDTO): OrderDTO {
	return reactive({
		...order,
		subOrders: order.subOrders?.map(subOrder => reactive(subOrder)) ?? [],
	}) as OrderDTO
}

/**
 * Merge updated fields into an existing reactive order without losing references.
 * This preserves reactivity on nested objects (like subOrders).
 */
function mergeOrder(existingOrder: OrderDTO, updatedOrder: OrderDTO) {
	// Merge top-level properties
	Object.assign(existingOrder, updatedOrder)

	// Replace subOrders with new reactive copies to ensure correct reactivity
	existingOrder.subOrders = updatedOrder.subOrders.map((sub: SuborderDTO) => reactive(sub))
}

/**
 * Sort orders by creation time (ascending). Adjust if you prefer descending.
 */
function sortOrders(orderList: OrderDTO[]): OrderDTO[] {
	return orderList.sort((a, b) => new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime())
}

export function useOrderEventsService(initialFilter: OrderFilterOptions = {}) {
	const { toast } = useToast()
	const { printBarcode } = useBarcodePrinter()

	// WebSocket URL includes timezone, if desired
	const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone
	const url = `${wsUrl}/orders/ws?timezone=${timezone}`

	// Local filter is a ref so we can update it reactively
	const localFilter = ref({ ...initialFilter })

	// Configure our WebSocket
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
			console.error('WebSocket Error:', err)
			toast({
				title: 'Ошибка соединения',
				description: 'Не удалось подключиться к серверу.',
				variant: 'destructive',
			})
		},
		onDisconnected: () => {
			toast({
				title: 'Вы отключились',
				description: 'Соединение с заказами было прервано.',
			})
		},
	})

	// Handle incoming WebSocket messages
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
			console.error('Не удалось разобрать сообщение WebSocket:', err)
			toast({
				title: 'Ошибка обработки данных',
				description: 'Не удалось обработать данные от сервера.',
				variant: 'destructive',
			})
		}
	})

	// ----------------------------------
	// Public Computed + Methods
	// ----------------------------------

	/**
	 * Filter orders by the chosen status. If no status is specified,
	 * return all orders.
	 */
	const filteredOrders = computed(() => {
		if (!localFilter.value.status) return state.allOrders
		return state.allOrders.filter(o => o.status === localFilter.value.status)
	})

	/**
	 * Tally counts of each status.
	 */
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
			counts[order.status] += 1
		})
		return counts
	})

	/**
	 * Update our local filter.
	 */
	function setFilter(newFilter: OrderFilterOptions) {
		localFilter.value = { ...localFilter.value, ...newFilter }
	}

	// ----------------------------------
	// Message/Event Handlers
	// ----------------------------------

	/** Replace entire list with sorted reactive orders. */
	function handleInitialData(initialOrders: OrderDTO[]) {
		state.allOrders = sortOrders(initialOrders.map(convertOrderToReactive))
	}

	/**
	 * On "order_created" event: upsert the order.
	 * If it's truly new, print barcodes.
	 */
	async function handleOrderCreated(newOrder: OrderDTO) {
		const [reactiveOrder, isNew] = upsertOrder(newOrder)
		if (!isNew) return // If the order was already there, skip printing

		// Print barcodes for suborders if needed
		try {
			await Promise.all(
				reactiveOrder.subOrders.map(async sub => {
					const name = `${sub.productSize.productName} ${sub.productSize.sizeName}`
					await printBarcode(name, `suborder-${sub.id}`, { showModal: false })
					console.log(`Suborder with id ${sub.id} printed`)
				}),
			)
		} catch (err) {
			console.error('Ошибка при печати штрих-кода:', err)
			toast({
				title: 'Ошибка печати',
				description: 'Не удалось напечатать один или несколько штрих-кодов.',
				variant: 'destructive',
			})
		}
	}

	/**
	 * On "order_updated" event: upsert the order (merge updates).
	 */
	function handleOrderUpdated(updated: OrderDTO) {
		upsertOrder(updated)
	}

	/**
	 * On "order_deleted" event: remove the order by ID.
	 */
	function handleOrderDeleted(orderId: number) {
		state.allOrders = state.allOrders.filter(o => o.id !== orderId)
	}

	// ----------------------------------
	// Internal Helpers
	// ----------------------------------

	/**
	 * Upsert an order:
	 *  - If `order.id` exists, merge changes (preserve references).
	 *  - Otherwise, insert a new reactive order at the front.
	 * Returns [theReactiveOrder, isNew].
	 */
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
	// Return the Hook API
	// ----------------------------------

	return {
		// Reactive data
		status: socketStatus,
		filteredOrders,
		orderCountsByStatus,

		// Actions
		setFilter,
		send,
		open,
		close,
	}
}
