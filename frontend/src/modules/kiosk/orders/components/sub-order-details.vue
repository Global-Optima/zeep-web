<!-- SuborderDetails.vue -->
<template>
	<section class="col-span-2 border-r h-full overflow-y-auto no-scrollbar">
		<p class="top-0 z-10 sticky bg-gray-100 py-3 font-medium text-center">Детали подзаказа</p>

		<div
			class="px-2 rounded-xl"
			v-if="suborder"
		>
			<div class="relative flex flex-col gap-4 bg-white p-6 rounded-xl overflow-y-auto">
				<button
					class="top-6 right-6 absolute"
					:disabled="suborder.status === 'Done'"
					@click="printQrCode"
				>
					<Printer
						stroke-width="1.5"
						class="w-8 h-8"
					/>
				</button>

				<div>
					<p class="font-medium text-xl">{{ suborder.productName }}</p>
					<!-- Toppings List -->
					<ul
						v-if="suborder.toppings.length > 0"
						class="space-y-1 mt-2"
					>
						<li
							v-for="(topping, index) in suborder.toppings"
							:key="index"
							class="flex items-center"
						>
							<Plus class="mr-2 w-4 h-4 text-gray-500" />
							<span class="text-gray-700">{{ topping }}</span>
						</li>
					</ul>
					<p
						v-else
						class="mt-2 text-gray-700"
					>
						Без топпингов
					</p>
				</div>

				<div>
					<p class="font-medium text-lg">Комментарий</p>
					<p class="mt-1 text-gray-700">
						{{ suborder.comments || 'Стандартное приготовление' }}
					</p>
				</div>

				<div>
					<p class="font-medium text-lg">Время приготовления</p>
					<p class="mt-1 text-gray-700">
						{{ suborder.prepTime || 'Не указано' }}
					</p>
				</div>

				<div class="flex items-center gap-2 mt-4">
					<button
						@click="toggleSuborderStatus(suborder)"
						:disabled="suborder.status === 'Done'"
						:class="[
              'flex-1 px-4 py-4 rounded-xl text-primary-foreground',
              suborder.status === 'Done'
                ? 'bg-gray-300 text-gray-500 cursor-not-allowed'
                : 'bg-primary'
            ]"
					>
						{{ suborder.status === 'Done' ? 'Выполнено' : 'Выполнить' }}
					</button>
				</div>
			</div>
		</div>

		<p
			v-else
			class="mt-2 text-center text-gray-400 text-sm"
		>
			Выберите подзаказ
		</p>
	</section>
</template>

<script setup lang="ts">
import { Plus, Printer } from 'lucide-vue-next'

interface Suborder {
  id: number;
  productName: string;
  toppings: string[];
  status: 'In Progress' | 'Done';
  comments?: string;
  prepTime: string;
}

const props = defineProps<{
  suborder: Suborder | null;
}>();

const emits = defineEmits<{
  (e: 'toggleSuborderStatus', suborder: Suborder): void;
  (e: 'printQrCode', suborder: Suborder): void;
}>();

const toggleSuborderStatus = (suborder: Suborder) => {
  emits('toggleSuborderStatus', suborder);
};

const printQrCode = () => {
  if (props.suborder) {
    emits('printQrCode', props.suborder);
  }
};
</script>
