<template>
	<section class="col-span-2 border-r h-full overflow-y-auto no-scrollbar">
		<p class="top-0 z-10 sticky bg-gray-100 py-3 font-medium text-center">Детали подзаказа</p>

		<!-- If suborder is provided, show its details -->
		<div
			v-if="suborder"
			class="relative flex flex-col gap-4 bg-white mx-2 p-6 rounded-3xl overflow-y-auto"
		>
			<!-- Print Button -->
			<Button
				size="icon"
				variant="ghost"
				type="button"
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
					{{ suborder.productSize.productName }}
					{{ suborder.productSize.sizeName }} ({{suborder.productSize.size}}
					{{ suborder.productSize.unit.name.toLowerCase() }})
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
					Без модификаторов
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

			<!-- Complete (or Next Status) Button -->
			<div class="flex items-center gap-2 mt-4">
				<button
					@click="toggleSuborderStatus(suborder)"
					:disabled="disabledCompleteButton"
					:class="cn(
              'flex-1 px-4 py-4 rounded-xl text-primary-foreground font-medium',
              suborder.status === SubOrderStatus.COMPLETED
                ? 'bg-gray-300 text-gray-500 cursor-not-allowed'
                : suborder.status === SubOrderStatus.PENDING
                ? 'bg-blue-500'
                : 'bg-primary'
            )"
				>
					{{ completeButtonText }}
				</button>
			</div>
		</div>

		<!-- If no suborder is selected, show a placeholder -->
		<p
			v-else
			class="mt-2 text-gray-400 text-sm text-center"
		>
			Выберите подзаказ
		</p>
	</section>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { getSavedBaristaQRSettings } from '@/core/hooks/use-qr-print.hook'
import { cn } from '@/core/utils/tailwind.utils'
import { useOrderQRPrinter } from '@/modules/admin/store-orders/hooks/use-order-qr-print.hook'
import {
  SubOrderStatus,
  type OrderDTO,
  type SuborderDTO
} from '@/modules/admin/store-orders/models/orders.models'
import { Plus, Printer } from 'lucide-vue-next'
import { computed } from 'vue'

/**
 * Define props:
 * - suborder: a single SuborderDTO or null (if none is selected).
 *   This suborder is assumed to be "deeply reactive" from the parent.
 */
const {suborder, order} = defineProps<{
  order: OrderDTO | null;
  suborder: SuborderDTO | null;
}>()

/**
 * Define events:
 * - toggleSuborderStatus: inform the parent to toggle status.
 */
const emits = defineEmits<{
  (e: 'toggleSuborderStatus', suborder: SuborderDTO): void;
}>()

/**
 * Emit an event to the parent to toggle the suborder's status.
 */
function toggleSuborderStatus(s: SuborderDTO) {
  emits('toggleSuborderStatus', s)
}

/**
 * Dynamically compute the button label.
 */
const completeButtonText = computed(() => {
  if (!suborder) return 'Обновить статус'
  switch (suborder.status) {
    case SubOrderStatus.PENDING:
      return 'Начать приготовление'
    case SubOrderStatus.PREPARING:
      return 'Завершить'
    case SubOrderStatus.COMPLETED:
      return 'Выполнено'
    default:
      return 'Обновить статус'
  }
})

/**
 * Disable the complete button if the suborder is already completed.
 */
const disabledCompleteButton = computed(() =>
  suborder?.status === SubOrderStatus.COMPLETED
)

/**
 * A hook to print a barcode for this suborder.
 */
const { printSubOrderQR } = useOrderQRPrinter()

async function printQrCode() {
  if (suborder && order) {
    const {width, height} = getSavedBaristaQRSettings()

    await printSubOrderQR(order, suborder, {
      labelWidthMm: width,
      labelHeightMm: height,
      desktopOnly: false
    })
  }
}
</script>
