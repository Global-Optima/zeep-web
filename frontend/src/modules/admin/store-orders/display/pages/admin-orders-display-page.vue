<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'

import AdminOrdersDisplayList from '@/modules/admin/store-orders/display/components/admin-orders-display-list.vue'
import { OrderStatus, type OrdersTimeZoneFilter } from '@/modules/admin/store-orders/models/orders.models'
import { ordersService } from '@/modules/admin/store-orders/services/orders.service'

// ------------------------------------------------------------------
// Constants & Config
// ------------------------------------------------------------------
const ORDERS_PER_PAGE = 6
const AUTO_PAGE_INTERVAL = 6000 // 6 seconds


 const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone

 const filter: OrdersTimeZoneFilter = {
  timeGapMinutes: 15,
  timezone: timezone,
  includeYesterdayOrders: false,
  statuses: [OrderStatus.PREPARING, OrderStatus.COMPLETED]
 }

const {
  data: orders,
} = useQuery({
  queryKey: ['barista-orders', filter],
  queryFn: () => ordersService.getBaristaOrders(filter),
  placeholderData: [],
  initialData: [],
  refetchInterval: 5000,
  refetchOnWindowFocus: true,
  staleTime: 0,
  retry: 3,
})

// ------------------------------------------------------------------
// Order Filtering
// (Removed the 15-minute hide logic. We just split into "PREPARING" vs "COMPLETED".)
// ------------------------------------------------------------------
const inProgressOrders = computed(() => {
  return orders.value.filter(order => order.status === OrderStatus.PREPARING)
})

const readyOrders = computed(() => {
  return orders.value.filter(order => order.status === OrderStatus.COMPLETED)
})

// ------------------------------------------------------------------
// Pagination State
// ------------------------------------------------------------------
const inProgressPageIndex = ref(0)
const readyPageIndex = ref(0)

const totalInProgressPages = computed(() => {
  return Math.ceil(inProgressOrders.value.length / ORDERS_PER_PAGE)
})

const totalReadyPages = computed(() => {
  return Math.ceil(readyOrders.value.length / ORDERS_PER_PAGE)
})

function setInProgressPage(page: number) {
  inProgressPageIndex.value = page
}

function setReadyPage(page: number) {
  readyPageIndex.value = page
}

// ------------------------------------------------------------------
// Auto-Rotate Pages
// ------------------------------------------------------------------
let autoRotateInterval: ReturnType<typeof setInterval> | null = null

function rotatePages() {
  autoRotateInterval = setInterval(() => {
    // Guard so we don’t modulo by zero
    if (totalInProgressPages.value > 0) {
      inProgressPageIndex.value =
        (inProgressPageIndex.value + 1) % totalInProgressPages.value
    }
    if (totalReadyPages.value > 0) {
      readyPageIndex.value =
        (readyPageIndex.value + 1) % totalReadyPages.value
    }
  }, AUTO_PAGE_INTERVAL)
}

onMounted(() => {
  rotatePages()
})

onBeforeUnmount(() => {
  if (autoRotateInterval) {
    clearInterval(autoRotateInterval)
  }
})
</script>

<template>
	<div class="flex flex-col items-center h-screen overflow-hidden">
		<!-- If you want a loading or error state, you can add it here -->

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
