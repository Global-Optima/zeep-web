<template>
	<AdminStoreProductsToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<Card>
		<CardContent class="mt-4">
			<p
				v-if="!storeProductsResponse || storeProductsResponse.data.length === 0"
				class="text-muted-foreground"
			>
				Товары магазина не найдены
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
</template>

<script setup lang="ts">
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Card, CardContent } from '@/core/components/ui/card'
import CardFooter from '@/core/components/ui/card/CardFooter.vue'
import { DEFAULT_PAGINATION_META } from '@/core/utils/pagination.utils'
import AdminStoreProductsList from '@/modules/admin/store-products/components/list/admin-store-products-list.vue'
import AdminStoreProductsToolbar from '@/modules/admin/store-products/components/list/admin-store-products-toolbar.vue'
import type { StoreProductsFilterDTO } from '@/modules/admin/store-products/models/store-products.model'
import { storeProductsService } from '@/modules/admin/store-products/services/store-products.service'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const filter = ref<StoreProductsFilterDTO>({})

const { data: storeProductsResponse } = useQuery({
  queryKey: computed(() => ['admin-store-products', filter.value]),
  queryFn: () => storeProductsService.getStoreProducts(filter.value),
})

function updateFilter(updatedFilter: StoreProductsFilterDTO) {
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
