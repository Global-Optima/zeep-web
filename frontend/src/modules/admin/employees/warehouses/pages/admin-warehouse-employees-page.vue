<template>
	<AdminWarehousesEmployeesToolbar
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
			<AdminWarehousesEmployeesList
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
import type { EmployeesFilter } from '@/modules/admin/employees/models/employees.models'
import AdminWarehousesEmployeesList from '@/modules/admin/employees/warehouses/components/list/admin-warehouses-employees-list.vue'
import AdminWarehousesEmployeesToolbar from '@/modules/admin/employees/warehouses/components/list/admin-warehouses-employees-toolbar.vue'
import { warehouseEmployeeService } from '@/modules/admin/employees/warehouses/services/warehouse-employees.service'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const filter = ref<EmployeesFilter>({})

const { data: employees } = useQuery({
  queryKey: computed(() => ['warehouse-employees', filter.value]),
  queryFn: () => warehouseEmployeeService.getWarehouseEmployees(filter.value),
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
