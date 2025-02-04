<template>
	<AdminWarehousesToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<Card>
		<CardContent class="mt-4">
			<p
				v-if="!regionsResponse || regionsResponse.data.length === 0"
				class="text-muted-foreground"
			>
				Склады не найдены
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
</template>

<script setup lang="ts">
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Card, CardContent, CardFooter } from '@/core/components/ui/card'
import { DEFAULT_PAGINATION_META } from '@/core/utils/pagination.utils'
import AdminRegionsToolbar from '@/modules/admin/regions/components/list/admin-regions-toolbar.vue'
import AdminWarehousesList from '@/modules/admin/warehouses/components/list/admin-warehouses-list.vue'
import AdminWarehousesToolbar from '@/modules/admin/warehouses/components/list/admin-warehouses-toolbar.vue'
import type { WarehouseFilter } from '@/modules/admin/warehouses/models/warehouse.model'
import { warehouseService } from '@/modules/admin/warehouses/services/warehouse.service'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const filter = ref<WarehouseFilter>({})

const { data: regionsResponse } = useQuery({
  queryKey: computed(() => ['admin-warehouses', filter.value]),
  queryFn: () => warehouseService.getPaginated(filter.value),
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
