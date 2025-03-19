<template>
	<AdminIngredientsToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<AdminListLoader v-if="isPending" />

	<div v-else>
		<Card>
			<CardContent class="mt-4">
				<!-- No Data Message -->
				<p
					v-if="!ingredientsResponse || ingredientsResponse.data.length === 0"
					class="text-muted-foreground text-center"
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
	</div>
</template>

<script setup lang="ts">
import AdminListLoader from '@/core/components/admin-list-loader/AdminListLoader.vue'
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Card, CardContent, CardFooter } from '@/core/components/ui/card'
import { usePaginationFilter } from "@/core/hooks/use-pagination-filter.hook"
import AdminIngredientsList from '@/modules/admin/ingredients/components/list/admin-ingredients-list.vue'
import AdminIngredientsToolbar from '@/modules/admin/ingredients/components/list/admin-ingredients-toolbar.vue'
import { ingredientsService } from '@/modules/admin/ingredients/services/ingredients.service'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter({})

// Fetch data using Vue Query
const { data: ingredientsResponse, isPending } = useQuery({
	queryKey: computed(() => ['admin-ingredients', filter.value]),
	queryFn: () => ingredientsService.getIngredients(filter.value),
})
</script>

<style scoped></style>
