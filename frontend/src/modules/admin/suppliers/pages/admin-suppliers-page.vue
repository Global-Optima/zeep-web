<template>
	<AdminListLoader v-if="isPending" />

	<div v-else>
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
	</div>
</template>

<script setup lang="ts">
import AdminListLoader from '@/core/components/admin-list-loader/AdminListLoader.vue'
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Card, CardContent, CardFooter } from '@/core/components/ui/card'
import { usePaginationFilter } from '@/core/hooks/use-pagination-filter.hook'
import AdminSuppliersList from '@/modules/admin/suppliers/components/list/admin-suppliers-list.vue'
import AdminSuppliersToolbar from '@/modules/admin/suppliers/components/list/admin-suppliers-toolbar.vue'
import type { SuppliersFilterDTO } from '@/modules/admin/suppliers/models/suppliers.model'
import { suppliersService } from '@/modules/admin/suppliers/services/suppliers.service'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter<SuppliersFilterDTO>({});

const { data: suppliersResponse, isPending } = useQuery({
  queryKey: computed(() => ['admin-suppliers', filter.value]),
  queryFn: () => suppliersService.getSuppliers(filter.value),
})
</script>

<style scoped></style>
