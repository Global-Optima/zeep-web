<template>
	<section class="col-span-1 border-r h-full overflow-y-auto no-scrollbar">
		<p class="top-0 z-10 sticky bg-gray-100 py-3 font-medium text-center">Заказы</p>

		<!-- If we have orders, show the list -->
		<div
			v-if="hasOrders"
			class="flex flex-col gap-2 px-2"
		>
			<div
				v-for="order in sortedOrders"
				:key="order.id"
				@click="selectOrder(order)"
				:class="orderClasses(order)"
			>
				<div class="flex flex-col w-full">
					<!-- Order ID and Name -->
					<div class="flex justify-between">
						<div class="text-base">
							<p class="text-gray-600">#{{ order.id }}</p>
							<p class="mt-1 font-medium">
								{{ order.customerName }}
							</p>
						</div>

						<!-- Status Icon / ETA -->
						<div>
							<template v-if="order.status === OrderStatus.PENDING">
								<Clock class="w-5 h-5 text-gray-500" />
							</template>
							<template v-else-if="order.status === OrderStatus.PREPARING">
								<p class="text-blue-600">{{ formatEta(order.subOrdersQuantity) }}</p>
							</template>
							<template v-else-if="order.status === OrderStatus.IN_DELIVERY">
								<Truck class="w-5 h-5 text-yellow-500" />
							</template>
							<template v-else-if="order.status === OrderStatus.COMPLETED">
								<Check class="w-5 h-5 text-green-800" />
							</template>
						</div>
					</div>

					<!-- Delivery/Cafe + Suborder count -->
					<div
						v-if="order.status !== OrderStatus.COMPLETED"
						class="mt-1 text-gray-700 text-sm"
					>
						<span>{{ order.deliveryAddressId !== null ? 'Доставка' : 'Кафе' }}</span
						>,
						<span>{{ order.subOrdersQuantity }} шт.</span>
					</div>

					<!-- Progress Bar for suborders -->
					<div
						v-if="order.status !== OrderStatus.COMPLETED"
						class="relative bg-slate-200 mt-4 rounded-[6px] w-full h-4 overflow-hidden"
					>
						<div
							class="h-full transition-all duration-300 ease-in-out"
							:class="progressBarClasses(order)"
							:style="{ width: `${orderProgress(order)}%` }"
						></div>
					</div>

					<!-- Progress Count -->
					<p class="mt-2 text-gray-500 text-xs">
						{{ completedSubOrders(order) }} / {{ order.subOrdersQuantity }} выполнено
					</p>
				</div>
			</div>
		</div>

		<!-- If no orders, show a placeholder message -->
		<p
			v-else
			class="mt-2 text-gray-400 text-sm text-center"
		>
			Список заказов пуст
		</p>
	</section>
</template>

<script setup lang="ts">
import { cn } from '@/core/utils/tailwind.utils'
import { OrderStatus, SubOrderStatus, type OrderDTO } from '@/modules/admin/store-orders/models/orders.models'
import { Check, Clock, Truck } from 'lucide-vue-next'
import { computed, toRefs } from 'vue'

/**
 * Props
 * - `orders`: deeply reactive array of OrderDTO from the parent
 * - `selectedOrder`: currently selected order (if any)
 */
const props = defineProps<{
  orders: OrderDTO[];
  selectedOrder: OrderDTO | null;
}>()

const emits = defineEmits<{
  (e: 'selectOrder', order: OrderDTO): void;
}>()

const { orders, selectedOrder } = toRefs(props)

/**
 * Emits an event to the parent when an order is clicked.
 */
function selectOrder(order: OrderDTO) {
  emits('selectOrder', order)
}

/**
 * Computed: Check if we have any orders to display.
 */
const hasOrders = computed(() => orders.value.length > 0)

/**
 * Count how many suborders are completed for a given order.
 */
function completedSubOrders(order: OrderDTO): number {
  if (!order.subOrders) return 0
  return order.subOrders.filter(sub => sub.status === SubOrderStatus.COMPLETED).length
}

/**
 * Determine the progress of suborders (0 to 100).
 */
function orderProgress(order: OrderDTO): number {
  if (!order.subOrdersQuantity || !order.subOrders) return 0
  return (completedSubOrders(order) / order.subOrdersQuantity) * 100
}

/**
 * Determine the color of the progress bar based on the order's status.
 */
function progressBarClasses(order: OrderDTO) {
  if (order.status === OrderStatus.PENDING) return 'bg-green-600'
  if (order.status === OrderStatus.PREPARING) return 'bg-blue-500'
  if (order.status === OrderStatus.IN_DELIVERY) return 'bg-yellow-500'
  if (order.status === OrderStatus.COMPLETED) return 'bg-green-600'
  return 'bg-slate-100'
}

/**
 * Apply styling / highlight for the selected order and
 * optionally for different statuses.
 */
function orderClasses(order: OrderDTO) {
  return cn(
    'flex items-start gap-2 p-4 rounded-xl cursor-pointer border transition-all duration-200 bg-white',
    selectedOrder.value?.id === order.id ? '!border-primary' : 'border-transparent',

    // Example: highlight different statuses (optional)
    order.status === OrderStatus.PENDING ? '' : '',
    order.status === OrderStatus.PREPARING ? 'bg-blue-50 bg-opacity-50 border-blue-200' : '',
    order.status === OrderStatus.COMPLETED ? 'bg-green-50 bg-opacity-50 border-green-200' : '',
  )
}

/**
 * Sort orders so that COMPLETED (and other 'finished' statuses) appear last.
 * Since `orders` is already reactive (from parent), we simply copy and sort.
 */
const sortedOrders = computed(() => {
  const orderStatusPriority: Record<OrderStatus, number> = {
    [OrderStatus.PENDING]: 1,
    [OrderStatus.PREPARING]: 1,
    [OrderStatus.IN_DELIVERY]: 1,
    [OrderStatus.COMPLETED]: 2,
    [OrderStatus.DELIVERED]: 2,
    [OrderStatus.CANCELLED]: 2,
  }

  return [...orders.value].sort((a, b) => {
    return orderStatusPriority[a.status] - orderStatusPriority[b.status]
  })
})

/**
 * Estimate the total preparation time (example logic).
 * Adjust to your actual logic for each suborder if necessary.
 */
function formatEta(subOrdersCount: number): string {
  const baseSubOrderEta = 2 // 2 minutes per sub-order (example)
  return `${baseSubOrderEta * subOrdersCount} мин`
}
</script>
