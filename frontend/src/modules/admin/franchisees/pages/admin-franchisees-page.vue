<template>
	<AdminFranchiseesToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<Card>
		<CardContent class="mt-4">
			<PageLoader v-if="isPending" />

			<p
				v-else-if="!franchiseesResponse || franchiseesResponse.data.length === 0"
				class="text-muted-foreground text-center h-52 flex items-center justify-center"
			>
				Франчайзи не найдены
			</p>

			<AdminFranchiseesList
				v-else
				:franchisees="franchiseesResponse.data"
			/>
		</CardContent>
		<CardFooter class="flex justify-end">
			<PaginationWithMeta
				v-if="franchiseesResponse"
				:meta="franchiseesResponse.pagination"
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
import AdminFranchiseesList from '@/modules/admin/franchisees/components/list/admin-franchisees-list.vue'
import AdminFranchiseesToolbar from '@/modules/admin/franchisees/components/list/admin-franchisees-toolbar.vue'
import type { FranchiseeFilterDTO } from '@/modules/admin/franchisees/models/franchisee.model'
import { franchiseeService } from '@/modules/admin/franchisees/services/franchisee.service'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

// Use pagination filter composable
const defaultFilter: FranchiseeFilterDTO = {
	page: DEFAULT_PAGINATION_META.page,
	pageSize: DEFAULT_PAGINATION_META.pageSize,
}

const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter(defaultFilter)

// Fetch data using Vue Query
const { data: franchiseesResponse, isPending } = useQuery({
	queryKey: computed(() => ['admin-franchisees', filter.value]),
	queryFn: () => franchiseeService.getPaginated(filter.value),
})
</script>

<style scoped></style>
