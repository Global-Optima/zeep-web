import { parseLocalStorageItem } from '@/core/utils/local-storage.utils'
import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export const CURRENT_STORE_STORAGE_KEY = 'ZEEP_CURRENT_STORE_ID'

export const useCurrentStoreStore = defineStore('CURRENT_STORE', () => {
	const initialStore = parseLocalStorageItem<number>(CURRENT_STORE_STORAGE_KEY) || null
	const currentStoreId = ref<number | null>(initialStore)

	function setCurrentStore(store: number): void {
		currentStoreId.value = store
	}

	if (typeof window !== 'undefined') {
		watch(
			currentStoreId,
			newValue => {
				try {
					window.localStorage.setItem(CURRENT_STORE_STORAGE_KEY, JSON.stringify(newValue))
				} catch (error) {
					console.error('[useCurrentStoreStore] Failed to store current store data:', error)
				}
			},
			{ deep: true },
		)
	}

	return {
		currentStoreId,
		setCurrentStore,
	}
})
