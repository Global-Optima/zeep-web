<template>
	<AdminStoreProvisionsToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<AdminListLoader v-if="isPending" />

	<div v-else>
		<Card>
			<CardContent class="mt-4">
				<!-- No Data Message -->
				<p
					v-if="!provisionsResponse || provisionsResponse.data.length === 0"
					class="text-muted-foreground"
				>
					Заготовки не найдены
				</p>

				<!-- Ingredients List -->
				<AdminStoreProvisionsList
					v-else
					:storeProvisions="provisionsResponse.data"
				/>
			</CardContent>
			<CardFooter class="flex justify-end">
				<PaginationWithMeta
					v-if="provisionsResponse"
					:meta="provisionsResponse.pagination"
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
import AdminStoreProvisionsList from '@/modules/admin/store-provisions/components/list/admin-store-provisions-list.vue'
import AdminStoreProvisionsToolbar from '@/modules/admin/store-provisions/components/list/admin-store-provisions-toolbar.vue'
import type { StoreProvisionFilter } from '@/modules/admin/store-provisions/models/store-provision.models'
import { storeProvisionsService } from '@/modules/admin/store-provisions/services/store-provision.service'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter<StoreProvisionFilter>({})

// Fetch data using Vue Query
const { data: provisionsResponse, isPending } = useQuery({
  queryKey: computed(() => ['admin-store-provisions', filter.value]),
  queryFn: () => storeProvisionsService.getStoreProvisions(filter.value),
})
</script>

<style scoped></style>
