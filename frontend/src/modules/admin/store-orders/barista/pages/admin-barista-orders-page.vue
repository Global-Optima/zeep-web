<script setup lang="ts">
import { useToast } from '@/core/components/ui/toast'
import { getRouteName } from '@/core/config/routes.config'
import { useBarcodeScanner } from '@/core/hooks/use-barcode-listener.hook'
import AdminBaristaOrderStatusSelector from '@/modules/admin/store-orders/barista/components/admin-barista-order-status-selector.vue'
import AdminBaristaOrdersList from '@/modules/admin/store-orders/barista/components/admin-barista-orders-list.vue'
import AdminBaristaSubOrderDetails from '@/modules/admin/store-orders/barista/components/admin-barista-sub-order-details.vue'
import AdminBaristaSubordersList from '@/modules/admin/store-orders/barista/components/admin-barista-suborders-list.vue'
import { useOrderEventsService } from '@/modules/admin/store-orders/barista/hooks/use-orders-event.hook'
import {
  OrderStatus,
  SubOrderStatus,
  type OrderDTO,
  type SuborderDTO,
} from '@/modules/admin/store-orders/models/orders.models'
import { ordersService } from '@/modules/admin/store-orders/services/orders.service'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

/**
 * Status type for your top-level filter buttons
 */
interface Status {
  label: string
  count: number
  status?: OrderStatus
}

const router = useRouter()
const { toast } = useToast()

/**
 * Selected items in the UI
 */
const selectedOrder = ref<OrderDTO | null>(null)
const selectedSuborder = ref<SuborderDTO | null>(null)

/**
 * Keep track of the currently selected status.
 * This correlates with an `OrderStatus` filter in your list.
 */
const selectedStatus = ref<Status>({
  label: 'Все',
  count: 0,
  status: undefined
})

/**
 * Default statuses (shown until fetched data arrives)
 */
const defaultStatuses = ref<Status[]>([
  { label: 'Все', count: 0 },
  { label: 'Активные', count: 0 },
  { label: 'Завершенные', count: 0 },
  { label: 'В доставке', count: 0 },
])

/* ============================
   Fetching Status Counts
============================ */
async function fetchStatuses(): Promise<Status[]> {
  const data = await ordersService.getStatusesCount()
  return [
    { label: 'Все',         count: data.ALL,        status: undefined },
    { label: 'Активные',    count: data.PREPARING,  status: OrderStatus.PREPARING },
    { label: 'Завершенные', count: data.COMPLETED,  status: OrderStatus.COMPLETED },
    { label: 'В доставке',  count: data.IN_DELIVERY, status: OrderStatus.IN_DELIVERY },
  ]
}

/**
 * useQuery for status counts.
 * You might call `refetch()` or rely on invalidation when orders change.
 */
const { data: fetchedStatuses } = useQuery({
  queryKey: ['order-statuses'],
  queryFn: fetchStatuses,
})

/* ============================
   Real-Time Order Handling
============================ */
const { filteredOrders, setFilter } = useOrderEventsService({
  status: selectedStatus.value.status,
})

/**
 * If you’d like to automatically update the filter whenever `selectedStatus` changes,
 * you can watch it. Alternatively, you can call onSelectStatus() as you do now.
 */
watch(
  () => selectedStatus.value.status,
  (newStatus) => {
    setFilter({ status: newStatus })
  }
)

/* ============================
   Computed: Displayed Statuses
============================ */
const displayedStatuses = computed<Status[]>(() => {
  return fetchedStatuses.value || defaultStatuses.value
})

/* ============================
   Selection / Navigation
============================ */
/**
 * Select an Order from the list
 */
function selectOrder(order: OrderDTO) {
  if (selectedOrder.value?.id === order.id) return
  selectedOrder.value = order
  selectedSuborder.value = null
}

/**
 * Select a Suborder
 */
function selectSuborder(suborder: SuborderDTO) {
  if (selectedSuborder.value?.id === suborder.id) return
  selectedSuborder.value = suborder
}

/**
 * Mark a suborder as completed.
 * If *all* suborders are done, mark the entire order completed and reset selection.
 */
async function toggleSuborderStatus(suborder: SuborderDTO) {
  if (suborder.status === SubOrderStatus.COMPLETED) return
  if (!selectedOrder.value) return

  try {
    await ordersService.completeSubOrder(selectedOrder.value.id, suborder.id)
    suborder.status = SubOrderStatus.COMPLETED

    // Check if all suborders are now completed
    const allDone = selectedOrder.value.subOrders.every(
      (so) => so.status === SubOrderStatus.COMPLETED
    )

    if (allDone) {
      selectedOrder.value.status = OrderStatus.COMPLETED
      // If we’re currently filtering by an "active" status, the order
      // will vanish from the list, so clear the selected order/suborder
      selectedOrder.value = null
      selectedSuborder.value = null
    }
  } catch (error) {
    console.error('Failed to complete suborder:', error)
    toast({ description: 'Не удалось завершить подзаказ', variant: 'destructive' })
  }
}

/**
 * Handle the back button
 */
function onBackClick() {
  router.push({ name: getRouteName('ADMIN_DASHBOARD') })
}

/**
 * Handle status selection (e.g. "Все", "Активные", etc.)
 */
async function onSelectStatus(status: Status) {
  selectedStatus.value = status
  // Clear the current selections
  selectedOrder.value = null
  selectedSuborder.value = null

  setFilter({ status: status.status })

  await scrollToTop()
}

/**
 * Scroll to the top smoothly
 */
async function scrollToTop() {
  if (typeof window !== 'undefined' && window.scrollTo) {
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

/**
 * If an order’s data changes (via websockets or barcode scan),
 * we want to ensure the selected order/suborder still exist
 * in the current filtered list. Otherwise, clear them.
 */
watch(
  () => filteredOrders.value,
  () => {
    if (!selectedOrder.value) return

    const exists = filteredOrders.value.some(o => o.id === selectedOrder.value?.id)
    if (!exists) {
      // The selected order is no longer in the list, clear it
      selectedOrder.value = null
      selectedSuborder.value = null
    } else if (selectedOrder.value && selectedSuborder.value) {
      // Ensure the suborder is still in the order
      const subExists = selectedOrder.value.subOrders.some(
        so => so.id === selectedSuborder.value?.id
      )
      if (!subExists) {
        selectedSuborder.value = null
      }
    }
  }
)

/* ============================
   Barcode Scanner Integration
============================ */
useBarcodeScanner({
  prefix: 'suborder-',
  onScan: async (suborderIdStr: string) => {
    try {
      const suborderId = Number(suborderIdStr)
      // 1) Complete suborder on the server; returns updated suborder
      const updatedSuborder = await ordersService.completeSubOrderById(suborderId)

      // 2) Update local data for the matching Order in the filtered list
      const localOrder = filteredOrders.value.find(
        (o) => o.id === updatedSuborder.orderId
      )
      if (localOrder) {
        // Locate the suborder
        const localSub = localOrder.subOrders.find(
          (so) => so.id === updatedSuborder.id
        )
        // Update suborder status
        if (localSub) {
          localSub.status = updatedSuborder.status
        }
        // If all suborders completed, mark order as completed
        const allDone = localOrder.subOrders.every(
          (so) => so.status === SubOrderStatus.COMPLETED
        )
        if (allDone) {
          localOrder.status = OrderStatus.COMPLETED
        }

        // Select them in the UI if desired
        selectedOrder.value = localOrder
        selectedSuborder.value = localSub ?? null
      }

      // 3) Show success toast
      toast({
        description: `Подзаказ ${updatedSuborder.productSize.productName} ${updatedSuborder.productSize.sizeName} был выполнен`,
      })
    } catch (error) {
      console.error('Barcode Scan Error:', error)
      toast({ description: 'Не удалось завершить подзаказ', variant: 'destructive' })
    }
  },
  onError: (err: Error) => {
    console.error('Barcode Scan Error:', err)
    toast({ description: 'Произошла ошибка при сканировании', variant: 'destructive' })
  },
})
</script>

<template>
	<div class="relative bg-gray-100 pt-safe w-full h-screen overflow-hidden">
		<!-- Header: Order Status Selector -->
		<AdminBaristaOrderStatusSelector
			:statuses="displayedStatuses"
			:selectedStatus="selectedStatus"
			@selectStatus="onSelectStatus"
			@back="onBackClick"
		/>

		<!-- Main Layout -->
		<div class="relative grid grid-cols-4 bg-gray-100 pb-4 w-full h-[calc(100vh-74px)]">
			<!-- Left: Orders List -->
			<AdminBaristaOrdersList
				:orders="filteredOrders"
				:selectedOrder="selectedOrder"
				@selectOrder="selectOrder"
			/>

			<!-- Middle: Suborders List -->
			<AdminBaristaSubordersList
				:suborders="selectedOrder?.subOrders || null"
				:selectedSuborder="selectedSuborder"
				@selectSuborder="selectSuborder"
			/>

			<!-- Right: SubOrder Details -->
			<AdminBaristaSubOrderDetails
				:suborder="selectedSuborder"
				@toggleSuborderStatus="toggleSuborderStatus"
			/>
		</div>
	</div>
</template>
