<template>
	<AdminIngredientsToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<Card>
		<CardContent class="mt-4">
			<!-- Loading Indicator -->
			<PageLoader v-if="isPending" />

			<!-- No Data Message -->
			<p
				v-else-if="!ingredientsResponse || ingredientsResponse.data.length === 0"
				class="text-muted-foreground text-center h-52 flex items-center justify-center"
			>
				Ингредиенты не найдены
			</p>

			<!-- Ingredients List -->
			<AdminIngredientsList
				v-else
				:ingredients="ingredientsResponse.data"
			/>
		</CardContent>
		<CardFooter class="flex justify-end">
			<PaginationWithMeta
				v-if="ingredientsResponse"
				:meta="ingredientsResponse.pagination"
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
import AdminIngredientsList from '@/modules/admin/ingredients/components/list/admin-ingredients-list.vue'
import AdminIngredientsToolbar from '@/modules/admin/ingredients/components/list/admin-ingredients-toolbar.vue'
import type { IngredientFilter } from '@/modules/admin/ingredients/models/ingredients.model'
import { ingredientsService } from '@/modules/admin/ingredients/services/ingredients.service'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

// Use pagination filter composable
const defaultFilter: IngredientFilter = {
	page: DEFAULT_PAGINATION_META.page,
	pageSize: DEFAULT_PAGINATION_META.pageSize,
}

const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter(defaultFilter)

// Fetch data using Vue Query
const { data: ingredientsResponse, isPending } = useQuery({
	queryKey: computed(() => ['admin-ingredients', filter.value]),
	queryFn: () => ingredientsService.getIngredients(filter.value),
})
</script>

<style scoped></style>
