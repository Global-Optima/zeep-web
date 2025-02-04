<template>
	<div class="flex-1 gap-4 grid auto-rows-max mx-auto max-w-6xl">
		<!-- Header -->
		<div class="flex items-center gap-4">
			<Button
				variant="outline"
				size="icon"
				@click="onCancel"
			>
				<ChevronLeft class="w-5 h-5" />
				<span class="sr-only">Назад</span>
			</Button>
			<h1 class="flex-1 sm:grow-0 font-semibold text-xl tracking-tight whitespace-nowrap shrink-0">
				Магазины франчайзи {{ franchisee.name }}
			</h1>
		</div>

		<!-- Stores List -->
		<Card>
			<CardContent class="mt-4">
				<p
					v-if="!storesResponse || storesResponse.data.length === 0"
					class="text-muted-foreground"
				>
					Магазины франчайзи не найдены
				</p>
				<AdminStoresList
					v-else
					:stores="storesResponse.data"
				/>
			</CardContent>
			<CardFooter class="flex justify-end">
				<PaginationWithMeta
					v-if="storesResponse"
					:meta="storesResponse.pagination"
					@update:page="updatePage"
					@update:pageSize="updatePageSize"
				/>
			</CardFooter>
		</Card>
	</div>
</template>

<script setup lang="ts">
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardFooter } from '@/core/components/ui/card'
import { DEFAULT_PAGINATION_META } from '@/core/utils/pagination.utils'
import type { FranchiseeDTO } from '@/modules/admin/franchisees/models/franchisee.model'
import AdminStoresList from '@/modules/admin/stores/components/list/admin-stores-list.vue'
import type { StoresFilter } from '@/modules/admin/stores/models/stores-dto.model'
import { storesService } from '@/modules/admin/stores/services/stores.service'
import { useQuery } from '@tanstack/vue-query'
import { ChevronLeft } from 'lucide-vue-next'
import { computed, ref } from 'vue'

// Props
const { franchisee } = defineProps<{ franchisee: FranchiseeDTO }>()

const emits = defineEmits<{
  (e: 'onCancel'): void
}>()

const onCancel = () => {
  emits('onCancel')
}

// Filtering State
const filter = ref<StoresFilter>({})

// Fetch Stores by Franchisee
const { data: storesResponse } = useQuery({
  queryKey: computed(() => ['admin-franchisee-stores', franchisee.id]),
  queryFn: () => storesService.getPaginated({ ...filter.value, franchiseeId: franchisee.id }),
})

function updateFilter(updatedFilter: StoresFilter) {
  filter.value = { ...filter.value, ...updatedFilter }
}

function updatePage(page: number) {
  updateFilter({ pageSize: DEFAULT_PAGINATION_META.pageSize, page: page })
}

function updatePageSize(pageSize: number) {
  updateFilter({ pageSize: pageSize, page: DEFAULT_PAGINATION_META.page })
}
</script>

<style scoped></style>
