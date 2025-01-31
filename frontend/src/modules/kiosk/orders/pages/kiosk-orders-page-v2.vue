<script setup lang="ts">
import { useToast } from '@/core/components/ui/toast'
import { getRouteName } from '@/core/config/routes.config'
import { useBarcodeScanner } from '@/core/hooks/use-barcode-listener.hook'
import OrderStatusSelector from '@/modules/kiosk/orders/components/order-status-selector.vue'
import OrdersList from '@/modules/kiosk/orders/components/orders-list.vue'
import SubOrderDetails from '@/modules/kiosk/orders/components/sub-order-details.vue'
import SubordersList from '@/modules/kiosk/orders/components/suborders-list.vue'
import { useOrderEventsService } from '@/modules/kiosk/orders/services/orders-event.service'
import {
  OrderStatus,
  SubOrderStatus,
  type OrderDTO,
  type SuborderDTO,
} from '@/modules/orders/models/orders.models'
import { ordersService } from '@/modules/orders/services/orders.service'
import { useQuery } from '@tanstack/vue-query'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

interface Status {
  label: string
  count: number
  status?: OrderStatus
}

const router = useRouter()
const { toast } = useToast()

const selectedOrder = ref<OrderDTO | null>(null)
const selectedSuborder = ref<SuborderDTO | null>(null)
const selectedStatus = ref<Status>({ label: 'Все', count: 0, status: undefined })

// Default statuses (shown if not yet fetched)
const statuses = ref<Status[]>([
  { label: 'Все', count: 0 },
  { label: 'Активные', count: 0 },
  { label: 'Завершенные', count: 0 },
  { label: 'В доставке', count: 0 },
])

// Fetch status counts
async function fetchStatuses(): Promise<Status[]> {
  const data = await ordersService.getStatusesCount()
  return [
    { label: 'Все', count: data.ALL, status: undefined },
    { label: 'Активные', count: data.PREPARING, status: OrderStatus.PREPARING },
    { label: 'Завершенные', count: data.COMPLETED, status: OrderStatus.COMPLETED },
    { label: 'В доставке', count: data.IN_DELIVERY, status: OrderStatus.IN_DELIVERY },
  ]
}

const { data: fetchedStatuses } = useQuery({
  queryKey: ['order-statuses'],
  queryFn: fetchStatuses,
})


const { filteredOrders, setFilter } = useOrderEventsService({
  status: selectedStatus.value.status,
})

/**
 * Select an Order in the UI
 */
function selectOrder(order: OrderDTO) {
  if (selectedOrder.value?.id === order.id) return
  selectedOrder.value = order
  selectedSuborder.value = null
}

/**
 * Select a Suborder in the UI
 */
function selectSuborder(suborder: SuborderDTO) {
  if (selectedSuborder.value?.id === suborder.id) return
  selectedSuborder.value = suborder
}

/**
 * Toggle the Suborder status via an action from SubOrderDetails component
 */
async function toggleSuborderStatus(suborder: SuborderDTO) {
  if (suborder.status === SubOrderStatus.COMPLETED) return
  if (!selectedOrder.value) return

  try {
    await ordersService.completeSubOrder(selectedOrder.value.id, suborder.id)
    suborder.status = SubOrderStatus.COMPLETED

    const allDone = selectedOrder.value?.subOrders.every(
      (so) => so.status === SubOrderStatus.COMPLETED
    )

    if (allDone) {
      selectedOrder.value.status = OrderStatus.COMPLETED
      selectedOrder.value = null
      selectedSuborder.value = null
    }
  } catch (error) {
    console.error('Failed to complete suborder:', error)
  }
}

/**
 * Back button logic
 */
function onBackClick() {
  router.push({ name: getRouteName('ADMIN_DASHBOARD') })
}

/**
 * Re-select status (reset Order/Suborder) and scroll
 */
async function onSelectStatus(status: Status) {
  selectedStatus.value = status
  selectedOrder.value = null
  selectedSuborder.value = null

  setFilter({status: status.status})

  await scrollToTop()
}

/**
 * Smooth scroll to top
 */
async function scrollToTop() {
  if (window && window.scrollTo) {
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

/**
 * Barcode Scanner
 *
 * - We look for barcodes with prefix: "suborder-"
 * - On successful scan, we complete that suborder by ID,
 *   then we update local UI state so that the user sees it completed
 *   and the correct order/suborder get highlighted.
 */
useBarcodeScanner({
  prefix: 'suborder-',
  onScan: async (suborderIdStr: string) => {
    try {
      const suborderId = Number(suborderIdStr)
      // 1) Complete suborder on the server; returns updated suborder
      const updatedSuborder = await ordersService.completeSubOrderById(suborderId)

      // 2) Update local data for the matching Order -> Suborder
      const localOrder = filteredOrders.value.find(
        (o) => o.id === updatedSuborder.orderId
      )

      if (localOrder) {
        // Locate the suborder in local state
        const localSub = localOrder.subOrders.find(
          (so) => so.id === updatedSuborder.id
        )

        // Update suborder status in local store
        if (localSub) {
          localSub.status = updatedSuborder.status
        }

        // If all suborders are completed, mark the entire order completed
        const allDone = localOrder.subOrders.every(
          (so) => so.status === SubOrderStatus.COMPLETED
        )
        if (allDone) {
          localOrder.status = OrderStatus.COMPLETED
        }

        // Now select them in the UI
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
    toast({ description: 'Произошла ошибка при сканировании', variant: 'destructive'})
  },
})
</script>

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
			<!-- Left: Orders List -->
			<OrdersList
				:orders="filteredOrders"
				:selectedOrder="selectedOrder"
				@selectOrder="selectOrder"
			/>

			<!-- Middle: Suborders List -->
			<SubordersList
				:suborders="selectedOrder?.subOrders || null"
				:selectedSuborder="selectedSuborder"
				@selectSuborder="selectSuborder"
			/>

			<!-- Right: SubOrder Details -->
			<SubOrderDetails
				:suborder="selectedSuborder"
				@toggleSuborderStatus="toggleSuborderStatus"
			/>
		</div>
	</div>
</template>
