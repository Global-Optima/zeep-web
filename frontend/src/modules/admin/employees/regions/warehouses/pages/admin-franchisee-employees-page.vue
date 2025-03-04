<template>
	<AdminListLoader v-if="isPending" />

	<div v-else>
		<AdminFranchiseesEmployeesToolbar
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
				<AdminFranchiseesEmployeesList
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
import AdminFranchiseesEmployeesList from '@/modules/admin/employees/franchisees/components/list/admin-franchisees-employees-list.vue'
import AdminFranchiseesEmployeesToolbar from '@/modules/admin/employees/franchisees/components/list/admin-franchisees-employees-toolbar.vue'
import { franchiseeEmployeeService } from '@/modules/admin/employees/franchisees/services/franchisee-employees.service'
import type { EmployeesFilter } from '@/modules/admin/employees/models/employees.models'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter<EmployeesFilter>({})

const { data: employees, isPending } = useQuery({
  queryKey: computed(() => ['franchisee-employees', filter.value]),
  queryFn: () => franchiseeEmployeeService.getFranchiseeEmployees(filter.value),
})
</script>

<style scoped></style>
