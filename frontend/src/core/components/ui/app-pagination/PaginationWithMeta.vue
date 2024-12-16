<script setup lang="ts">
import {
  Button,
} from '@/core/components/ui/button'

import {
  Pagination,
  PaginationEllipsis,
  PaginationFirst,
  PaginationLast,
  PaginationList,
  PaginationListItem,
  PaginationNext,
  PaginationPrev,
} from '@/core/components/ui/pagination'
import type { PaginationMeta } from '@/core/utils/pagination.utils'

import { computed } from 'vue'

// Props: Accepts PaginationMeta and sibling count
const props = defineProps<{
  meta: PaginationMeta,
  siblingCount?: number
}>()

// Emit: Emits an event when the page changes
const emit = defineEmits<{
  (e: 'update:page', page: number): void
}>()

// Compute the total pages and default page from meta
const totalPages = computed(() => props.meta.totalPages)
const currentPage = computed(() => props.meta.page)

// Handle page change
const handlePageChange = (newPage: number) => {
  emit('update:page', newPage)
}
</script>

<template>
	<Pagination
		v-slot="{ page }"
		:total="totalPages"
		:sibling-count="siblingCount || 1"
		show-edges
		:default-page="currentPage"
		@update:page="handlePageChange"
	>
		<PaginationList
			v-slot="{ items }"
			class="flex items-center gap-1"
		>
			<PaginationFirst />
			<PaginationPrev />

			<template
				v-for="(item, index) in items"
				:key="index"
			>
				<PaginationListItem
					v-if="item.type === 'page'"
					:value="item.value"
					as-child
				>
					<Button
						class="p-0 w-10 h-10"
						:variant="item.value === page ? 'default' : 'outline'"
						@click="handlePageChange(item.value)"
					>
						{{ item.value }}
					</Button>
				</PaginationListItem>
				<PaginationEllipsis
					v-else
					:key="item.type"
					:index="index"
				/>
			</template>

			<PaginationNext />
			<PaginationLast />
		</PaginationList>
	</Pagination>
</template>
