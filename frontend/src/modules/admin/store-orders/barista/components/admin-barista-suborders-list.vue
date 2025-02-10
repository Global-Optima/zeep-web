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
import { SubOrderStatus, type SuborderDTO } from '@/modules/admin/store-orders/models/orders.models'
import { Check, Clock } from 'lucide-vue-next'
import { computed, toRefs } from 'vue'

const props = defineProps<{
  suborders: SuborderDTO[] | null;
  selectedSuborder: SuborderDTO | null;
}>()

const emits = defineEmits<{
  (e: 'selectSuborder', suborder: SuborderDTO): void;
}>()

const { suborders, selectedSuborder } = toRefs(props)

/**
 * Selects a suborder and emits the event.
 */
function selectSuborder(suborder: SuborderDTO) {
  emits('selectSuborder', suborder)
}

/**
 * Checks if there are any suborders.
 */
const hasSuborders = computed(() => suborders.value && suborders.value.length > 0)

/**
 * Computes dynamic classes for styling each suborder based on its status.
 */
function suborderClasses(suborder: SuborderDTO) {
  return [
    'flex justify-between items-start gap-2 bg-white p-4 rounded-xl cursor-pointer border',
    selectedSuborder.value?.id === suborder.id ? '!border-primary' : 'border-transparent',
    suborder.status === SubOrderStatus.PENDING ? 'bg-gray-100' : '',
    suborder.status === SubOrderStatus.PREPARING ? 'border-blue-400' : '',
    suborder.status === SubOrderStatus.COMPLETED ? 'bg-green-50 bg-opacity-50 border-green-200' : '',
  ]
}
</script>
