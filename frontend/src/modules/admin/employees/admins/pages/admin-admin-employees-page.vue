<template>
	<AdminAdminEmployeesToolbar
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
</template>

<script setup lang="ts">
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Card, CardContent, CardFooter } from '@/core/components/ui/card'
import { DEFAULT_PAGINATION_META } from '@/core/utils/pagination.utils'
import AdminAdminEmployeesList from '@/modules/admin/employees/admins/components/list/admin-admin-employees-list.vue'
import AdminAdminEmployeesToolbar from '@/modules/admin/employees/admins/components/list/admin-admin-employees-toolbar.vue'
import { adminEmployeeService } from '@/modules/admin/employees/admins/services/admin-employees.service'
import type { EmployeesFilter } from '@/modules/admin/employees/models/employees.models'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const filter = ref<EmployeesFilter>({})

const { data: employees } = useQuery({
  queryKey: computed(() => ['admin-employees', filter.value]),
  queryFn: () => adminEmployeeService.getAdminEmployees(filter.value),
})

function updateFilter(updatedFilter: EmployeesFilter) {
  filter.value = {...filter.value, ...updatedFilter}
}

function updatePage(page: number) {
  updateFilter({ pageSize: DEFAULT_PAGINATION_META.pageSize, page: page})

}

function updatePageSize(pageSize: number) {
  updateFilter({ pageSize: pageSize, page: DEFAULT_PAGINATION_META.page})
}
</script>

<style scoped></style>
