<template>
	<AdminAdditiveCategoriesToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<AdminListLoader v-if="isPending" />

	<div v-else>
		<Card>
			<CardContent class="mt-4">
				<!-- Loading Indicator -->

				<!-- No Data Message -->
				<p
					v-if="!categoriesResponse || categoriesResponse.data.length === 0"
					class="text-muted-foreground"
				>
					Категории модификаторов не найдены
				</p>

				<!-- Category List -->
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
	</div>
</template>

<script setup lang="ts">
import AdminListLoader from '@/core/components/admin-list-loader/AdminListLoader.vue'
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Card, CardContent, CardFooter } from '@/core/components/ui/card'
import { usePaginationFilter } from "@/core/hooks/use-pagination-filter.hook"
import AdminAdditiveCategoriesList from '@/modules/admin/additive-categories/components/list/admin-additive-categories-list.vue'
import AdminAdditiveCategoriesToolbar from '@/modules/admin/additive-categories/components/list/admin-additive-categories-toolbar.vue'
import type { AdditiveCategoriesFilterQuery } from '@/modules/admin/additives/models/additives.model'
import { additivesService } from '@/modules/admin/additives/services/additives.service'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

// Use the pagination filter hook
const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter<AdditiveCategoriesFilterQuery>({
  includeEmpty: true,
});

// Fetch data using Vue Query
const { data: categoriesResponse, isPending } = useQuery({
  queryKey: computed(() => ['admin-additive-categories', filter.value]),
  queryFn: () => additivesService.getAdditiveCategories(filter.value),
});
</script>

<style scoped></style>
