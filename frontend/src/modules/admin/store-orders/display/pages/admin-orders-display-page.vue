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
const { filteredOrders } = useOrderEventsService({timeGapMinutes: 3})

// State for "In Progress" and "Ready" orders
const inProgressOrders = computed(() =>
	filteredOrders.value.filter(order => order.status === OrderStatus.PREPARING)
)
const readyOrders = computed(() =>
	filteredOrders.value.filter(order => order.status === OrderStatus.COMPLETED)
)

// Pagination State
const inProgressPageIndex = ref(0)
const readyPageIndex = ref(0)

// Pagination Logic
const totalInProgressPages = computed(() =>
	Math.ceil(inProgressOrders.value.length / ORDERS_PER_PAGE)
)
const totalReadyPages = computed(() =>
	Math.ceil(readyOrders.value.length / ORDERS_PER_PAGE)
)

// Methods to update pages
const setInProgressPage = (page: number) => {
	inProgressPageIndex.value = page
}
const setReadyPage = (page: number) => {
	readyPageIndex.value = page
}

// Automatically change pages every 6 seconds
let interval: ReturnType<typeof setInterval>

const rotatePages = () => {
	interval = setInterval(() => {
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
	clearInterval(interval)
})
</script>
