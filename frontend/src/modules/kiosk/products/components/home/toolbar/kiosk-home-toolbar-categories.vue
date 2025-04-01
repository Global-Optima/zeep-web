<template>
	<div
		class="flex bg-white bg-clip-padding bg-opacity-70 backdrop-filter backdrop-blur-sm p-1 border border-slate-100 rounded-full min-w-min overflow-x-auto no-scrollbar"
		data-testid="category-buttons"
	>
		<button
			v-for="category in categories"
			:key="category.id"
			@click="selectCategory(category.id)"
			:aria-pressed="selectedCategoryId === category.id"
			:class="[
        'px-7 py-5 text-base sm:text-2xl rounded-full font-medium whitespace-nowrap',
        selectedCategoryId === category.id
          ? 'bg-primary text-primary-foreground'
          : 'bg-transparent text-gray-700 hover:bg-gray-100',
      ]"
			data-testid="category-button"
		>
			{{ category.name }}
		</button>
	</div>
</template>

<script setup lang="ts">
const emit = defineEmits<{ (e: 'update:category', categoryId: number): void }>()
defineProps<{
  selectedCategoryId: number | null
  categories: Array<{ id: number; name: string }>
}>()

function selectCategory(categoryId: number) {
  emit('update:category', categoryId)
}
</script>
