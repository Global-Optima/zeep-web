<template>
	<AdminIngredientCategoriesToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<AdminListLoader v-if="isPending" />

	<div v-else>
		<Card>
			<CardContent class="mt-4">
				<!-- No Data Message -->
				<p
					v-if="!productCategoriesResponse || productCategoriesResponse.data.length === 0"
					class="text-muted-foreground"
				>
					Категории сырья не найдены
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
