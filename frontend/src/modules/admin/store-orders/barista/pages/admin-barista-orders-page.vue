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
import { computed, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter();
const { toast } = useToast();

// Selected items in the UI
const selectedOrder = ref<OrderDTO | null>(null);
const selectedSuborder = ref<SuborderDTO | null>(null);

// Keep track of the currently selected status
const selectedStatus = ref<{ label: string; count: number; status?: OrderStatus }>({
  label: 'Все',
  count: 0,
  status: undefined,
});

// Real-Time Order Handling
const { filteredOrders, orderCountsByStatus, setFilter } = useOrderEventsService({
  status: selectedStatus.value.status,
});

// Map orderCountsByStatus to the displayed statuses format
const displayedStatuses = computed(() => {
  const counts = orderCountsByStatus.value;
  return [
    { label: 'Все', count: Object.values(counts).reduce((sum, count) => sum + count, 0), status: undefined },
    { label: 'Активные', count: counts[OrderStatus.PREPARING] || 0, status: OrderStatus.PREPARING },
    { label: 'Завершенные', count: counts[OrderStatus.COMPLETED] || 0, status: OrderStatus.COMPLETED },
    { label: 'В доставке', count: counts[OrderStatus.IN_DELIVERY] || 0, status: OrderStatus.IN_DELIVERY },
  ];
});

// Watch for changes in selectedStatus and update the filter
watch(
  () => selectedStatus.value.status,
  (newStatus) => {
    setFilter({ status: newStatus });
  }
);

// Utility function to check if all suborders are completed
function areAllSubordersCompleted(order: OrderDTO): boolean {
  return order.subOrders.every((so) => so.status === SubOrderStatus.COMPLETED);
}

// Select an Order from the list
function selectOrder(order: OrderDTO) {
  if (selectedOrder.value?.id === order.id) return;
  selectedOrder.value = order;
  selectedSuborder.value = null;
}

// Select a Suborder
function selectSuborder(suborder: SuborderDTO) {
  if (selectedSuborder.value?.id === suborder.id) return;
  selectedSuborder.value = suborder;
}

// Mark a suborder as completed
async function toggleSuborderStatus(suborder: SuborderDTO) {
  if (suborder.status === SubOrderStatus.COMPLETED) return;
  if (!selectedOrder.value) return;

  try {
    await ordersService.completeSubOrder(selectedOrder.value.id, suborder.id);
    suborder.status = SubOrderStatus.COMPLETED;

    // Check if all suborders are now completed
    if (areAllSubordersCompleted(selectedOrder.value)) {
      selectedOrder.value.status = OrderStatus.COMPLETED;
      selectedOrder.value = null;
      selectedSuborder.value = null;
    }
  } catch (error) {
    console.error('Failed to complete suborder:', error);
    toast({ description: 'Не удалось завершить подзаказ', variant: 'destructive' });
  }
}

// Handle the back button
function onBackClick() {
  router.push({ name: getRouteName('ADMIN_DASHBOARD') });
}

// Handle reload
function onReloadClick() {
  window.location.reload();
}

// Handle status selection
function onSelectStatus(status: { label: string; count: number; status?: OrderStatus }) {
  selectedStatus.value = status;
  selectedOrder.value = null;
  selectedSuborder.value = null;
  setFilter({ status: status.status });
  scrollToTop();
}

// Scroll to the top smoothly
function scrollToTop() {
  if (typeof window !== 'undefined' && window.scrollTo) {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }
}

// Ensure selected order/suborder still exist in the filtered list
watch(
  () => filteredOrders.value,
  () => {
    if (!selectedOrder.value) return;

    const exists = filteredOrders.value.some((o) => o.id === selectedOrder.value?.id);
    if (!exists) {
      selectedOrder.value = null;
      selectedSuborder.value = null;
    } else if (selectedOrder.value && selectedSuborder.value) {
      const subExists = selectedOrder.value.subOrders.some(
        (so) => so.id === selectedSuborder.value?.id
      );
      if (!subExists) {
        selectedSuborder.value = null;
      }
    }
  }
);

// Barcode Scanner Integration
useBarcodeScanner({
  prefix: 'suborder-',
  onScan: async (suborderIdStr: string) => {
    try {
      const suborderId = Number(suborderIdStr);
      const updatedSuborder = await ordersService.completeSubOrderById(suborderId);

      const localOrder = filteredOrders.value.find((o) => o.id === updatedSuborder.orderId);
      if (localOrder) {
        const localSub = localOrder.subOrders.find((so) => so.id === updatedSuborder.id);
        if (localSub) {
          localSub.status = updatedSuborder.status;
        }

        // Check if all suborders are now completed
        if (areAllSubordersCompleted(localOrder)) {
          localOrder.status = OrderStatus.COMPLETED;
        }

        selectedOrder.value = localOrder;
        selectedSuborder.value = localSub ?? null;
      }

      toast({
        description: `Подзаказ ${updatedSuborder.productSize.productName} ${updatedSuborder.productSize.sizeName} был выполнен`,
      });
    } catch (error) {
      console.error('Barcode Scan Error:', error);
      toast({ description: 'Не удалось завершить подзаказ', variant: 'destructive' });
    }
  },
  onError: (err: Error) => {
    console.error('Barcode Scan Error:', err);
    toast({ description: 'Произошла ошибка при сканировании', variant: 'destructive' });
  },
});
</script>

<template>
	<div class="relative bg-gray-100 pt-safe w-full h-screen overflow-hidden">
		<!-- Header: Order Status Selector -->
		<AdminBaristaOrderStatusSelector
			:statuses="displayedStatuses"
			:selectedStatus="selectedStatus"
			@selectStatus="onSelectStatus"
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
