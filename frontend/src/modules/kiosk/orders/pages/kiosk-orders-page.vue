<!-- OrderManagement.vue -->
<template>
	<div class="relative bg-gray-100 pt-safe w-full h-screen overflow-hidden">
		<!-- Header: Order Status Selector -->
		<OrderStatusSelector
			:statuses="statuses"
			:selectedStatus="selectedStatus"
			@selectStatus="onSelectStatus"
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
				:suborders="selectedOrder?.suborders || null"
				:selectedSuborder="selectedSuborder"
				@selectSuborder="selectSuborder"
			/>

			<SubOrderDetails
				:suborder="selectedSuborder"
				@toggleSuborderStatus="toggleSuborderStatus"
				@printQrCode="printQrCode"
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
import { computed, nextTick, ref, watch } from 'vue'
import { useRouter } from 'vue-router'


interface Suborder {
  id: number;
  productName: string;
  toppings: string[];
  status: 'In Progress' | 'Done';
  comments?: string;
  prepTime: string;
}

interface Order {
  id: number;
  customerName: string;
  customerEmail: string;
  details: string;
  eta: string;
  suborders: Suborder[];
  status: 'Active' | 'Completed' | 'In Delivery';
  type: 'Delivery' | 'In-Store';
}

interface Status {
  label: string;
  count: number;
}

const router = useRouter();

const onBackClick = () => {
  router.push({ name: getRouteName('ADMIN_DASHBOARD') });
};

const scrollToTop = async () => {
  await nextTick();
  if (window && window.scrollTo) {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }
};

/**
 * Mock data with different statuses and order types
 */
const orders = ref<Order[]>([
  // ... (Your mock data remains the same)
]);

/**
 * Reactive states
 */
const statuses = ref<Status[]>([
  { label: 'Все', count: orders.value.length },
  { label: 'Активные', count: orders.value.filter((o) => o.status === 'Active').length },
  { label: 'Завершенные', count: orders.value.filter((o) => o.status === 'Completed').length },
  { label: 'В доставке', count: orders.value.filter((o) => o.status === 'In Delivery').length },
]);

const selectedStatus = ref<Status>(statuses.value[0]);
const selectedOrder = ref<Order | null>(null);
const selectedSuborder = ref<Suborder | null>(null);

/**
 * Computed values
 */
const filteredOrders = computed(() => {
  switch (selectedStatus.value.label) {
    case 'Активные':
      return orders.value.filter((order) => order.status === 'Active');
    case 'Завершенные':
      return orders.value.filter((order) => order.status === 'Completed');
    case 'В доставке':
      return orders.value.filter((order) => order.status === 'In Delivery');
    default:
      return orders.value;
  }
});

/**
 * Methods
 */
const onSelectStatus = async (status: Status) => {
  selectedStatus.value = status;
  selectedOrder.value = null;
  selectedSuborder.value = null;
  await scrollToTop();
};

const selectOrder = async (order: Order) => {
  if (selectedOrder.value?.id === order.id) return;
  selectedOrder.value = order;
  selectedSuborder.value = null;
  await scrollToTop();
};

const selectSuborder = async (suborder: Suborder) => {
  if (selectedSuborder.value?.id === suborder.id) return;
  selectedSuborder.value = suborder;
  await scrollToTop();
};

/**
 * Toggle suborder status between 'In Progress' and 'Done'
 */
const toggleSuborderStatus = (suborder: Suborder) => {
  if (suborder.status === 'Done') return;
  suborder.status = 'Done';

  // Check if all suborders are done to mark the order as completed
  const allDone = selectedOrder.value?.suborders.every((so) => so.status === 'Done');
  if (allDone && selectedOrder.value?.status === 'Active') {
    selectedOrder.value.status = 'Completed';
    updateStatusCounts();
    // Unselect order and suborder when order is completed
    selectedOrder.value = null;
    selectedSuborder.value = null;
  }
};

const printQrCode = () => {
  console.log('Print QR Code');
};

/**
 * Watchers
 */
watch(
  () => orders.value,
  () => {
    updateStatusCounts();
  },
  { deep: true }
);

/**
 * Update counts in statuses based on orders
 */
const updateStatusCounts = () => {
  statuses.value = statuses.value.map((status) => {
    switch (status.label) {
      case 'Все':
        return { ...status, count: orders.value.length };
      case 'Активные':
        return { ...status, count: orders.value.filter((o) => o.status === 'Active').length };
      case 'Завершенные':
        return { ...status, count: orders.value.filter((o) => o.status === 'Completed').length };
      case 'В доставке':
        return { ...status, count: orders.value.filter((o) => o.status === 'In Delivery').length };
      default:
        return status;
    }
  });
};
</script>
