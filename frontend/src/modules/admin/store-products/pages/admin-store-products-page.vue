<template>
	<div
		v-if=" !isStoreSyncPending && isStoreSyncResponse && !isStoreSyncResponse.isSync"
		class="mb-4"
	>
		<AdminStoresSyncCard />
	</div>

	<AdminStoreProductsToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<AdminListLoader v-if="isLoading" />

	<div v-else>
		<Card>
			<CardContent class="mt-4">
				<p
					v-if="!storeProductsResponse || storeProductsResponse.data.length === 0"
					class="text-muted-foreground"
				>
					Товары кафе не найдены
				</p>

				<AdminStoreProductsList
					v-else
					:storeProducts="storeProductsResponse.data"
				/>
			</CardContent>
			<CardFooter class="flex justify-end">
				<PaginationWithMeta
					v-if="storeProductsResponse"
					:meta="storeProductsResponse.pagination"
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
import { Card, CardContent } from '@/core/components/ui/card'
import CardFooter from '@/core/components/ui/card/CardFooter.vue'
import { usePaginationFilter } from '@/core/hooks/use-pagination-filter.hook'
import AdminStoreProductsList from '@/modules/admin/store-products/components/list/admin-store-products-list.vue'
import AdminStoreProductsToolbar from '@/modules/admin/store-products/components/list/admin-store-products-toolbar.vue'
import AdminStoresSyncCard from '@/modules/admin/store-products/components/stores-sync/admin-stores-sync-card.vue'
import type { StoreProductsFilterDTO } from '@/modules/admin/store-products/models/store-products.model'
import { storeProductsService } from '@/modules/admin/store-products/services/store-products.service'
import { storeSyncService } from '@/modules/admin/stores/services/stores-sync.service'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter<StoreProductsFilterDTO>({})

const { data: storeProductsResponse, isLoading } = useQuery({
  queryKey: computed(() => ['admin-store-products', filter.value]),
  queryFn: () => storeProductsService.getStoreProducts(filter.value),
  enabled: computed(() => Boolean(filter.value.storeId))
})

const { data: isStoreSyncResponse, isLoading: isStoreSyncPending } = useQuery({
  queryKey: ['admin-store-is-sync'],
  queryFn: () => storeSyncService.isStoreSynchronized(),
  enabled: computed(() => Boolean(filter.value.storeId))
})
</script>

<style scoped></style>
