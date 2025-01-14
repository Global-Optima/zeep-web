<template>
	<AdminAdditiveCategoriesToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<Card>
		<CardContent class="mt-4">
			<p
				v-if="!categoriesResponse || categoriesResponse.data.length === 0"
				class="text-muted-foreground"
			>
				Категории топпингов не найдены
			</p>
			<AdminAdditiveCategoriesList
				v-else
				:categories="categoriesResponse.data"
			/>
		</CardContent>
		<CardFooter class="flex justify-end">
			<PaginationWithMeta
				v-if="categoriesResponse"
				:meta="categoriesResponse.pagination"
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
import AdminAdditiveCategoriesList from '@/modules/admin/additive-categories/components/list/admin-additive-categories-list.vue'
import AdminAdditiveCategoriesToolbar from '@/modules/admin/additive-categories/components/list/admin-additive-categories-toolbar.vue'
import type { AdditiveCategoriesFilterQuery } from '@/modules/admin/additives/models/additives.model'
import { additivesService } from '@/modules/admin/additives/services/additives.service'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const filter = ref<AdditiveCategoriesFilterQuery>({
  showAll: true
})

const { data: categoriesResponse } = useQuery({
  queryKey: computed(() => ['admin-additive-categories', filter.value]),
  queryFn: () => additivesService.getAdditiveCategories(filter.value),
})

function updateFilter(updatedFilter: AdditiveCategoriesFilterQuery) {
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
