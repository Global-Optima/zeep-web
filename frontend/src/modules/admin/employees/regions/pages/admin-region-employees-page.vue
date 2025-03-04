<template>
	<AdminRegionsEmployeesToolbar
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
				<AdminRegionsEmployeesList
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
import AdminRegionsEmployeesList from '@/modules/admin/employees/regions/components/list/admin-regions-employees-list.vue'
import AdminRegionsEmployeesToolbar from '@/modules/admin/employees/regions/components/list/admin-regions-employees-toolbar.vue'
import { regionEmployeeService } from '@/modules/admin/employees/regions/services/region-employees.service'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter<EmployeesFilter>({})

const { data: employees, isPending } = useQuery({
  queryKey: computed(() => ['region-employees', filter.value]),
  queryFn: () => regionEmployeeService.getRegionEmployees(filter.value),
})
</script>

<style scoped></style>
