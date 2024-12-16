<template>
	<AdminProductCategoriesToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<Card>
		<CardContent class="mt-4">
			<p
				v-if="!productCategories || productCategories.length === 0"
				class="text-muted-foreground"
			>
				Категории не найдены
			</p>
			<AdminProductCategoriesList
				v-else
				:productCategories="productCategories"
			/>
		</CardContent>
	</Card>
</template>

<script setup lang="ts">
import { Card, CardContent } from '@/core/components/ui/card'
import AdminProductCategoriesList from '@/modules/admin/product-categories/components/list/admin-product-categories-list.vue'
import AdminProductCategoriesToolbar from '@/modules/admin/product-categories/components/list/admin-product-categories-toolbar.vue'
import type { ProductCategoriesFilter } from '@/modules/admin/product-categories/models/product-categories.model'
import { productCategoriesService } from '@/modules/admin/product-categories/services/product-categories.service'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const filter = ref<ProductCategoriesFilter>({
  search: '',
})

const { data: productCategories } = useQuery({
  queryKey: computed(() => ['product-categories', filter.value]),
  queryFn: () => productCategoriesService.getProductCategories(filter.value),
  initialData: []
})

function updateFilter(updatedFilter: ProductCategoriesFilter) {
  filter.value = {...filter.value, ...updatedFilter}
}
</script>

<style scoped></style>
