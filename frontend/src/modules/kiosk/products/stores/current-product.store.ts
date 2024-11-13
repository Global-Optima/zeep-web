import { defineStore } from 'pinia'

interface CurrentProductState {
	selectedProductId: number | null
	isBottomSheetOpen: boolean
}

export const useProductStore = defineStore('CURRENT_PRODUCT_STORE', {
	state: (): CurrentProductState => ({
		selectedProductId: null,
		isBottomSheetOpen: false,
	}),

	actions: {
		selectProduct(productId: number) {
			this.selectedProductId = productId
			this.isBottomSheetOpen = true
		},
		closeBottomSheet() {
			this.isBottomSheetOpen = false
			this.selectedProductId = null
		},
	},

	getters: {
		isProductSelected: (state): boolean => state.selectedProductId !== null,
	},
})
