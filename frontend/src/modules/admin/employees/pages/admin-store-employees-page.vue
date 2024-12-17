<template>
	<AdminStoreEmployeesToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<Card>
		<CardContent class="mt-4">
			<p
				v-if="!employees || employees.length === 0"
				class="text-muted-foreground"
			>
				Сотрудники не найдены
			</p>
			<AdminStoreEmployeesList :employees="employees" />
		</CardContent>
	</Card>
</template>

<script setup lang="ts">

import { Card, CardContent } from '@/core/components/ui/card'
import AdminStoreEmployeesList from '@/modules/admin/employees/components/list/admin-store-employees-list.vue'
import AdminStoreEmployeesToolbar from '@/modules/admin/employees/components/list/admin-store-employees-toolbar.vue'
import type { EmployeesFilter } from '@/modules/employees/models/employees.models'
import { employeesService } from '@/modules/employees/services/employees.service'
import { useCurrentStoreStore } from '@/modules/stores/store/current-store.store'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const filter = ref<EmployeesFilter>({
  search: '',
})

const {currentStoreId} = useCurrentStoreStore()

const queryKey = computed(() => (['employees', filter.value]))

const { data: employees } = useQuery({
  queryKey: queryKey,
  queryFn: () => {
    if (!currentStoreId) throw new Error('No store ID available')
    return employeesService.getStoreEmployees(currentStoreId, filter.value)
  },
  initialData: [],
  enabled: computed(() => !!currentStoreId),
})

function updateFilter(updatedFilter: EmployeesFilter) {
  filter.value = {...filter.value, ...updatedFilter}
}
</script>
