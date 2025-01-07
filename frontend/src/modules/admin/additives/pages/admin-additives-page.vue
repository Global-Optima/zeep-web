<template>
	<AdminAdditivesToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<Card>
		<CardContent class="mt-4">
			<p
				v-if="!additivesResponse || additivesResponse.data.length === 0"
				class="text-muted-foreground"
			>
				Топпингиы не найдены
			</p>
			<AdminAdditivesList
				v-else
				:additives="additivesResponse.data"
			/>
		</CardContent>
		<CardFooter class="flex justify-end">
			<PaginationWithMeta
				v-if="additivesResponse"
				:meta="additivesResponse.pagination"
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
import AdminAdditivesList from '@/modules/admin/additives/components/list/admin-additives-list.vue'
import AdminAdditivesToolbar from '@/modules/admin/additives/components/list/admin-additives-toolbar.vue'
import type { AdditiveFilterQuery } from '@/modules/admin/additives/models/additives.model'
import { additivesService } from '@/modules/admin/additives/services/additives.service'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const filter = ref<AdditiveFilterQuery>({
})

const { data: additivesResponse } = useQuery({
  queryKey: computed(() => ['admin-additives', filter.value]),
  queryFn: () => additivesService.getAdditives(filter.value),
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
