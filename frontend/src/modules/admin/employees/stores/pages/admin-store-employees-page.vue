<template>
	<AdminStoreEmployeesToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<AdminListLoader v-if="isPending" />

	<div v-else>
		<Card>
			<CardContent class="mt-4">
				<p
					v-if="!employees || employees.data.length === 0"
					class="text-muted-foreground"
				>
					Сотрудники не найдены
				</p>
				<AdminStoreEmployeesList
					v-else
					:employees="employees.data"
				/>
			</CardContent>
			<CardFooter class="flex justify-end">
				<PaginationWithMeta
					v-if="employees"
					:meta="employees.pagination"
					@update:page="updatePage"
					@update:pageSize="updatePageSize"
				/>
			</CardFooter>
		</Card>
	</div>
</template>

<script setup lang="ts">
import AdminListLoader from '@/core/components/admin-list-loader/AdminListLoader.vue'
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Card, CardContent, CardFooter } from '@/core/components/ui/card'
import { usePaginationFilter } from '@/core/hooks/use-pagination-filter.hook'
import type { EmployeesFilter } from '@/modules/admin/employees/models/employees.models'
import AdminStoreEmployeesList from '@/modules/admin/employees/stores/components/list/admin-store-employees-list.vue'
import AdminStoreEmployeesToolbar from '@/modules/admin/employees/stores/components/list/admin-store-employees-toolbar.vue'
import { storeEmployeeService } from '@/modules/admin/employees/stores/services/store-employees.service'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter<EmployeesFilter>({})

const { data: employees, isPending } = useQuery({
  queryKey: computed(() => ['store-employees', filter.value]),
  queryFn: () => storeEmployeeService.getStoreEmployees(filter.value),
})
</script>

<style scoped></style>
