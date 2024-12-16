<template>
	<AdminStoresToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<Card>
		<CardContent class="mt-4">
			<AdminStoresList :stores="stores ?? []" />
		</CardContent>
	</Card>
</template>

<script setup lang="ts">
import { Card, CardContent } from '@/core/components/ui/card'
import AdminStoresList from '@/modules/admin/stores/components/list/admin-stores-list.vue'
import AdminStoresToolbar from '@/modules/admin/stores/components/list/admin-stores-toolbar.vue'
import type { StoresFilter } from '@/modules/stores/models/stores-dto.model'
import { storesService } from '@/modules/stores/services/stores.service'
import { useQuery } from '@tanstack/vue-query'
import { ref } from 'vue'

// Reactive filter object
const filter = ref<StoresFilter>({
  searchTerm: '',
  isFranchise: undefined,
})

// Query stores data
const { data: stores } = useQuery({
  queryKey: ['stores', filter],
  queryFn: () => storesService.getStores(filter.value),
})

// Update filter handler
function updateFilter(updatedFilter: Partial<StoresFilter>) {
  Object.assign(filter.value, updatedFilter)
}
</script>

<style scoped></style>
