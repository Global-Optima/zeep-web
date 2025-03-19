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
					Сотрудники региона {{ region.name }}
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
				<AdminRegionsEmployeesList
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
import AdminRegionsEmployeesList from '@/modules/admin/employees/regions/components/list/admin-regions-employees-list.vue'
import type { RegionEmployeeFilter } from '@/modules/admin/employees/regions/models/region-employees.model'
import { regionEmployeeService } from '@/modules/admin/employees/regions/services/region-employees.service'
import type { RegionDTO } from '@/modules/admin/regions/models/regions.model'
import type { WarehouseFilter } from '@/modules/admin/warehouses/models/warehouse.model'
import { useQuery } from '@tanstack/vue-query'
import { ChevronLeft } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

const {region} = defineProps<{region: RegionDTO}>()

const emits = defineEmits<{
  (e: 'onCancel'): void
}>()

const router = useRouter()

const onCancel = () => {
  emits('onCancel')
}

const filter = ref<RegionEmployeeFilter>({})

const { data: regionsResponse } = useQuery({
  queryKey: computed(() => ['admin-region-employees', region.id]),
  queryFn: () => regionEmployeeService.getRegionEmployees({...filter.value, regionId: region.id}),
})

const onAddClick = () => {
  router.push(`/admin/regions/${region.id}/employees/create`)
}

function updateFilter(updatedFilter: WarehouseFilter) {
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
