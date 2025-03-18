<template>
	<div class="flex flex-col items-center h-screen overflow-hidden">
		<!-- Main Layout -->
		<div class="flex items-start w-full h-full">
			<!-- In Progress Orders Section -->
			<AdminOrdersDisplayList
				title="В работе"
				:orders="inProgressOrders"
				:currentPageIndex="inProgressPageIndex"
				:totalPages="totalInProgressPages"
				class="bg-slate-800"
				type="PREPARING"
				@pageChange="setInProgressPage"
			/>

			<!-- Ready Orders Section -->
			<AdminOrdersDisplayList
				title="Готовы"
				:orders="readyOrders"
				:currentPageIndex="readyPageIndex"
				:totalPages="totalReadyPages"
				class="bg-slate-900"
				type="COMPLETED"
				@pageChange="setReadyPage"
			/>
		</div>
	</div>
</template>

<script setup lang="ts">
import { useOrderEventsService } from '@/modules/admin/store-orders/barista/hooks/use-orders-event.hook'
import AdminOrdersDisplayList from '@/modules/admin/store-orders/display/components/admin-orders-display-list.vue'
import { OrderStatus } from '@/modules/admin/store-orders/models/orders.models'
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'

// Constants
const ORDERS_PER_PAGE = 6
const AUTO_PAGE_INTERVAL = 6000 // 6 seconds

// WebSocket Hook Integration
const { filteredOrders } = useOrderEventsService({ timeGapMinutes: 15 })

// Create a reactive currentTime value that updates every second,
// so that our computed filtering will refresh as time passes.
const currentTime = ref(Date.now())
const timeInterval = setInterval(() => {
	currentTime.value = Date.now()
}, 1000)

// Clean up the interval when the component is unmounted.
onBeforeUnmount(() => {
	clearInterval(timeInterval)
})

// Orders that are "In Progress" are filtered by status.
const inProgressOrders = computed(() =>
	filteredOrders.value.filter(order => order.status === OrderStatus.PREPARING)
)

// For "Ready" orders, we only show orders that are completed
// and whose completedAt is within the last 15 minutes.
const readyOrders = computed(() =>
	filteredOrders.value.filter(order => {
		if (order.status !== OrderStatus.COMPLETED) return false
		if (!order.completedAt) return true // if no completedAt date, include by default
		const completedAtTime = new Date(order.completedAt).getTime()
		return (currentTime.value - completedAtTime) < 15 * 60 * 1000
	})
)

// Pagination State
const inProgressPageIndex = ref(0)
const readyPageIndex = ref(0)

// Calculate total pages for each order section
const totalInProgressPages = computed(() =>
	Math.ceil(inProgressOrders.value.length / ORDERS_PER_PAGE)
)
const totalReadyPages = computed(() =>
	Math.ceil(readyOrders.value.length / ORDERS_PER_PAGE)
)

// Methods to update the current page indices
const setInProgressPage = (page: number) => {
	inProgressPageIndex.value = page
}
const setReadyPage = (page: number) => {
	readyPageIndex.value = page
}

// Automatically rotate pages every 6 seconds
let autoRotateInterval: ReturnType<typeof setInterval>

const rotatePages = () => {
	autoRotateInterval = setInterval(() => {
		inProgressPageIndex.value =
			(inProgressPageIndex.value + 1) % totalInProgressPages.value
		readyPageIndex.value =
			(readyPageIndex.value + 1) % totalReadyPages.value
	}, AUTO_PAGE_INTERVAL)
}

onMounted(() => {
	rotatePages()
})

onBeforeUnmount(() => {
	clearInterval(autoRotateInterval)
})
</script>
