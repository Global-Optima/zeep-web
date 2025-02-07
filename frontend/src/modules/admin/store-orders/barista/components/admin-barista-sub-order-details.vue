<!-- SuborderDetails.vue -->
<template>
	<section class="col-span-2 border-r h-full overflow-y-auto no-scrollbar">
		<p class="top-0 z-10 sticky bg-gray-100 py-3 font-medium text-center">Детали подзаказа</p>

		<!-- If suborder is provided, show its details -->
		<div
			class="px-2 rounded-xl"
			v-if="suborder"
		>
			<div class="relative flex flex-col gap-4 bg-white p-6 rounded-xl overflow-y-auto">
				<!-- Print Button -->
				<Button
					size="icon"
					variant="ghost"
					class="top-6 right-6 absolute"
					:disabled="disabledCompleteButton"
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
					<p class="mt-1 text-gray-700">Стандартное приготовление</p>
				</div>

				<div>
					<p class="font-medium text-lg">Время приготовления</p>
					<p class="mt-1 text-gray-700">2 мин</p>
				</div>

				<!-- Complete Button -->
				<div class="flex items-center gap-2 mt-4">
					<button
						@click="toggleSuborderStatus(suborder)"
						:disabled="disabledCompleteButton"
						:class="[
              'flex-1 px-4 py-4 rounded-xl text-primary-foreground',
              suborder.status === SubOrderStatus.COMPLETED
                ? 'bg-gray-300 text-gray-500 cursor-not-allowed'
                : 'bg-primary'
            ]"
					>
						{{ completeButtonText }}
					</button>
				</div>
			</div>
		</div>

		<!-- If no suborder is selected, show a placeholder -->
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
import { usePrinter } from '@/core/hooks/use-print.hook'
import { SubOrderStatus, type SuborderDTO } from '@/modules/admin/store-orders/models/orders.models'
import { ordersService } from '@/modules/admin/store-orders/services/orders.service'
import { Plus, Printer } from 'lucide-vue-next'
import { computed, toRefs } from 'vue'

/**
 * Define props
 */
const props = defineProps<{
  suborder: SuborderDTO | null;
}>()

const { suborder } = toRefs(props)

/**
 * Define emits
 */
const emits = defineEmits<{
  (e: 'toggleSuborderStatus', suborder: SuborderDTO): void;
}>()

function toggleSuborderStatus(s: SuborderDTO) {
  emits('toggleSuborderStatus', s)
}

/**
 * Computed properties for the "Complete" button state
 */
const completeButtonText = computed(() =>
  suborder.value?.status === SubOrderStatus.COMPLETED
    ? 'Выполнено'
    : 'Выполнить'
)

const disabledCompleteButton = computed(() =>
  suborder.value?.status === SubOrderStatus.COMPLETED
)

/**
 * Use your printer hook
 */
const { print } = usePrinter()

async function printQrCode() {
  if (suborder.value) {
    const blob = await ordersService.getSuborderBarcodeFile(suborder.value.id)
    await print(blob)
  }
}
</script>
