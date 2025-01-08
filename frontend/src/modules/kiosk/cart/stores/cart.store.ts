import type { AdditiveCategoryItemDTO } from '@/modules/admin/additives/models/additives.model'
import type { StoreProductDetailsDTO } from '@/modules/admin/store-products/models/store-products.model'
import md5 from 'md5'
import { defineStore } from 'pinia'
import type { ProductSizeDTO } from '../../products/models/product.model'

export interface CartItem {
	key: string
	product: StoreProductDetailsDTO
	size: ProductSizeDTO
	additives: AdditiveCategoryItemDTO[]
	quantity: number
}

interface CartState {
	cartItems: { [key: string]: CartItem }
}

export const useCartStore = defineStore('ZEEP_CART', {
	state: (): CartState => ({
		cartItems: {},
	}),

	getters: {
		totalItems(state): number {
			return Object.values(state.cartItems).reduce((total, item) => total + item.quantity, 0)
		},
		isEmpty(state): boolean {
			return Object.values(state.cartItems).length === 0
		},
		totalPrice(state): number {
			return Object.values(state.cartItems).reduce((total, item) => {
				const additivesPrice = item.additives.reduce((sum, additive) => sum + additive.price, 0)
				return total + (item.size.basePrice + additivesPrice) * item.quantity
			}, 0)
		},
	},

	actions: {
		generateCartItemKey(
			product: StoreProductDetailsDTO,
			size: ProductSizeDTO,
			additives: AdditiveCategoryItemDTO[],
		): string {
			const additiveIds = additives
				.map(a => a.id)
				.sort()
				.join('-')
			return md5(`${product.id}-${size.id}-${additiveIds}`)
		},

		addToCart(
			product: StoreProductDetailsDTO,
			size: ProductSizeDTO,
			additives: AdditiveCategoryItemDTO[],
			quantity: number = 1,
		) {
			const key = this.generateCartItemKey(product, size, additives)
			if (this.cartItems[key]) {
				this.cartItems[key].quantity += quantity
			} else {
				this.cartItems[key] = {
					key,
					product,
					size,
					additives,
					quantity,
				}
			}
		},

		removeFromCart(
			product: StoreProductDetailsDTO,
			size: ProductSizeDTO,
			additives: AdditiveCategoryItemDTO[],
		) {
			const key = this.generateCartItemKey(product, size, additives)
			delete this.cartItems[key]
		},

		removeCartItemByKey(key: string) {
			delete this.cartItems[key]
		},

		incrementQuantity(key: string, amount: number = 1) {
			if (this.cartItems[key]) {
				this.cartItems[key].quantity += amount
			}
		},

		decrementQuantity(key: string, amount: number = 1) {
			if (this.cartItems[key]) {
				this.cartItems[key].quantity -= amount
				if (this.cartItems[key].quantity <= 0) {
					delete this.cartItems[key]
				}
			}
		},

		clearCart() {
			this.cartItems = {}
		},

		updateCartItem(
			key: string,
			updates: {
				size?: ProductSizeDTO
				additives?: AdditiveCategoryItemDTO[]
				quantity?: number
			},
		) {
			const existingItem = this.cartItems[key]
			if (!existingItem) return

			const updatedSize = updates.size || existingItem.size
			const updatedAdditives = updates.additives || existingItem.additives
			const updatedQuantity = updates.quantity ?? existingItem.quantity

			const newKey = this.generateCartItemKey(existingItem.product, updatedSize, updatedAdditives)

			if (newKey !== key) {
				delete this.cartItems[key]
			}

			this.cartItems[newKey] = {
				key: newKey,
				product: existingItem.product,
				size: updatedSize,
				additives: updatedAdditives,
				quantity: updatedQuantity,
			}
		},
	},

	persist: true,
})
