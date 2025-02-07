<!-- OrdersList.vue -->
<template>
	<section class="col-span-1 border-r h-full overflow-y-auto no-scrollbar">
		<p class="top-0 z-10 sticky bg-gray-100 py-3 font-medium text-center">Заказы</p>

		<!-- Show the list if we have any orders -->
		<div
			v-if="hasOrders"
			class="flex flex-col gap-2 px-2"
		>
			<div
				v-for="order in orders"
				:key="order.id"
				@click="selectOrder(order)"
				:class="orderClasses(order)"
			>
				<div>
					<div class="flex flex-col text-lg">
						<p class="text-gray-500">#{{ order.id }}</p>
						<p class="font-medium">{{ order.customerName }}</p>
					</div>
					<div class="mt-1 text-gray-700 text-sm">
						<!-- Translate orderType -->
						<span>{{ order.deliveryAddressId !== null ? 'Доставка' : 'Кафе' }}</span
						>,
						<span>{{ order.subOrdersQuantity }} шт.</span>
					</div>
				</div>
				<div>
					<!-- Adjust status icons/text as needed -->
					<template v-if="order.status === OrderStatus.PREPARING">
						<!-- If you want to show ETA as text -->
						<p class="text-blue-600">{{ formatEta(order.subOrdersQuantity) }}</p>
					</template>
					<template v-else-if="order.status === OrderStatus.IN_DELIVERY">
						<Truck class="w-5 h-5 text-yellow-500" />
					</template>
					<template v-else-if="order.status === OrderStatus.COMPLETED">
						<Check class="w-5 h-5 text-green-500" />
					</template>
				</div>
			</div>
		</div>

		<!-- If no orders, show a placeholder message -->
		<p
			v-else
			class="mt-2 text-center text-gray-400 text-sm"
		>
			Список заказов пуст
		</p>
	</section>
</template>

<script setup lang="ts">
import { OrderStatus, type OrderDTO } from '@/modules/admin/store-orders/models/orders.models'
import { Check, Truck } from 'lucide-vue-next'
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

function orderClasses(order: OrderDTO) {
  return [
    'flex justify-between items-start gap-2 bg-white p-4 rounded-xl cursor-pointer border',
    selectedOrder.value?.id === order.id ? 'border-primary' : 'border-transparent'
  ]
}

const hasOrders = computed(() => orders.value.length > 0)

function formatEta(subOrdersCount: number): string {
  const baseSubOrderEta = 2
  return `${baseSubOrderEta * subOrdersCount} мин`
}
</script>
