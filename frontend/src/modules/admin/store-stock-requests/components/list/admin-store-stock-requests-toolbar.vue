<template>
	<div
		class="flex md:flex-row flex-col justify-between items-start md:items-center space-y-4 md:space-y-0 mb-4"
	>
		<!-- Left Side: Filter Menu -->
		<div class="flex items-center space-x-2 w-full md:w-auto">
			<MultiSelectFilter
				title="Статусы"
				:options="statusOptions"
				v-model="selectedStatuses"
			/>
		</div>

		<div class="flex items-center space-x-2 w-full md:w-auto">
			<Button variant="outline">Экспорт</Button>
		</div>
	</div>
</template>

<script setup lang="ts">
import MultiSelectFilter from '@/core/components/multi-select-filter/MultiSelectFilter.vue'
import { Button } from '@/core/components/ui/button'
import { StockRequestStatus, type GetStockRequestsFilter } from '@/modules/admin/store-stock-requests/models/stock-requests.model'
import { ref, watch } from 'vue'

const props = defineProps<{ filter?: GetStockRequestsFilter }>()
const emit = defineEmits(['update:filter'])


const selectedStatuses = ref<StockRequestStatus[]>(props.filter?.statuses ?? [])


watch(selectedStatuses, (newStatuses) => {
  emit('update:filter', {
    ...props.filter,
    statuses: newStatuses.length ? newStatuses : undefined,
  })
})

const statusOptions = [
  { label: 'Созданные', value: StockRequestStatus.CREATED },
  { label: 'Обработанные', value: StockRequestStatus.PROCESSED },
  { label: 'В доставке', value: StockRequestStatus.IN_DELIVERY },
  { label: 'Завершённые', value: StockRequestStatus.COMPLETED },
  { label: 'Отклонённые магазином', value: StockRequestStatus.REJECTED_BY_STORE },
  { label: 'Отклонённые складом', value: StockRequestStatus.REJECTED_BY_WAREHOUSE },
  { label: 'Принятые с изменениями', value: StockRequestStatus.ACCEPTED_WITH_CHANGE },
]
</script>
