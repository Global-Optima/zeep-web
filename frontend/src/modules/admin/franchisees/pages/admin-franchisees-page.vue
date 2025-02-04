<template>
	<AdminFranchiseesToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

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
</template>

<script setup lang="ts">
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Card, CardContent, CardFooter } from '@/core/components/ui/card'
import { DEFAULT_PAGINATION_META } from '@/core/utils/pagination.utils'
import AdminFranchiseesList from '@/modules/admin/franchisees/components/list/admin-franchisees-list.vue'
import AdminFranchiseesToolbar from '@/modules/admin/franchisees/components/list/admin-franchisees-toolbar.vue'
import type { FranchiseeFilterDTO } from '@/modules/admin/franchisees/models/franchisee.model'
import { franchiseeService } from '@/modules/admin/franchisees/services/franchisee.service'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const filter = ref<FranchiseeFilterDTO>({})

const { data: franchiseesResponse } = useQuery({
  queryKey: computed(() => ['admin-franchisees', filter.value]),
  queryFn: () => franchiseeService.getPaginated(filter.value),
})

function updateFilter(updatedFilter: FranchiseeFilterDTO) {
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
