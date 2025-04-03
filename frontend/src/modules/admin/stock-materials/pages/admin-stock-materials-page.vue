<template>
	<AdminStockMaterialsToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<AdminListLoader v-if="isPending" />

	<div v-else>
		<Card>
			<CardContent class="mt-4">
				<p
					v-if="!stockMaterialsResponse || stockMaterialsResponse.data.length === 0"
					class="text-muted-foreground"
				>
					Складские продукты не найдены
				</p>
				<AdminStockMaterialsList
					v-else
					:stock-materials="stockMaterialsResponse.data"
				/>
			</CardContent>
			<CardFooter class="flex justify-end">
				<PaginationWithMeta
					v-if="stockMaterialsResponse"
					:meta="stockMaterialsResponse.pagination"
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
import AdminStockMaterialsList from '@/modules/admin/stock-materials/components/list/admin-stock-materials-list.vue'
import AdminStockMaterialsToolbar from '@/modules/admin/stock-materials/components/list/admin-stock-materials-toolbar.vue'
import type { StockMaterialsFilter } from '@/modules/admin/stock-materials/models/stock-materials.model'
import { stockMaterialsService } from '@/modules/admin/stock-materials/services/stock-materials.service'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter<StockMaterialsFilter>({});

const { data: stockMaterialsResponse, isPending } = useQuery({
  queryKey: computed(() => ['admin-stock-materials', filter.value]),
  queryFn: () => stockMaterialsService.getAllStockMaterials(filter.value),
})
</script>

<style scoped></style>
