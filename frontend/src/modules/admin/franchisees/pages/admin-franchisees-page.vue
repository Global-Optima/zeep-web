<template>
	<AdminFranchiseesToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<AdminListLoader v-if="isPending" />

	<div v-else>
		<Card>
			<CardContent class="mt-4">
				<p
					v-if="!franchiseesResponse || franchiseesResponse.data.length === 0"
					class="text-muted-foreground"
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
	</div>
</template>

<script setup lang="ts">
import AdminListLoader from '@/core/components/admin-list-loader/AdminListLoader.vue'
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Card, CardContent, CardFooter } from '@/core/components/ui/card'
import { usePaginationFilter } from "@/core/hooks/use-pagination-filter.hook"
import AdminFranchiseesList from '@/modules/admin/franchisees/components/list/admin-franchisees-list.vue'
import AdminFranchiseesToolbar from '@/modules/admin/franchisees/components/list/admin-franchisees-toolbar.vue'
import type { FranchiseeFilterDTO } from '@/modules/admin/franchisees/models/franchisee.model'
import { franchiseeService } from '@/modules/admin/franchisees/services/franchisee.service'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter<FranchiseeFilterDTO>({})

// Fetch data using Vue Query
const { data: franchiseesResponse, isPending } = useQuery({
	queryKey: computed(() => ['admin-franchisees', filter.value]),
	queryFn: () => franchiseeService.getPaginated(filter.value),
})
</script>

<style scoped></style>
