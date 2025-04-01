<template>
	<section
		class="top-0 z-10 sticky flex items-center gap-2 px-8 py-6 w-full overflow-x-auto transition-all no-scrollbar"
		data-testid="search-bar-wrapper"
	>
		<KioskHomeToolbarSearch
			:searchTerm="searchTerm"
			@update:searchTerm="emit('update:searchTerm', $event)"
		/>

		<KioskHomeToolbarCategories
			:categories="categories"
			:selectedCategoryId="selectedCategoryId"
			@update:category="emit('update:category', $event)"
		/>
	</section>
</template>

<script setup lang="ts">
import KioskHomeToolbarCategories from '@/modules/kiosk/products/components/home/toolbar/kiosk-home-toolbar-categories.vue'
import KioskHomeToolbarSearch from '@/modules/kiosk/products/components/home/toolbar/kiosk-home-toolbar-search.vue'

defineProps<{
  selectedCategoryId: number | null
  searchTerm: string
  categories: Array<{ id: number; name: string }>
}>()

const emit = defineEmits<{
  (e: 'update:category', categoryId: number): void
  (e: 'update:searchTerm', term: string): void
}>()
</script>

<style scoped>
/* Hide scrollbar for horizontal overflow */
.no-scrollbar::-webkit-scrollbar {
  display: none;
}
.no-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
