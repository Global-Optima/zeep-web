<script setup lang="ts">
import { useOrderEvents } from '@/modules/admin/store-orders/barista/hooks/use-orders-event.hook'
import AdminOrdersDisplayList from '@/modules/admin/store-orders/display/components/admin-orders-display-list.vue'
import { OrderStatus } from '@/modules/admin/store-orders/models/orders.models'
import { ordersService } from '@/modules/admin/store-orders/services/orders.service'
import { useQuery } from '@tanstack/vue-query'
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'

// Constants
const ORDERS_PER_PAGE = 6
const AUTO_PAGE_INTERVAL = 6000 // 6 seconds
const TIME_GAP_MINUTES = 15

// Initial HTTP fetch for immediate data
const { data: initialOrders, isPending: isInitialLoading } = useQuery({
  queryKey: ['initial-orders', TIME_GAP_MINUTES],
  queryFn: () => ordersService.getBaristaOrders({timeGapMinutes: TIME_GAP_MINUTES }),
  retry: 3,
  refetchInterval: 20_000
})

// WebSocket integration
const { filteredOrders, status: socketStatus } = useOrderEvents({
  timeGapMinutes: TIME_GAP_MINUTES
})

// Merge HTTP and WebSocket data
const mergedOrders = computed(() => {
  if (socketStatus.value === "OPEN") {
    return filteredOrders.value
  }
  return initialOrders.value ?? []
})

// Time tracking
const currentTime = ref(Date.now())
const timeInterval = setInterval(() => {
  currentTime.value = Date.now()
}, 1000)

onBeforeUnmount(() => clearInterval(timeInterval))

// Order filtering
const inProgressOrders = computed(() =>
  mergedOrders.value.filter(order =>
    order.status === OrderStatus.PREPARING
  )
)

const readyOrders = computed(() =>
  mergedOrders.value.filter(order => {
    if (order.status !== OrderStatus.COMPLETED) return false
    if (!order.completedAt) return true
    const completedTime = new Date(order.completedAt).getTime()
    return (currentTime.value - completedTime) < TIME_GAP_MINUTES * 60 * 1000
  })
)

// Pagination state
const inProgressPageIndex = ref(0)
const readyPageIndex = ref(0)

// Total pages calculation
const totalInProgressPages = computed(() =>
  Math.ceil(inProgressOrders.value.length / ORDERS_PER_PAGE)
)

const totalReadyPages = computed(() =>
  Math.ceil(readyOrders.value.length / ORDERS_PER_PAGE)
)

// Page change handlers
const setInProgressPage = (page: number) => {
  inProgressPageIndex.value = page
}

const setReadyPage = (page: number) => {
  readyPageIndex.value = page
}

// Auto-rotate pages
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

<template>
	<div class="flex flex-col items-center h-screen overflow-hidden">
		<!-- Loading state -->
		<div
			v-if="isInitialLoading"
			class="flex justify-center items-center h-full"
		>
			<div class="border-primary border-t-4 rounded-full w-12 h-12 animate-spin"></div>
		</div>

		<!-- Main content -->
		<div class="flex items-start w-full h-full">
			<!-- In Progress Orders -->
			<AdminOrdersDisplayList
				title="В работе"
				:orders="inProgressOrders"
				:currentPageIndex="inProgressPageIndex"
				:totalPages="totalInProgressPages"
				class="bg-slate-800"
				type="PREPARING"
				@pageChange="setInProgressPage"
			/>

			<!-- Ready Orders -->
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
