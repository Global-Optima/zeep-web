<template>
	<div class="flex-1 gap-4 grid auto-rows-max mx-auto max-w-6xl">
		<div class="flex items-center gap-4">
			<Button
				variant="outline"
				size="icon"
				@click="onCancel"
			>
				<ChevronLeft class="w-5 h-5" />
				<span class="sr-only">Назад</span>
			</Button>
			<h1 class="flex-1 sm:grow-0 font-semibold text-xl tracking-tight whitespace-nowrap shrink-0">
				Склады региона {{ region.name }}
			</h1>
		</div>

		<Card>
			<CardContent class="mt-4">
				<p
					v-if="!regionsResponse || regionsResponse.data.length === 0"
					class="text-muted-foreground"
				>
					Склады региона не найдены
				</p>
				<AdminWarehousesList
					v-else
					:warehouses="regionsResponse.data"
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
import type { RegionDTO } from '@/modules/admin/regions/models/regions.model'
import AdminWarehousesList from '@/modules/admin/warehouses/components/list/admin-warehouses-list.vue'
import type { WarehouseFilter } from '@/modules/admin/warehouses/models/warehouse.model'
import { warehouseService } from '@/modules/admin/warehouses/services/warehouse.service'
import { useQuery } from '@tanstack/vue-query'
import { ChevronLeft } from 'lucide-vue-next'
import { computed, ref } from 'vue'

const {region} = defineProps<{region: RegionDTO}>()

const emits = defineEmits<{
  (e: 'onCancel'): void
}>()

const onCancel = () => {
  emits('onCancel')
}

const filter = ref<WarehouseFilter>({})

const { data: regionsResponse } = useQuery({
  queryKey: computed(() => ['admin-region-warehouses', region.id]),
  queryFn: () => warehouseService.getPaginated({...filter.value, regionId: region.id}),
})

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
