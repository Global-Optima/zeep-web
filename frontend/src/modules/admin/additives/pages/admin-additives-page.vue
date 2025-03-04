<template>
	<AdminAdditivesToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>
	<Card>
		<CardContent class="mt-4">
			<!-- Loading Indicator -->
			<PageLoader v-if="isPending" />

			<!-- No Data Message -->
			<p
				v-else-if="!additivesResponse || additivesResponse.data.length === 0"
				class="text-muted-foreground text-center"
			>
				Топпинги не найдены
			</p>

			<!-- Additive List -->
			<AdminAdditivesList
				v-else
				:additives="additivesResponse.data"
			/>
		</CardContent>
		<CardFooter class="flex justify-end">
			<PaginationWithMeta
				v-if="additivesResponse"
				:meta="additivesResponse.pagination"
				@update:page="updatePage"
				@update:pageSize="updatePageSize"
			/>
		</CardFooter>
	</Card>
</template>

<script setup lang="ts">
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue';
import { Card, CardContent, CardFooter } from '@/core/components/ui/card';
import PageLoader from '@/core/components/page-loader/PageLoader.vue';
import { usePaginationFilter } from '@/core/hooks/use-pagination-filter.hook';
import { DEFAULT_PAGINATION_META } from '@/core/utils/pagination.utils';
import AdminAdditivesList from '@/modules/admin/additives/components/list/admin-additives-list.vue';
import AdminAdditivesToolbar from '@/modules/admin/additives/components/list/admin-additives-toolbar.vue';
import type { AdditiveFilterQuery } from '@/modules/admin/additives/models/additives.model';
import { additivesService } from '@/modules/admin/additives/services/additives.service';
import { useQuery } from '@tanstack/vue-query';
import { computed } from 'vue';

// Define default filter
const defaultFilter: AdditiveFilterQuery = {
  page: DEFAULT_PAGINATION_META.page,
  pageSize: DEFAULT_PAGINATION_META.pageSize,
};

// Use the pagination filter hook
const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter(defaultFilter);

// Fetch data using Vue Query
const { data: additivesResponse, isPending } = useQuery({
  queryKey: computed(() => ['admin-additives', filter.value]),
  queryFn: () => additivesService.getAdditives(filter.value),
});
</script>

<style scoped></style>
