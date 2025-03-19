<template>
	<section class="col-span-1 border-r h-full overflow-y-auto no-scrollbar">
		<p class="top-0 z-10 sticky bg-gray-100 py-3 font-medium text-center">Подзаказы</p>

		<!-- If suborders exist -->
		<div
			v-if="hasSuborders"
			class="flex flex-col gap-2 px-2"
		>
			<div
				v-for="suborder in suborders"
				:key="suborder.id"
				@click="selectSuborder(suborder)"
				:class="suborderClasses(suborder)"
			>
				<div>
					<p class="font-medium text-lg">
						{{ suborder.productSize.productName }} {{ suborder.productSize.sizeName }}
					</p>
					<p class="text-gray-700 text-sm line-clamp-2">
						{{ suborder.additives.map(a => a.additive.name).join(', ') || 'Без топпингов' }}
					</p>
				</div>

				<div>
					<!-- Status handling -->
					<template v-if="suborder.status === SubOrderStatus.PENDING">
						<Clock class="w-5 h-5 text-gray-500" />
					</template>
					<template v-else-if="suborder.status === SubOrderStatus.PREPARING">
						<p class="text-blue-600">2 мин.</p>
					</template>
					<template v-else-if="suborder.status === SubOrderStatus.COMPLETED">
						<Check class="w-5 h-5 text-green-500" />
					</template>
				</div>
			</div>
		</div>

		<!-- If no suborders exist -->
		<p
			v-else
			class="mt-2 text-gray-400 text-sm text-center"
		>
			{{ suborders ? 'Список подзаказов пуст' : 'Выберите заказ' }}
		</p>
	</section>
</template>

<script setup lang="ts">
import { cn } from '@/core/utils/tailwind.utils'
import {
  type SuborderDTO,
  SubOrderStatus,
} from '@/modules/admin/store-orders/models/orders.models'
import { Check, Clock } from 'lucide-vue-next'
import { computed, toRefs } from 'vue'

/**
 * Props
 * - `suborders`: an array of SuborderDTO or null if no order is selected
 * - `selectedSuborder`: the currently selected suborder (if any)
 */
const props = defineProps<{
  suborders: SuborderDTO[] | null;
  selectedSuborder: SuborderDTO | null;
}>()

const emits = defineEmits<{
  (e: 'selectSuborder', suborder: SuborderDTO): void;
}>()

const { suborders, selectedSuborder } = toRefs(props)

/**
 * Check if there are suborders to display.
 */
const hasSuborders = computed(() => suborders.value && suborders.value.length > 0)

/**
 * Emit an event to the parent when a suborder is clicked.
 */
function selectSuborder(suborder: SuborderDTO) {
  emits('selectSuborder', suborder)
}

/**
 * Dynamically style each suborder card based on status and selection.
 */
function suborderClasses(suborder: SuborderDTO) {
  return cn(
    'flex items-start justify-between gap-2 p-4 rounded-xl cursor-pointer border transition-all duration-200 bg-white',
    // Highlight if it's the selected suborder
    selectedSuborder.value?.id === suborder.id ? '!border-primary' : 'border-transparent',

    // Optionally style suborders by status
    suborder.status === SubOrderStatus.PENDING ? '' : '',
    suborder.status === SubOrderStatus.PREPARING ? 'bg-blue-50 bg-opacity-50 border-blue-200' : '',
    suborder.status === SubOrderStatus.COMPLETED ? 'bg-green-50 bg-opacity-50 border-green-200' : '',
  )
}
</script>
