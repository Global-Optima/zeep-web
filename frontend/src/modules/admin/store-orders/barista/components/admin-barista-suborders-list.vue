<template>
	<section class="col-span-1 border-r h-full overflow-y-auto no-scrollbar">
		<p class="top-0 z-10 sticky bg-gray-100 py-3 font-medium text-center">Подзаказы</p>

		<!-- If suborders is non-null and has items -->
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
					<p class="line-clamp-2 text-gray-700 text-sm">
						{{ suborder.additives.map(a => a.additive.name).join(', ') }}
					</p>
				</div>

				<div>
					<!-- Show "2 мин" if PREPARING; otherwise, show a check icon -->
					<p
						v-if="suborder.status === SubOrderStatus.PREPARING"
						class="text-blue-600"
					>
						2 мин.
					</p>
					<Check
						v-else
						class="w-5 h-5 text-green-500"
					/>
				</div>
			</div>
		</div>

		<!-- If suborders is null or empty -->
		<p
			v-else
			class="mt-2 text-center text-gray-400 text-sm"
		>
			{{ suborders ? 'Список подзаказов пуст' : 'Выберите заказ' }}
		</p>
	</section>
</template>

<script setup lang="ts">
import { SubOrderStatus, type SuborderDTO } from '@/modules/admin/store-orders/models/orders.models'
import { Check } from 'lucide-vue-next'
import { computed, toRefs } from 'vue'

/**
 * Define props as a single `props` object, then
 * destructure them with `toRefs` to ensure reactivity.
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
 * Helper function to emit the selected suborder.
 */
function selectSuborder(suborder: SuborderDTO) {
  emits('selectSuborder', suborder)
}

/**
 * A computed boolean to check if we have suborders.
 * This cleans up the template's v-if logic.
 */
const hasSuborders = computed(() => {
  return suborders.value && suborders.value.length > 0
})

/**
 * A small helper to compute dynamic classes for each suborder item.
 */
function suborderClasses(suborder: SuborderDTO) {
  return [
    'flex justify-between items-start gap-2 bg-white p-4 rounded-xl cursor-pointer border',
    selectedSuborder.value?.id === suborder.id ? 'border-primary' : 'border-transparent'
  ]
}

/**
 * Exported references or helper functions are optional with `<script setup>`.
 * Everything declared is automatically "scoped" to this component.
 */
</script>
