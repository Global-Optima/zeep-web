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
  type SuborderDTO
} from '@/modules/admin/store-orders/models/orders.models'
import { ordersService } from '@/modules/admin/store-orders/services/orders.service'

import { computed, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

/* ------------------------------------------------------
   Basic Setup
------------------------------------------------------ */
const router = useRouter()
const { toast } = useToast()

/** Currently selected order + suborder in the UI */
const selectedOrder = ref<OrderDTO | null>(null)
const selectedSuborder = ref<SuborderDTO | null>(null)

/** Track the currently selected status (for filtering). */
const selectedStatus = ref<{ label: string; count: number; status?: OrderStatus }>({
  label: 'Все',
  count: 0,
  status: undefined
})

/**
 * Real-Time Orders:
 * Pull from our updated `useOrderEventsService`, which returns deeply reactive data.
 */
const { filteredOrders, orderCountsByStatus, setFilter } = useOrderEventsService({
  status: selectedStatus.value.status,
})

/**
 * Convert orderCountsByStatus into an array structure for the
 * AdminBaristaOrderStatusSelector or any other consumer.
 */
const displayedStatuses = computed(() => {
  const counts = orderCountsByStatus.value
  return [
    {
      label: 'Все',
      count: Object.values(counts).reduce((sum, val) => sum + val, 0),
      status: undefined,
    },
    {
      label: 'В ожидании',
      count: counts[OrderStatus.PENDING] || 0,
      status: OrderStatus.PENDING,
    },
    {
      label: 'Активные',
      count: counts[OrderStatus.PREPARING] || 0,
      status: OrderStatus.PREPARING,
    },
    {
      label: 'Завершенные',
      count: counts[OrderStatus.COMPLETED] || 0,
      status: OrderStatus.COMPLETED,
    },
    {
      label: 'В доставке',
      count: counts[OrderStatus.IN_DELIVERY] || 0,
      status: OrderStatus.IN_DELIVERY,
    },
  ]
})

/* ------------------------------------------------------
   Watchers & Filters
------------------------------------------------------ */

/** Update filter whenever `selectedStatus` changes */
watch(
  () => selectedStatus.value.status,
  newStatus => {
    setFilter({ status: newStatus })
  }
)

/**
 * Ensure selected items remain valid in the current `filteredOrders`.
 * If the order or suborder disappears (e.g., changed status), clear the selection.
 */
watch(
  () => filteredOrders.value,
  () => {
    if (!selectedOrder.value) return

    // Still in list?
    const stillExists = filteredOrders.value.some(o => o.id === selectedOrder.value?.id)
    if (!stillExists) {
      selectedOrder.value = null
      selectedSuborder.value = null
      return
    }

    // If order still exists, ensure suborder does too
    if (selectedOrder.value && selectedSuborder.value) {
      const subExists = selectedOrder.value.subOrders.some(
        so => so.id === selectedSuborder.value?.id
      )
      if (!subExists) {
        selectedSuborder.value = null
      }
    }
  }
)

/* ------------------------------------------------------
   Selection & Actions
------------------------------------------------------ */

/** User clicks an order in the list */
function selectOrder(order: OrderDTO) {
  if (selectedOrder.value?.id === order.id) return
  selectedOrder.value = order
  selectedSuborder.value = null
}

/** User clicks a suborder in the suborders list */
function selectSuborder(suborder: SuborderDTO) {
  if (selectedSuborder.value?.id === suborder.id) return
  selectedSuborder.value = suborder
}

/**
 * Toggle suborder status by calling the server.
 *  - Provide immediate feedback by assigning the returned status.
 *  - If all suborders are completed => order is effectively `COMPLETED` => unselect.
 */
async function toggleSuborderStatus(suborder: SuborderDTO) {
  if (!selectedOrder.value) return

  try {
    // Request the next status from the server
    const updatedSuborder = await ordersService.toggleNextStatus(suborder.id)

    console.log("UPDATEEEEDDD", updatedSuborder)

    // Immediate local feedback (optional, since WS will update eventually)
    Object.assign(suborder, updatedSuborder)

    // If all suborders are completed, close the selection
    if (areAllSubordersCompleted(selectedOrder.value)) {
      // We can also set the order's status to COMPLETED for local feedback:
      selectedOrder.value.status = OrderStatus.COMPLETED

      // Deselect the order and suborder to close the details panel
      selectedOrder.value = null
      selectedSuborder.value = null
    }
  } catch (error) {
    console.error('Failed to complete suborder:', error)
    toast({ description: 'Не удалось изменить статус подзаказа', variant: 'destructive' })
  }
}

/* ------------------------------------------------------
   Barcode Scanner
------------------------------------------------------ */

/**
 * Barcode scanning:
 *  1. We expect suborder barcodes with prefix `suborder-`.
 *  2. On scan, we toggle the suborder via the same `ordersService`.
 *  3. Optionally do local feedback or rely on websockets.
 */
useBarcodeScanner({
  prefix: 'suborder-',
  onScan: async (suborderIdStr: string) => {
    try {
      const suborderId = Number(suborderIdStr)
      const updatedSuborder = await ordersService.toggleNextStatus(suborderId)

      // Find the local order in filteredOrders
      const localOrder = filteredOrders.value.find(o => o.id === updatedSuborder.orderId)
      if (localOrder) {
        // Optionally update the local suborder for immediate feedback
        const localSub = localOrder.subOrders.find(so => so.id === updatedSuborder.id)
        if (localSub) Object.assign(localSub, updatedSuborder)

        // If all suborders are completed, unselect
        if (areAllSubordersCompleted(localOrder)) {
          localOrder.status = OrderStatus.COMPLETED
          selectedOrder.value = null
          selectedSuborder.value = null
        } else {
          // Otherwise, select the updated suborder
          selectedOrder.value = localOrder
          selectedSuborder.value = localSub ?? null
        }
      }

      toast({
        description: `Статус подзаказа ${updatedSuborder.productSize.productName} ${updatedSuborder.productSize.sizeName} был изменен`,
      })
    } catch (error) {
      console.error('Barcode Scan Error:', error)
      toast({ description: 'Не удалось изменить статус подзаказа', variant: 'destructive' })
    }
  },
  onError: (err: Error) => {
    console.error('Barcode Scan Error:', err)
    toast({ description: 'Произошла ошибка при сканировании', variant: 'destructive' })
  },
})

/* ------------------------------------------------------
   Utility: Check if all suborders in an order are COMPLETED
------------------------------------------------------ */
function areAllSubordersCompleted(order: OrderDTO) {
  return order.subOrders.every(sub => sub.status === SubOrderStatus.COMPLETED)
}

/* ------------------------------------------------------
   Navigation + Reload
------------------------------------------------------ */
function onBackClick() {
  router.push({ name: getRouteName('ADMIN_DASHBOARD') })
}

function onReloadClick() {
  window.location.reload()
}

/* ------------------------------------------------------
   Scroll to top
------------------------------------------------------ */
function scrollToTop() {
  if (typeof window !== 'undefined' && window.scrollTo) {
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

function onSelectStatus(status: { label: string; count: number; status?: OrderStatus }) {
  // Update the currently selected status
  selectedStatus.value = status

  // Clear any currently selected order/suborder
  selectedOrder.value = null
  selectedSuborder.value = null

  setFilter({ status: status.status })
  scrollToTop()
}
</script>

<template>
	<div class="relative bg-gray-100 pt-safe w-full h-screen overflow-hidden">
		<!-- Header: Order Status Selector -->
		<AdminBaristaOrderStatusSelector
			:statuses="displayedStatuses"
			:selectedStatus="selectedStatus"
			@select-status="onSelectStatus"
			@back="onBackClick"
			@reload="onReloadClick"
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
