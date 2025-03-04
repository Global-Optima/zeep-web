import { useDebounce } from '@vueuse/core'
import { computed, ref } from 'vue'

export function useDebouncedSearch(initialValue = '', debounceDelay = 500) {
	const searchTerm = ref(initialValue)

	// Create a debounced version of the searchTerm
	const debouncedSearchTerm = useDebounce(
		computed(() => searchTerm.value),
		debounceDelay,
	)

	// Function to update searchTerm
	function updateSearch(value: string) {
		searchTerm.value = value
	}

	return {
		searchTerm, // Directly usable for v-model
		debouncedSearchTerm, // The debounced version for API requests
		updateSearch, // Method to manually update search term
	}
}
