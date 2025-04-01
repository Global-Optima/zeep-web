<template>
	<div
		class="flex bg-white bg-clip-padding bg-opacity-70 backdrop-filter backdrop-blur-sm p-1 border border-slate-100 rounded-full min-w-min overflow-x-auto no-scrollbar"
		data-testid="category-buttons"
		@touchstart.stop
		@touchmove.stop
	>
		<button
			v-for="category in categories"
			:key="category.id"
			@click="selectCategory(category.id)"
			:aria-pressed="selectedCategoryId === category.id"
			:ref="el => setButtonRef(category.id, el)"
			:class="[
          'px-8 py-6 text-3xl rounded-full whitespace-nowrap',
          selectedCategoryId === category.id
            ? 'bg-primary text-primary-foreground font-medium'
            : 'bg-transparent text-gray-600 hover:bg-gray-100',
        ]"
			data-testid="category-button"
		>
			{{ category.name }}
		</button>
	</div>
</template>

<script setup lang="ts">
import { nextTick, ref, watch, type ComponentPublicInstance } from 'vue'

const emit = defineEmits<{ (e: 'update:category', categoryId: number): void }>()
const props = defineProps<{
  selectedCategoryId: number | null
  categories: Array<{ id: number; name: string }>
}>()

// Store refs for each category button
const categoryButtons = ref<Record<number, HTMLElement>>({})

function selectCategory(categoryId: number) {
  emit('update:category', categoryId)
}

function setButtonRef(categoryId: number, el: Element | ComponentPublicInstance | null) {
  let element: HTMLElement | null = null

  if (el) {
    if (el instanceof HTMLElement) {
      element = el
    } else if ('$el' in el && el.$el instanceof HTMLElement) {
      element = el.$el
    }
  }

  if (element) {
    categoryButtons.value[categoryId] = element
  } else {
    delete categoryButtons.value[categoryId]
  }
}

// Watch for changes in the active category and scroll it into view.
watch(
  () => props.selectedCategoryId,
  async (newVal) => {
    if (newVal !== null) {
      // Wait for DOM update
      await nextTick()
      const btn = categoryButtons.value[newVal]
      if (btn) {
        // You can adjust "inline" to "center" or "start" for left alignment.
        btn.scrollIntoView({ behavior: 'smooth', inline: 'center' })
      }
    }
  }
)
</script>
