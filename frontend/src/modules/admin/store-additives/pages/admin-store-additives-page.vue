<template>
	<AdminStoreAdditivesToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<Card>
		<CardContent class="mt-4">
			<p
				v-if="!storeAdditivesResponse || storeAdditivesResponse.data.length === 0"
				class="text-muted-foreground"
			>
				Топпинги магазина не найдены
			</p>

			<AdminStoreAdditivesList
				v-else
				:additives="storeAdditivesResponse.data"
			/>
		</CardContent>
		<CardFooter class="flex justify-end">
			<PaginationWithMeta
				v-if="storeAdditivesResponse"
				:meta="storeAdditivesResponse.pagination"
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
import AdminStoreAdditivesList from '@/modules/admin/store-additives/components/list/admin-store-additives-list.vue'
import AdminStoreAdditivesToolbar from '@/modules/admin/store-additives/components/list/admin-store-additives-toolbar.vue'
import type { StoreAdditivesFilterDTO } from '@/modules/admin/store-additives/models/store-additves.model'
import { storeAdditivesService } from '@/modules/admin/store-additives/services/store-additives.service'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const filter = ref<StoreAdditivesFilterDTO>({})

const { data: storeAdditivesResponse } = useQuery({
  queryKey: computed(() => ['admin-store-additives', filter.value]),
  queryFn: () => storeAdditivesService.getStoreAdditives(filter.value),
})

function updateFilter(updatedFilter: StoreAdditivesFilterDTO) {
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
