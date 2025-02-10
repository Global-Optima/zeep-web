<template>
	<section class="col-span-1 border-r h-full overflow-y-auto no-scrollbar">
		<p class="top-0 z-10 sticky bg-gray-100 py-3 font-medium text-center">Заказы</p>

		<!-- Show the list if we have any orders -->
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
						<div class="text-lg">
							<p class="text-gray-600">#{{ order.id }}</p>
							<p class="font-medium">
								{{ order.customerName }}
							</p>
						</div>
						<!-- Status Icon -->
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

					<!-- Order Details -->
					<div class="mt-1 text-gray-700 text-sm">
						<span>{{ order.deliveryAddressId !== null ? 'Доставка' : 'Кафе' }}</span
						>,
						<span>{{ order.subOrdersQuantity }} шт.</span>
					</div>

					<!-- Progress Bar for Suborders -->
					<div class="relative bg-gray-200 mt-4 rounded-sm w-full h-4 overflow-hidden">
						<div
							class="h-full transition-all duration-300 ease-in-out"
							:class="progressBarClasses(order)"
							:style="{ width: `${orderProgress(order)}%` }"
						></div>
					</div>

					<!-- Progress Count -->
					<p class="mt-1 text-gray-500 text-xs text-right">
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

const props = defineProps<{
  orders: OrderDTO[];
  selectedOrder: OrderDTO | null;
}>()

const { orders, selectedOrder } = toRefs(props)

const emits = defineEmits<{
  (e: 'selectOrder', order: OrderDTO): void;
}>()

function selectOrder(order: OrderDTO) {
  emits('selectOrder', order)
}

/**
 * 🧠 Computed: Determines if there are any orders.
 */
const hasOrders = computed(() => orders.value.length > 0)

/**
 * 🧠 Computed: Calculates completed suborders for an order.
 */
function completedSubOrders(order: OrderDTO): number {
  if (!order.subOrders) return 0
  return order.subOrders.filter(sub => sub.status === SubOrderStatus.COMPLETED).length
}

/**
 * 🧠 Computed: Determines progress percentage for an order.
 */
function orderProgress(order: OrderDTO): number {
  if (!order.subOrdersQuantity || !order.subOrders) return 0
  return (completedSubOrders(order) / order.subOrdersQuantity) * 100
}

/**
 * 🎨 Determines the color of the progress bar based on order status.
 */
function progressBarClasses(order: OrderDTO) {
  if (order.status === OrderStatus.PENDING) return 'bg-green-600'
  if (order.status === OrderStatus.PREPARING) return 'bg-blue-500'
  if (order.status === OrderStatus.IN_DELIVERY) return 'bg-yellow-500'
  if (order.status === OrderStatus.COMPLETED) return 'bg-green-600'
  return 'bg-gray-300'
}

/**
 * 🎨 Sorts orders to move `COMPLETED` ones to the bottom dynamically.
 */
const sortedOrders = computed(() => {
  return [...orders.value].sort((a, b) => {
    const orderStatusPriority: Record<OrderStatus, number> = {
      [OrderStatus.PENDING]: 1,
      [OrderStatus.PREPARING]: 2,
      [OrderStatus.IN_DELIVERY]: 3,
      [OrderStatus.COMPLETED]: 4,
      [OrderStatus.DELIVERED]: 5,
      [OrderStatus.CANCELLED]: 6
    }

    return orderStatusPriority[a.status] - orderStatusPriority[b.status]
  })
})

/**
 * 🎨 Determines the styling for each order card.
 */
function orderClasses(order: OrderDTO) {
  return cn(
    'flex items-start gap-2 p-4 rounded-xl cursor-pointer border transition-all duration-200 bg-white',
    selectedOrder.value?.id === order.id ? 'border-primary' : 'border-transparent',
    order.status === OrderStatus.PENDING ? '' : '',
    order.status === OrderStatus.PREPARING ? 'border-blue-400' : '',
    order.status === OrderStatus.COMPLETED ? 'bg-green-50 bg-opacity-50 border-green-200' : '',
  )
}

/**
 * ⏳ Estimates the total preparation time.
 */
function formatEta(subOrdersCount: number): string {
  const baseSubOrderEta = 2
  return `${baseSubOrderEta * subOrdersCount} мин`
}
</script>
