<template>
	<AdminAdditiveCategoriesToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>
	<Card>
		<CardContent class="mt-4">
			<!-- Loading Indicator -->
			<PageLoader v-if="isPending" />

			<!-- No Data Message -->
			<p
				v-else-if="!categoriesResponse || categoriesResponse.data.length === 0"
				class="text-muted-foreground text-center h-52 flex items-center justify-center"
			>
				Категории топпингов не найдены
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
</template>

<script setup lang="ts">
import PageLoader from "@/core/components/page-loader/PageLoader.vue";
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue';
import { Card, CardContent, CardFooter } from '@/core/components/ui/card';
import { usePaginationFilter } from "@/core/hooks/use-pagination-filter.hook";
import AdminAdditiveCategoriesList from '@/modules/admin/additive-categories/components/list/admin-additive-categories-list.vue';
import AdminAdditiveCategoriesToolbar from '@/modules/admin/additive-categories/components/list/admin-additive-categories-toolbar.vue';
import type { AdditiveCategoriesFilterQuery } from '@/modules/admin/additives/models/additives.model';
import { additivesService } from '@/modules/admin/additives/services/additives.service';
import { useQuery } from '@tanstack/vue-query';
import { computed } from 'vue';

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
