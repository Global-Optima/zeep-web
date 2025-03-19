<template>
	<div class="flex items-center gap-4 mb-4">
		<Button
			variant="outline"
			size="icon"
			@click="$router.back"
		>
			<ChevronLeft class="w-5 h-5" />
			<span class="sr-only">Назад</span>
		</Button>
		<h1 class="flex-1 sm:grow-0 font-semibold text-xl tracking-tight whitespace-nowrap shrink-0">
			Аудит сотрудника
		</h1>
	</div>

	<AdminStoreEmployeesAuditToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<AdminListLoader v-if="isPending" />

	<div v-else>
		<Card class="mt-4">
			<CardContent class="p-4">
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
	</div>
</template>

<script setup lang="ts">
import AdminListLoader from '@/core/components/admin-list-loader/AdminListLoader.vue'
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardFooter } from '@/core/components/ui/card'
import { usePaginationFilter } from '@/core/hooks/use-pagination-filter.hook'
import AdminStoreEmployeesAuditList from '@/modules/admin/employees/components/audit/admin-store-employees-audit-list.vue'
import AdminStoreEmployeesAuditToolbar from '@/modules/admin/employees/components/audit/admin-store-employees-audit-toolbar.vue'
import type { EmployeeAuditFilter } from '@/modules/admin/employees/models/employees-audit.models'
import { employeeAuditService } from '@/modules/admin/employees/services/employees-audit.service'
import { useQuery } from '@tanstack/vue-query'
import { ChevronLeft } from 'lucide-vue-next'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter<EmployeeAuditFilter>({})

const route = useRoute()
const employeeId = route.params.id as string

const { data: employeeAudits, isPending } = useQuery({
  queryKey: computed(() => ['employee-audits', {...filter.value, employeeId}]),
	queryFn: () => employeeAuditService.getAudits({...filter.value, employeeId: Number(employeeId)}),
  enabled: !!employeeId,
})
</script>

<style scoped></style>
