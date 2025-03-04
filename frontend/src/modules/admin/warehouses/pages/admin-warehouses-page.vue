<template>
	<AdminWarehousesToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<AdminListLoader v-if="isPending" />

	<div v-else>
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
	</div>
</template>

<script setup lang="ts">
import AdminListLoader from '@/core/components/admin-list-loader/AdminListLoader.vue'
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Card, CardContent, CardFooter } from '@/core/components/ui/card'
import { usePaginationFilter } from '@/core/hooks/use-pagination-filter.hook'
import AdminWarehousesList from '@/modules/admin/warehouses/components/list/admin-warehouses-list.vue'
import AdminWarehousesToolbar from '@/modules/admin/warehouses/components/list/admin-warehouses-toolbar.vue'
import type { WarehouseFilter } from '@/modules/admin/warehouses/models/warehouse.model'
import { warehouseService } from '@/modules/admin/warehouses/services/warehouse.service'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter<WarehouseFilter>({})

const { data: regionsResponse, isPending } = useQuery({
  queryKey: computed(() => ['admin-warehouses', filter.value]),
  queryFn: () => warehouseService.getPaginated(filter.value),
})
</script>

<style scoped></style>
