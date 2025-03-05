<template>
	<AdminIngredientCategoriesToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<Card>
		<CardContent class="mt-4">
			<!-- Loading Indicator -->
			<PageLoader v-if="isPending" />

			<!-- No Data Message -->
			<p
				v-else-if="!productCategoriesResponse || productCategoriesResponse.data.length === 0"
				class="text-muted-foreground text-center h-52 flex items-center justify-center"
			>
				Категории ингредиентов не найдены
			</p>

			<!-- Ingredient Categories List -->
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
import PageLoader from "@/core/components/page-loader/PageLoader.vue"
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Card, CardContent, CardFooter } from '@/core/components/ui/card'
import { usePaginationFilter } from "@/core/hooks/use-pagination-filter.hook"
import { DEFAULT_PAGINATION_META } from '@/core/utils/pagination.utils'
import AdminIngredientCategoriesList from '@/modules/admin/ingredient-categories/components/list/admin-ingredient-categories-list.vue'
import AdminIngredientCategoriesToolbar from '@/modules/admin/ingredient-categories/components/list/admin-ingredient-categories-toolbar.vue'
import type { IngredientCategoryFilter } from '@/modules/admin/ingredients/models/ingredients.model'
import { ingredientsService } from '@/modules/admin/ingredients/services/ingredients.service'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

// Use pagination filter composable
const defaultFilter: IngredientCategoryFilter = {
	page: DEFAULT_PAGINATION_META.page,
	pageSize: DEFAULT_PAGINATION_META.pageSize,
}

const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter(defaultFilter)

// Fetch data using Vue Query
const { data: productCategoriesResponse, isPending } = useQuery({
	queryKey: computed(() => ['admin-ingredient-categories', filter.value]),
	queryFn: () => ingredientsService.getIngredientCategories(filter.value),
})
</script>

<style scoped></style>
