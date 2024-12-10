<template>
	<AdminStoresToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>
	<AdminStoresList :stores="stores" />
</template>

<script setup lang="ts">
import AdminStoresList from '@/modules/admin/stores/components/list/admin-stores-list.vue'
import AdminStoresToolbar from '@/modules/admin/stores/components/list/admin-stores-toolbar.vue'
import type { StoresFilter } from '@/modules/stores/models/stores-dto.model'
import { storesService } from '@/modules/stores/services/stores.service'
import { useQuery } from '@tanstack/vue-query'
import { reactive } from 'vue'

// Reactive filter object
const filter = reactive<StoresFilter>({
  searchTerm: '',
  isFranchise: undefined,
})

// Query stores data
const { data: stores } = useQuery({
  queryKey: ['stores', filter],
  queryFn: () => storesService.getStores(filter),
  initialData: [],
})

// Update filter handler
function updateFilter(updatedFilter: Partial<StoresFilter>) {
  Object.assign(filter, updatedFilter)
}
</script>

<style scoped></style>
