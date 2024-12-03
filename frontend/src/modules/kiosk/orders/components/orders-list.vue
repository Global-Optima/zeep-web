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
						<span>{{ order.type === 'Delivery' ? 'Доставка' : 'Кафе' }}</span
						>, <span>{{ order.suborders.length }} шт.</span>
					</div>
				</div>
				<div>
					<template v-if="order.status === 'Active'">
						<p class="text-blue-600">{{ order.eta }}</p>
					</template>
					<template v-else-if="order.status === 'In Delivery'">
						<Truck class="w-5 h-5 text-yellow-500" />
					</template>
					<template v-else-if="order.status === 'Completed'">
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
import { Check, Truck } from 'lucide-vue-next'

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

defineProps<{
  orders: Order[];
  selectedOrder: Order | null;
}>();

const emits = defineEmits<{
  (e: 'selectOrder', order: Order): void;
}>();

const selectOrder = (order: Order) => {
  emits('selectOrder', order);
};
</script>
