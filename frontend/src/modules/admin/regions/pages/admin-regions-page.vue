<template>
	<AdminListLoader v-if="isPending" />

	<div v-else>
		<AdminRegionsToolbar
			:filter="filter"
			@update:filter="updateFilter"
		/>

		<Card>
			<CardContent class="mt-4">
				<p
					v-if="!regionsResponse || regionsResponse.data.length === 0"
					class="text-muted-foreground"
				>
					Регионы не найдены
				</p>
				<AdminRegionsList
					v-else
					:regions="regionsResponse.data"
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
import AdminRegionsList from '@/modules/admin/regions/components/list/admin-regions-list.vue'
import AdminRegionsToolbar from '@/modules/admin/regions/components/list/admin-regions-toolbar.vue'
import type { RegionFilterDTO } from '@/modules/admin/regions/models/regions.model'
import { regionsService } from '@/modules/admin/regions/services/regions.service'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter<RegionFilterDTO>({})

const { data: regionsResponse, isPending } = useQuery({
  queryKey: computed(() => ['admin-regions', filter.value]),
  queryFn: () => regionsService.getPaginated(filter.value),
})
</script>

<style scoped></style>
