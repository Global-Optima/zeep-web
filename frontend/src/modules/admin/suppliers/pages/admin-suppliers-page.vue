<template>
	<AdminSuppliersToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<Card>
		<CardContent class="mt-4">
			<AdminSuppliersList :suppliers="suppliers" />
		</CardContent>
	</Card>
</template>

<script setup lang="ts">
import { Card, CardContent } from '@/core/components/ui/card'
import AdminSuppliersList from '@/modules/admin/suppliers/components/list/admin-suppliers-list.vue'
import AdminSuppliersToolbar from '@/modules/admin/suppliers/components/list/admin-suppliers-toolbar.vue'
import type { SuppliersFilter } from '@/modules/admin/suppliers/models/suppliers.model'
import { suppliersService } from '@/modules/admin/suppliers/services/suppliers.service'
import type { StoresFilter } from '@/modules/stores/models/stores-dto.model'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const filter = ref<SuppliersFilter>({
  search: '',
})

const { data: suppliers } = useQuery({
  queryKey: computed(() => ['suppliers', filter.value]),
  queryFn: () => suppliersService.getSuppliers(filter.value),
  initialData: []
})

function updateFilter(updatedFilter: Partial<StoresFilter>) {
  filter.value = {...filter.value, ...updatedFilter}
}
</script>

<style scoped></style>
