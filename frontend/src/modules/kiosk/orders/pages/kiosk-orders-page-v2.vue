<template>
	<div class="relative bg-gray-100 pt-safe w-full h-screen overflow-hidden">
		<!-- Header: Order Status Selector -->
		<OrderStatusSelector
			:statuses="fetchedStatuses || statuses"
			:selectedStatus="selectedStatus"
			@select-status="onSelectStatus"
			@back="onBackClick"
		/>

		<!-- Main Layout -->
		<div class="relative grid grid-cols-4 bg-gray-100 pb-4 w-full h-[calc(100vh-74px)]">
			<OrdersList
				:orders="filteredOrders"
				:selectedOrder="selectedOrder"
				@selectOrder="selectOrder"
			/>
			<SubordersList
				:suborders="selectedOrder?.subOrders || null"
				:selectedSuborder="selectedSuborder"
				@selectSuborder="selectSuborder"
			/>
			<SubOrderDetails
				:suborder="selectedSuborder"
				@toggleSuborderStatus="toggleSuborderStatus"
			/>
		</div>
	</div>
</template>

<script setup lang="ts">
import { getRouteName } from '@/core/config/routes.config'
import OrderStatusSelector from '@/modules/kiosk/orders/components/order-status-selector.vue'
import OrdersList from '@/modules/kiosk/orders/components/orders-list.vue'
import SubOrderDetails from '@/modules/kiosk/orders/components/sub-order-details.vue'
import SubordersList from '@/modules/kiosk/orders/components/suborders-list.vue'
import { useOrderEventsService } from '@/modules/kiosk/orders/services/orders-event.service'
import { OrderStatus, SubOrderStatus, type OrderDTO, type SubOrderDTO } from '@/modules/orders/models/orders.models'
import { orderService } from '@/modules/orders/services/orders.service'
import { useCurrentStoreStore } from '@/modules/stores/store/current-store.store'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

interface Status {
	label: string
	count: number
}

const router = useRouter()
const {currentStoreId} = useCurrentStoreStore()

const onBackClick = () => {
	router.push({ name: getRouteName('ADMIN_DASHBOARD') })
}

const selectedOrder = ref<OrderDTO | null>(null)
const selectedSuborder = ref<SubOrderDTO | null>(null)
const selectedStatus = ref<Status>({ label: 'Все', count: 0 })

const statusMap: Record<string, OrderStatus | null> = {
	"Все": null,
	"Активные": OrderStatus.PREPARING,
	"Завершенные": OrderStatus.COMPLETED,
	"В доставке": OrderStatus.IN_DELIVERY,
}

const statuses = ref<Status[]>([
	{ label: 'Все', count: 0 },
	{ label: 'Активные', count: 0 },
	{ label: 'Завершенные', count: 0 },
	{ label: 'В доставке', count: 0 },
])

const fetchStatuses = async (): Promise<Status[]> => {
  if (!currentStoreId) return []

	const data = await orderService.getStatusesCount(currentStoreId)
	return [
		{ label: 'Все', count: data.ALL },
		{ label: 'Активные', count: data.PREPARING },
		{ label: 'Завершенные', count: data.COMPLETED },
		{ label: 'В доставке', count: data.IN_DELIVERY },
	]
}

const {data: fetchedStatuses} = useQuery({
  queryKey: ['statuses', currentStoreId],
  queryFn: fetchStatuses ,
})

// Use the composable
const { getOrdersRef } = useOrderEventsService(currentStoreId!)

// Filtered orders based on selectedStatus
const filteredOrders = computed(() => {
	const backendStatus = statusMap[selectedStatus.value.label]
	let result = getOrdersRef().value

	if (backendStatus) {
		result = result.filter(o => o.status === backendStatus)
	}

	return result
})

function selectOrder(order: OrderDTO) {
	if (selectedOrder.value?.id === order.id) return
	selectedOrder.value = order
	selectedSuborder.value = null
}

function selectSuborder(suborder: SubOrderDTO) {
	if (selectedSuborder.value?.id === suborder.id) return
	selectedSuborder.value = suborder
}

const onSelectStatus = async (status: Status) => {
  selectedStatus.value = status;
  selectedOrder.value = null;
  selectedSuborder.value = null;
  await scrollToTop();
};

const scrollToTop = async () => {
  if (window && window.scrollTo) {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }
};

async function toggleSuborderStatus(suborder: SubOrderDTO) {
	if (suborder.status === SubOrderStatus.COMPLETED) return

	if (!selectedOrder.value) return

	const orderId = selectedOrder.value.id
	const suborderId = suborder.id
  if (!currentStoreId) return

	try {
		await orderService.completeSubOrder(currentStoreId, orderId, suborderId)
		suborder.status = SubOrderStatus.COMPLETED

		const allDone = selectedOrder.value?.subOrders.every((so) => so.status === SubOrderStatus.COMPLETED)
		if (allDone) {
			selectedOrder.value.status = OrderStatus.COMPLETED
			selectedOrder.value = null
			selectedSuborder.value = null
		}
	} catch (error) {
		console.error('Failed to complete suborder:', error)
	}
}
</script>
