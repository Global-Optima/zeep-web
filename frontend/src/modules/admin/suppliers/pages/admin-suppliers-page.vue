<template>
	<AdminSuppliersToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<Card>
		<CardContent class="mt-4">
			<p
				v-if="!suppliersResponse || suppliersResponse.data.length === 0"
				class="text-muted-foreground"
			>
				Поставщики не найдены
			</p>
			<AdminSuppliersList
				v-else
				:suppliers="suppliersResponse.data"
			/>
		</CardContent>
		<CardFooter class="flex justify-end">
			<PaginationWithMeta
				v-if="suppliersResponse"
				:meta="suppliersResponse.pagination"
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
import AdminSuppliersList from '@/modules/admin/suppliers/components/list/admin-suppliers-list.vue'
import AdminSuppliersToolbar from '@/modules/admin/suppliers/components/list/admin-suppliers-toolbar.vue'
import type { SuppliersFilterDTO } from '@/modules/admin/suppliers/models/suppliers.model'
import { suppliersService } from '@/modules/admin/suppliers/services/suppliers.service'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const filter = ref<SuppliersFilterDTO>({})

const { data: suppliersResponse } = useQuery({
  queryKey: computed(() => ['admin-suppliers', filter.value]),
  queryFn: () => suppliersService.getSuppliers(filter.value),
})

function updateFilter(updatedFilter: SuppliersFilterDTO) {
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
