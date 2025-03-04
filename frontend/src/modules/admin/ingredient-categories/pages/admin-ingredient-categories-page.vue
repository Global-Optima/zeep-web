<template>
	<AdminListLoader v-if="isPending" />

	<div v-else>
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
					class="flex justify-center items-center h-52 text-muted-foreground text-center"
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
	</div>
</template>

<script setup lang="ts">
import AdminListLoader from '@/core/components/admin-list-loader/AdminListLoader.vue'
import PageLoader from "@/core/components/page-loader/PageLoader.vue"
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Card, CardContent, CardFooter } from '@/core/components/ui/card'
import { usePaginationFilter } from "@/core/hooks/use-pagination-filter.hook"
import AdminIngredientCategoriesList from '@/modules/admin/ingredient-categories/components/list/admin-ingredient-categories-list.vue'
import AdminIngredientCategoriesToolbar from '@/modules/admin/ingredient-categories/components/list/admin-ingredient-categories-toolbar.vue'
import { ingredientsService } from '@/modules/admin/ingredients/services/ingredients.service'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'


const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter({})

// Fetch data using Vue Query
const { data: productCategoriesResponse, isPending } = useQuery({
	queryKey: computed(() => ['admin-ingredient-categories', filter.value]),
	queryFn: () => ingredientsService.getIngredientCategories(filter.value),
})
</script>

<style scoped></style>
