import { defineStore } from 'pinia'

interface ModalState {
	isOpen: boolean
	productId: number | null
}

export const useCurrentProductStore = defineStore('CURRENT_PRODUCT_STORE', {
	state: (): ModalState => ({
		isOpen: false,
		productId: null,
	}),
	actions: {
		openModal(productId: number) {
			this.isOpen = true
			this.productId = productId
		},
		closeModal() {
			this.isOpen = false
			this.productId = null
		},
	},
	getters: {
		isModalOpen(state): boolean {
			return state.isOpen
		},
		currentProductId(state): number | null {
			return state.productId
		},
	},
})
