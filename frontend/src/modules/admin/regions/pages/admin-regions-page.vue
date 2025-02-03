<template>
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
</template>

<script setup lang="ts">
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Card, CardContent, CardFooter } from '@/core/components/ui/card'
import { DEFAULT_PAGINATION_META } from '@/core/utils/pagination.utils'
import type { AdditiveFilterQuery } from '@/modules/admin/additives/models/additives.model'
import AdminRegionsList from '@/modules/admin/regions/components/list/admin-regions-list.vue'
import AdminRegionsToolbar from '@/modules/admin/regions/components/list/admin-regions-toolbar.vue'
import type { RegionFilterDTO } from '@/modules/admin/regions/models/regions.model'
import { regionsService } from '@/modules/admin/regions/services/regions.service'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const filter = ref<RegionFilterDTO>({})

const { data: regionsResponse } = useQuery({
  queryKey: computed(() => ['admin-regions', filter.value]),
  queryFn: () => regionsService.getPaginated(filter.value),
})

function updateFilter(updatedFilter: AdditiveFilterQuery) {
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
