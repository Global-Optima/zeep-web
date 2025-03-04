<template>
	<AdminStoresToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<AdminListLoader v-if="isPending" />

	<div v-else>
		<Card>
			<CardContent class="mt-4">
				<p
					v-if="!stores || stores.data.length === 0"
					class="text-muted-foreground"
				>
					Склады не найдены
				</p>
				<AdminStoresList
					v-else
					:stores="stores.data"
				/>
			</CardContent>
			<CardFooter class="flex justify-end">
				<PaginationWithMeta
					v-if="stores"
					:meta="stores.pagination"
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
import AdminStoresList from '@/modules/admin/stores/components/list/admin-stores-list.vue'
import AdminStoresToolbar from '@/modules/admin/stores/components/list/admin-stores-toolbar.vue'
import type { StoresFilter } from '@/modules/admin/stores/models/stores-dto.model'
import { storesService } from '@/modules/admin/stores/services/stores.service'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter<StoresFilter>({})

// Query stores data
const { data: stores, isPending } = useQuery({
  queryKey: computed(() => ['stores', filter.value]),
  queryFn: () => storesService.getPaginated(filter.value),
})
</script>

<style scoped></style>
