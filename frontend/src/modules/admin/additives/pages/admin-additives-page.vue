<template>
	<AdminAdditivesToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<AdminListLoader v-if="isPending" />

	<div v-else>
		<Card>
			<CardContent class="mt-4">
				<!-- No Data Message -->
				<p
					v-if="!additivesResponse || additivesResponse.data.length === 0"
					class="text-muted-foreground"
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
	</div>
</template>

<script setup lang="ts">
import AdminListLoader from '@/core/components/admin-list-loader/AdminListLoader.vue'
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Card, CardContent, CardFooter } from '@/core/components/ui/card'
import { usePaginationFilter } from '@/core/hooks/use-pagination-filter.hook'
import AdminAdditivesList from '@/modules/admin/additives/components/list/admin-additives-list.vue'
import AdminAdditivesToolbar from '@/modules/admin/additives/components/list/admin-additives-toolbar.vue'
import { additivesService } from '@/modules/admin/additives/services/additives.service'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

// Use the pagination filter hook
const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter({});

// Fetch data using Vue Query
const { data: additivesResponse, isPending } = useQuery({
  queryKey: computed(() => ['admin-additives', filter.value]),
  queryFn: () => additivesService.getAdditives(filter.value),
});
</script>

<style scoped></style>
