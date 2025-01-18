<template>
	<AdminStoreEmployeesToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<Card>
		<CardContent class="mt-4">
			<p
				v-if="!employees || employees.data.length === 0"
				class="text-muted-foreground"
			>
				Сотрудники не найдены
			</p>
			<!-- TODO: add pagination here -->
			<AdminStoreEmployeesList
				v-else
				:employees="employees.data"
			/>
		</CardContent>
	</Card>
</template>

<script setup lang="ts">

import { Card, CardContent } from '@/core/components/ui/card'
import AdminStoreEmployeesList from '@/modules/admin/employees/components/list/admin-store-employees-list.vue'
import AdminStoreEmployeesToolbar from '@/modules/admin/employees/components/list/admin-store-employees-toolbar.vue'
import type { StoreEmployeesFilter } from '@/modules/employees/models/employees.models'
import { employeesService } from '@/modules/employees/services/employees.service'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'


const filter = ref<StoreEmployeesFilter>({
  search: '',
})

const queryKey = computed(() => (['employees', filter.value]))

const { data: employees } = useQuery({
  queryKey: queryKey,
  queryFn: () => {
    return employeesService.getStoreEmployees({...filter.value})
  },
})

function updateFilter(updatedFilter: StoreEmployeesFilter) {
  filter.value = {...filter.value, ...updatedFilter}
}
</script>
