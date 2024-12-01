<!-- SubordersList.vue -->
<template>
	<section class="col-span-1 border-r h-full overflow-y-auto no-scrollbar">
		<p class="top-0 z-10 sticky bg-gray-100 py-3 font-medium text-center">Подзаказы</p>

		<div
			v-if="suborders && suborders.length > 0"
			class="flex flex-col gap-2 px-2"
		>
			<div
				v-for="suborder in suborders"
				:key="suborder.id"
				@click="selectSuborder(suborder)"
				:class="[
          'flex justify-between items-start gap-2 bg-white p-4 rounded-xl cursor-pointer border',
          selectedSuborder?.id === suborder.id ? 'border-primary' : 'border-transparent'
        ]"
			>
				<div>
					<p class="font-medium text-lg">{{ suborder.productName }}</p>
					<p class="line-clamp-2 text-gray-700 text-sm">{{ suborder.toppings.join(', ') }}</p>
				</div>
				<div>
					<p
						v-if="suborder.status === 'In Progress'"
						class="text-blue-600"
					>
						{{ suborder.prepTime }}
					</p>
					<Check
						v-else
						class="w-5 h-5 text-green-500"
					/>
				</div>
			</div>
		</div>

		<p
			v-else
			class="mt-2 text-center text-gray-400 text-sm"
		>
			{{ suborders ? 'Список подзаказов пуст' : 'Выберите заказ' }}
		</p>
	</section>
</template>

<script setup lang="ts">
import { Check } from 'lucide-vue-next'

interface Suborder {
  id: number;
  productName: string;
  toppings: string[];
  status: 'In Progress' | 'Done';
  comments?: string;
  prepTime: string;
}

defineProps<{
  suborders: Suborder[] | null;
  selectedSuborder: Suborder | null;
}>();

const emits = defineEmits<{
  (e: 'selectSuborder', suborder: Suborder): void;
}>();

const selectSuborder = (suborder: Suborder) => {
  emits('selectSuborder', suborder);
};
</script>
