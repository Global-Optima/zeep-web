<template>
	<AdminStoreAdditivesToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<AdminListLoader v-if="isLoading" />

	<div v-else>
		<Card>
			<CardContent class="mt-4">
				<p
					v-if="!storeAdditivesResponse || storeAdditivesResponse.data.length === 0"
					class="text-muted-foreground"
				>
					Топпинги кафе не найдены
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
	</div>
</template>

<script setup lang="ts">
import AdminListLoader from '@/core/components/admin-list-loader/AdminListLoader.vue'
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Card, CardContent } from '@/core/components/ui/card'
import CardFooter from '@/core/components/ui/card/CardFooter.vue'
import { usePaginationFilter } from '@/core/hooks/use-pagination-filter.hook'
import type { AdditiveFilterQuery } from '@/modules/admin/additives/models/additives.model'
import AdminStoreAdditivesList from '@/modules/admin/store-additives/components/list/admin-store-additives-list.vue'
import AdminStoreAdditivesToolbar from '@/modules/admin/store-additives/components/list/admin-store-additives-toolbar.vue'
import { storeAdditivesService } from '@/modules/admin/store-additives/services/store-additives.service'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter<AdditiveFilterQuery>({})

const { data: storeAdditivesResponse, isLoading } = useQuery({
  queryKey: computed(() => ['admin-store-additives', filter.value]),
  queryFn: () => storeAdditivesService.getStoreAdditives(filter.value),
  enabled: computed(() => Boolean(filter.value.storeId))
})
</script>

<style scoped></style>
