<template>
	<div
		class="relative bg-white bg-clip-padding bg-opacity-70 backdrop-filter backdrop-blur-sm border border-slate-100 rounded-full"
		data-testid="search-input-container"
	>
		<!-- Collapsed State (Icon Only) -->
		<div
			v-if="!isInputFocused && !searchTerm"
			@click="focusInput"
			class="flex justify-center items-center rounded-full size-[100px] hover:scale-105 transition-transform cursor-pointer select-none"
			data-testid="collapsed-search-trigger"
			aria-label="Раскрыть поиск"
		>
			<Search class="size-10 text-slate-500" />
		</div>

		<!-- Expanded Search Input -->
		<input
			v-if="isInputFocused || searchTerm"
			ref="searchInputRef"
			type="text"
			:value="searchTerm"
			@input="updateSearchTerm"
			@focus="handleFocus"
			@blur="handleBlur"
			class="bg-transparent px-10 rounded-full focus:outline-none w-56 sm:w-80 h-[100px] text-3xl transition-[width] duration-300"
			placeholder="Поиск"
			aria-label="Поиск продуктов"
			data-testid="search-input"
		/>

		<!-- Clear Button -->
		<button
			v-if="searchTerm"
			@click="clearSearchTerm"
			class="top-5 right-8 absolute p-2 focus:outline-none text-primary"
			aria-label="Очистить поле поиска"
			data-testid="clear-button"
		>
			<X class="size-10" />
		</button>
	</div>
</template>

<script setup lang="ts">
import { Search, X } from 'lucide-vue-next'
import { ref } from 'vue'

const emit = defineEmits<{ (e: 'update:searchTerm', term: string): void }>()
defineProps<{ searchTerm: string }>()

const isInputFocused = ref(false)
const searchInputRef = ref<HTMLInputElement | null>(null)

/**
 * When the search icon is clicked, focus the input.
 */
function focusInput() {
  isInputFocused.value = true
  requestAnimationFrame(() => searchInputRef.value?.focus())
}

function handleFocus() {
  isInputFocused.value = true
}

function handleBlur() {
  if (!searchInputRef.value?.value.trim()) {
    isInputFocused.value = false
  }
}

function updateSearchTerm(event: Event) {
  emit('update:searchTerm', (event.target as HTMLInputElement).value)
}

function clearSearchTerm() {
  emit('update:searchTerm', '')
  isInputFocused.value = false
}
</script>

<style scoped>
/* Smooth scaling animation */
.fade-scale-enter-active,
.fade-scale-leave-active {
  transition: transform 0.2s ease, opacity 0.2s ease;
}
.fade-scale-enter-from,
.fade-scale-leave-to {
  opacity: 0;
  transform: scale(0.95);
}
</style>
