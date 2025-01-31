<!-- SuborderDetails.vue -->
<template>
	<section class="col-span-2 border-r h-full overflow-y-auto no-scrollbar">
		<p class="top-0 z-10 sticky bg-gray-100 py-3 font-medium text-center">Детали подзаказа</p>

		<div
			class="px-2 rounded-xl"
			v-if="suborder"
		>
			<div class="relative flex flex-col gap-4 bg-white p-6 rounded-xl overflow-y-auto">
				<Button
					size="icon"
					variant="ghost"
					class="top-6 right-6 absolute"
					:disabled="suborder.status === SubOrderStatus.COMPLETED"
					@click="printQrCode"
				>
					<Printer
						stroke-width="1.5"
						class="!size-8"
					/>
				</Button>

				<div>
					<p class="font-medium text-xl">
						{{ suborder.productSize.productName }} {{ suborder.productSize.sizeName }}
					</p>
					<!-- Toppings List -->
					<ul
						v-if="suborder.additives.length > 0"
						class="space-y-1 mt-2"
					>
						<li
							v-for="(topping, index) in suborder.additives"
							:key="index"
							class="flex items-center"
						>
							<Plus class="mr-2 w-4 h-4 text-gray-500" />
							<span class="text-gray-700">{{ topping.additive.name }}</span>
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
						{{ 'Стандартное приготовление' }}
					</p>
				</div>

				<div>
					<p class="font-medium text-lg">Время приготовления</p>
					<p class="mt-1 text-gray-700">
						{{ '2 мин'  }}
					</p>
				</div>

				<div class="flex items-center gap-2 mt-4">
					<button
						@click="toggleSuborderStatus(suborder)"
						:disabled="suborder.status === SubOrderStatus.COMPLETED"
						:class="[
              'flex-1 px-4 py-4 rounded-xl text-primary-foreground',
              suborder.status === 'COMPLETED'
                ? 'bg-gray-300 text-gray-500 cursor-not-allowed'
                : 'bg-primary'
            ]"
					>
						{{ suborder.status === "COMPLETED" ? 'Выполнено' : 'Выполнить' }}
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
import { Button } from '@/core/components/ui/button'
import { SubOrderStatus, type SuborderDTO } from '@/modules/orders/models/orders.models'
import { ordersService } from '@/modules/orders/services/orders.service'
import { Plus, Printer } from 'lucide-vue-next'

const {suborder} = defineProps<{
  suborder: SuborderDTO | null;
}>();

const emits = defineEmits<{
  (e: 'toggleSuborderStatus', suborder: SuborderDTO): void;
}>();

const toggleSuborderStatus = (suborder: SuborderDTO) => {
  emits('toggleSuborderStatus', suborder);
};

const printQrCode = async () => {
  if (suborder) {
    await ordersService.getSuborderBarcodeFile(suborder.id);
  }
};
</script>
