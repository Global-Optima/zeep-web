<template>
	<AdminAdditiveCategoriesToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<Card>
		<CardContent class="mt-4">
			<p
				v-if="!additiveCategories || additiveCategories.length === 0"
				class="text-muted-foreground"
			>
				Категории не найдены
			</p>
			<AdminAdditiveCategoriesList
				v-else
				:additiveCategories="additiveCategories"
			/>
		</CardContent>
	</Card>
</template>

<script setup lang="ts">
import { Card, CardContent } from '@/core/components/ui/card'
import AdminAdditiveCategoriesList from '@/modules/admin/additive-categories/components/list/admin-additive-categories-list.vue'
import AdminAdditiveCategoriesToolbar from '@/modules/admin/additive-categories/components/list/admin-additive-categories-toolbar.vue'
import type { AdditiveCategoriesFilterQuery } from '@/modules/admin/additives/models/additives.model'
import { additivesService } from '@/modules/admin/additives/services/additives.service'
import type { ProductCategoriesFilter } from '@/modules/admin/product-categories/models/product-categories.model'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const filter = ref<AdditiveCategoriesFilterQuery>({
})

const { data: additiveCategories } = useQuery({
  queryKey: computed(() => ['additive-categories', filter.value]),
  queryFn: () => additivesService.getAdditiveCategories(filter.value),
  initialData: []
})

function updateFilter(updatedFilter: ProductCategoriesFilter) {
  filter.value = {...filter.value, ...updatedFilter}
}
</script>

<style scoped></style>
