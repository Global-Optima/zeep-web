<template>
	<AdminIngredientsToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<Card>
		<CardContent class="mt-4">
			<p
				v-if="!ingredientsResponse || ingredientsResponse.data.length === 0"
				class="text-muted-foreground"
			>
				Ингредиенты не найдены
			</p>
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
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Card, CardContent, CardFooter } from '@/core/components/ui/card'
import { DEFAULT_PAGINATION_META } from '@/core/utils/pagination.utils'
import AdminIngredientsList from '@/modules/admin/ingredients/components/list/admin-ingredients-list.vue'
import AdminIngredientsToolbar from '@/modules/admin/ingredients/components/list/admin-ingredients-toolbar.vue'
import type { IngredientFilter } from '@/modules/admin/ingredients/models/ingredients.model'
import { ingredientsService } from '@/modules/admin/ingredients/services/ingredients.service'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const filter = ref<IngredientFilter>({})

const { data: ingredientsResponse } = useQuery({
  queryKey: computed(() => ['admin-ingredients', filter.value]),
  queryFn: () => ingredientsService.getIngredients(filter.value),
})

function updateFilter(updatedFilter: IngredientFilter) {
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
