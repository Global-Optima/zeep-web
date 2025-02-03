<template>
	<AdminStoresToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<Card>
		<CardContent class="mt-4">
			<p
				v-if="!stores || stores.data.length === 0"
				class="text-muted-foreground"
			>
				Магазины не найдены
			</p>

			<AdminStoresList
				v-else
				:stores="stores.data"
			/>
		</CardContent>
	</Card>
</template>

<script setup lang="ts">
import { Card, CardContent } from '@/core/components/ui/card'
import AdminStoresList from '@/modules/admin/stores/components/list/admin-stores-list.vue'
import AdminStoresToolbar from '@/modules/admin/stores/components/list/admin-stores-toolbar.vue'
import type { StoresFilter } from '@/modules/admin/stores/models/stores-dto.model'
import { storesService } from '@/modules/admin/stores/services/stores.service'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

// Reactive filter object
const filter = ref<StoresFilter>({})

// Query stores data
const { data: stores } = useQuery({
  queryKey: computed(() => ['stores', filter.value]),
  queryFn: () => storesService.getPaginated(filter.value),
})

// Update filter handler
function updateFilter(updatedFilter: Partial<StoresFilter>) {
  filter.value = {...filter.value, ...updatedFilter}
}
</script>

<style scoped></style>
