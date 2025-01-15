<template>
	<AdminIngredientCategoriesToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<Card>
		<CardContent class="mt-4">
			<p
				v-if="!productCategoriesResponse || productCategoriesResponse.data.length === 0"
				class="text-muted-foreground"
			>
				Категории ингредиентов не найдены
			</p>
			<AdminIngredientCategoriesList
				v-else
				:ingredient-categories="productCategoriesResponse.data"
			/>
		</CardContent>
		<CardFooter class="flex justify-end">
			<PaginationWithMeta
				v-if="productCategoriesResponse"
				:meta="productCategoriesResponse.pagination"
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
import AdminIngredientCategoriesList from '@/modules/admin/ingredient-categories/components/list/admin-ingredient-categories-list.vue'
import AdminIngredientCategoriesToolbar from '@/modules/admin/ingredient-categories/components/list/admin-ingredient-categories-toolbar.vue'
import type { IngredientCategoryFilter } from '@/modules/admin/ingredients/models/ingredients.model'
import { ingredientsService } from '@/modules/admin/ingredients/services/ingredients.service'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const filter = ref<IngredientCategoryFilter>({})

const { data: productCategoriesResponse } = useQuery({
  queryKey: computed(() => ['admin-ingredient-categories', filter.value]),
  queryFn: () => ingredientsService.getIngredientCategories(filter.value),
})

function updateFilter(updatedFilter: IngredientCategoryFilter) {
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
