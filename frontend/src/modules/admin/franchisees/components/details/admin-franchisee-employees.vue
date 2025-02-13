<template>
	<div class="flex-1 gap-4 grid auto-rows-max mx-auto max-w-6xl">
		<div class="flex justify-between items-center gap-4">
			<div class="flex items-center gap-4">
				<Button
					variant="outline"
					size="icon"
					@click="onCancel"
				>
					<ChevronLeft class="w-5 h-5" />
					<span class="sr-only">Назад</span>
				</Button>
				<h1
					class="flex-1 sm:grow-0 font-semibold text-xl tracking-tight whitespace-nowrap shrink-0"
				>
					Сотрудники франчайзи {{ franchisee.name }}
				</h1>
			</div>

			<div>
				<Button @click="onAddClick"> Создать </Button>
			</div>
		</div>

		<Card>
			<CardContent class="mt-4">
				<p
					v-if="!regionsResponse || regionsResponse.data.length === 0"
					class="text-muted-foreground"
				>
					Сотрудники региона не найдены
				</p>
				<AdminFranchiseesEmployeesList
					v-else
					:employees="regionsResponse.data"
				/>
			</CardContent>
			<CardFooter class="flex justify-end">
				<PaginationWithMeta
					v-if="regionsResponse"
					:meta="regionsResponse.pagination"
					@update:page="updatePage"
					@update:pageSize="updatePageSize"
				/>
			</CardFooter>
		</Card>
	</div>
</template>

<script setup lang="ts">
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardFooter } from '@/core/components/ui/card'
import { DEFAULT_PAGINATION_META } from '@/core/utils/pagination.utils'
import AdminFranchiseesEmployeesList from '@/modules/admin/employees/franchisees/components/list/admin-franchisees-employees-list.vue'
import type { FranchiseEmployeeFilter } from '@/modules/admin/employees/franchisees/models/franchisees-employees.model'
import { franchiseeEmployeeService } from '@/modules/admin/employees/franchisees/services/franchisee-employees.service'
import type { FranchiseeDTO } from '@/modules/admin/franchisees/models/franchisee.model'
import { useQuery } from '@tanstack/vue-query'
import { ChevronLeft } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

const {franchisee} = defineProps<{franchisee: FranchiseeDTO}>()

const emits = defineEmits<{
  (e: 'onCancel'): void
}>()

const router = useRouter()

const onCancel = () => {
  emits('onCancel')
}

const filter = ref<FranchiseEmployeeFilter>({})

const { data: regionsResponse } = useQuery({
  queryKey: computed(() => ['admin-franchisee-employees', franchisee.id]),
  queryFn: () => franchiseeEmployeeService.getFranchiseeEmployees({...filter.value, franchiseeId: franchisee.id}),
})

const onAddClick = () => {
  router.push(`/admin/franchisees/${franchisee.id}/employees/create`)
}

function updateFilter(updatedFilter: FranchiseEmployeeFilter) {
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
