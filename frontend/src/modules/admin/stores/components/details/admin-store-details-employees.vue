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
					Сотрудники кафе {{ store.name }}
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
					Сотрудники кафе не найдены
				</p>
				<AdminStoreEmployeesList
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
import AdminStoreEmployeesList from '@/modules/admin/employees/stores/components/list/admin-store-employees-list.vue'
import type { StoreEmployeeFilter } from '@/modules/admin/employees/stores/models/store-employees.model'
import { storeEmployeeService } from '@/modules/admin/employees/stores/services/store-employees.service'
import type { StoreDTO } from '@/modules/admin/stores/models/stores.models'
import { useQuery } from '@tanstack/vue-query'
import { ChevronLeft } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

const {store} = defineProps<{store: StoreDTO}>()

const emits = defineEmits<{
  (e: 'onCancel'): void
}>()

const router = useRouter()

const filter = ref<StoreEmployeeFilter>({})

const { data: regionsResponse } = useQuery({
  queryKey: computed(() => ['admin-store-employees', store.id]),
  queryFn: () => storeEmployeeService.getStoreEmployees({...filter.value, storeId: store.id}),
})

const onCancel = () => {
  emits('onCancel')
}

const onAddClick = () => {
  router.push(`/admin/stores/${store.id}/employees/create`)
}

function updateFilter(updatedFilter: StoreEmployeeFilter) {
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
