<template>
	<AdminAdminEmployeesToolbar
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
				<AdminAdminEmployeesList
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
import AdminAdminEmployeesList from '@/modules/admin/employees/admins/components/list/admin-admin-employees-list.vue'
import AdminAdminEmployeesToolbar from '@/modules/admin/employees/admins/components/list/admin-admin-employees-toolbar.vue'
import { adminEmployeeService } from '@/modules/admin/employees/admins/services/admin-employees.service'
import type { EmployeesFilter } from '@/modules/admin/employees/models/employees.models'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter<EmployeesFilter>({})

const { data: employees, isPending } = useQuery({
  queryKey: computed(() => ['admin-employees', filter.value]),
  queryFn: () => adminEmployeeService.getAdminEmployees(filter.value),
})
</script>

<style scoped></style>
