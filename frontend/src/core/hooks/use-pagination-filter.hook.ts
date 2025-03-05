import { shallowRef } from 'vue'
import { DEFAULT_PAGINATION_META, type PaginationParams } from '../utils/pagination.utils'

export function usePaginationFilter<T extends PaginationParams>(defaultFilter: T) {
	// Reactive filter state (shallow ref for performance optimization)
	const filter = shallowRef<T>({ ...defaultFilter })

	// Function to update filter reactively
	function updateFilter(updatedFilter: Partial<T>) {
		filter.value = { ...filter.value, ...updatedFilter }
	}

	// Function to update pagination page
	function updatePage(page: number) {
		filter.value = { ...filter.value, page }
	}

	// Function to update pagination page size
	function updatePageSize(pageSize: number) {
		filter.value = {
			...filter.value,
			pageSize,
			page: DEFAULT_PAGINATION_META.page,
		}
	}

	return {
		filter,
		updateFilter,
		updatePage,
		updatePageSize,
	}
}
