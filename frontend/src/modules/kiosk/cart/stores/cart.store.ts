// useCartStore.ts
import md5 from 'md5'
import { defineStore } from 'pinia'
import type {
	Additive,
	ProductSize,
	StoreProductDetails,
} from '../../products/models/product.model'

export interface CartItem {
	key: string
	product: StoreProductDetails
	size: ProductSize
	additives: Additive[]
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
			product: StoreProductDetails,
			size: ProductSize,
			additives: Additive[],
		): string {
			const additiveIds = additives
				.map(a => a.id)
				.sort()
				.join('-')
			return md5(`${product.id}-${size.id}-${additiveIds}`)
		},

		addToCart(
			product: StoreProductDetails,
			size: ProductSize,
			additives: Additive[],
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

		removeFromCart(product: StoreProductDetails, size: ProductSize, additives: Additive[]) {
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
	},

	persist: true,
})
