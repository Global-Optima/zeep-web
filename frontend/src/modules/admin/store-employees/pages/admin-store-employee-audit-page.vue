<template>
	<div
		v-if="employee"
		class="mx-auto w-full max-w-7xl"
	>
		<div class="flex items-center gap-4">
			<Button
				variant="outline"
				size="icon"
				@click="$router.back"
			>
				<ChevronLeft class="w-5 h-5" />
				<span class="sr-only">Назад</span>
			</Button>
			<h1 class="flex-1 sm:grow-0 font-semibold text-xl tracking-tight whitespace-nowrap shrink-0">
				Активности сотрудника {{ getEmployeeShortName(employee) }}
			</h1>

			<div class="md:flex items-center gap-2 hidden md:ml-auto">
				<Button
					variant="outline"
					type="button"
					@click="$router.back"
				>
					Отменить
				</Button>
			</div>
		</div>

		<Card class="mt-4">
			<CardContent class="p-4">
				<AdminStoreEmployeesAuditToolbar
					:filter="filter"
					@update:filter="updateFilter"
				/>

				<Card>
					<CardContent class="mt-4">
						<p
							v-if="!employeeAudits || employeeAudits.data.length === 0"
							class="text-muted-foreground"
						>
							Активности не найдены
						</p>
						<AdminStoreEmployeesAuditList
							v-else
							:audits="employeeAudits.data"
						/>
					</CardContent>
					<CardFooter class="flex justify-end">
						<PaginationWithMeta
							v-if="employeeAudits"
							:meta="employeeAudits.pagination"
							@update:page="updatePage"
							@update:pageSize="updatePageSize"
						/>
					</CardFooter>
				</Card>
			</CardContent>
		</Card>
	</div>
</template>

<script setup lang="ts">
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardFooter } from '@/core/components/ui/card'
import { DEFAULT_PAGINATION_META } from '@/core/utils/pagination.utils'
import { getEmployeeShortName } from '@/core/utils/user-formatting.utils'
import AdminStoreEmployeesAuditList from '@/modules/admin/store-employees/components/audit/admin-store-employees-audit-list.vue'
import AdminStoreEmployeesAuditToolbar from '@/modules/admin/store-employees/components/audit/admin-store-employees-audit-toolbar.vue'
import type { EmployeeAuditFilter } from '@/modules/admin/store-employees/models/employees-audit.models'
import { employeeAuditService } from '@/modules/admin/store-employees/services/employees-audit.service'
import { employeesService } from '@/modules/admin/store-employees/services/employees.service'
import { useQuery } from '@tanstack/vue-query'
import { ChevronLeft } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useRoute } from 'vue-router'

const filter = ref<EmployeeAuditFilter>({})

const route = useRoute()
const employeeId = route.params.id as string

const { data: employee } = useQuery({
  queryKey: ['store-employee', employeeId],
	queryFn: () => employeesService.getStoreEmployeeById(Number(employeeId)),
  enabled: !!employeeId,
})

const { data: employeeAudits } = useQuery({
  queryKey: computed(() => ['employee-audits', {...filter.value, employeeId}]),
	queryFn: () => employeeAuditService.getAudits({...filter.value, employeeId: Number(employeeId)}),
  enabled: !!employeeId,
})

function updateFilter(updatedFilter: EmployeeAuditFilter) {
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
