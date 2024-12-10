<!-- OrdersList.vue -->
<template>
	<section class="col-span-1 border-r h-full overflow-y-auto no-scrollbar">
		<p class="top-0 z-10 sticky bg-gray-100 py-3 font-medium text-center">Заказы</p>

		<div
			v-if="orders.length > 0"
			class="flex flex-col gap-2 px-2"
		>
			<div
				v-for="order in orders"
				:key="order.id"
				@click="selectOrder(order)"
				:class="[
          'flex justify-between items-start gap-2 bg-white p-4 rounded-xl cursor-pointer border',
          selectedOrder?.id === order.id ? 'border-primary' : 'border-transparent'
        ]"
			>
				<div>
					<div class="flex flex-col text-lg">
						<p class="text-gray-500">#{{ order.id }}</p>
						<p class="font-medium">{{ order.customerName }}</p>
					</div>
					<div class="mt-1 text-gray-700 text-sm">
						<!-- Translate orderType -->
						<span> {{ order.deliveryAddressId !== null ? 'Доставка' : 'Кафе' }} </span>,
						<span>{{ order.subOrdersQuantity }} шт.</span>
					</div>
				</div>
				<div>
					<!-- Adjust status icons/text as needed -->
					<template v-if="order.status === OrderStatus.PREPARING">
						<!-- If you want to show ETA as text -->
						<p class="text-blue-600">{{ formatEta(order.subOrdersQuantity) }}</p>
					</template>
					<template v-else-if="order.status === 'IN_DELIVERY'">
						<Truck class="w-5 h-5 text-yellow-500" />
					</template>
					<template v-else-if="order.status === 'COMPLETED'">
						<Check class="w-5 h-5 text-green-500" />
					</template>
				</div>
			</div>
		</div>

		<p
			v-else
			class="mt-2 text-center text-gray-400 text-sm"
		>
			Список заказов пуст
		</p>
	</section>
</template>

<script setup lang="ts">
import { OrderStatus, type OrderDTO } from '@/modules/orders/models/orders.models'
import { Check, Truck } from 'lucide-vue-next'

defineProps<{
  orders: OrderDTO[];
  selectedOrder: OrderDTO | null;
}>()

const emits = defineEmits<{
  (e: 'selectOrder', order: OrderDTO): void;
}>()

const selectOrder = (order: OrderDTO) => {
  emits('selectOrder', order);
}

function formatEta(subOrdersCount: number): string {
  const baseSubOrderEta = 2
  return `${baseSubOrderEta*subOrdersCount} мин`
}
</script>
