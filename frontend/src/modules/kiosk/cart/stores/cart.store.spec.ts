// src/core/stores/cart.spec.ts
import md5 from 'md5'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it } from 'vitest'
import { useCartStore } from './cart.store'

// Mock minimal DTOs
const mockProduct = {
	id: 1,
	name: 'Prod',
	description: '',
	imageUrl: '',
	videoUrl: '',
	category: { id: 1, name: '', description: '', machineCategory: 'TEA' },
	productId: 1,
	storePrice: 100,
	basePrice: 0,
	storeProductSizeCount: 0,
	productSizeCount: 0,
	isAvailable: true,
	isOutOfStock: false,
	sizes: [],
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
} as any

const mockSize = {
	id: 2,
	name: 'M',
	storePrice: 50,
	unit: { id: 1, name: 'u', conversionFactor: 1 },
	additives: [],
	ingredients: [],
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
} as any

const defaultAdditive = {
	id: 10,
	additiveId: 10,
	storePrice: 5,
	isDefault: false,
	isOutOfStock: false,
	name: '',
	description: '',
	basePrice: 0,
	imageUrl: '',
	size: 0,
	unit: { id: 1, name: '', conversionFactor: 1 },
	category: {
		id: 1,
		name: '',
		description: '',
		isMultipleSelect: false,
		isRequired: false,
		machineCategory: 0,
	},
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
} as any

describe('Cart Store', () => {
	let store: ReturnType<typeof useCartStore>

	beforeEach(() => {
		setActivePinia(createPinia())
		store = useCartStore()
		store.clearCart()
		store.closeModal()
	})

	it('initial state is empty', () => {
		expect(store.isEmpty).toBe(true)
		expect(store.totalItems).toBe(0)
		expect(store.totalPrice).toBe(0)
		expect(store.isModalOpen).toBe(false)
	})

	it('toggles and closes modal', () => {
		store.toggleModal()
		expect(store.isModalOpen).toBe(true)
		store.closeModal()
		expect(store.isModalOpen).toBe(false)
	})

	it('generates stable key from product, size, additives', () => {
		const key = store.generateCartItemKey(mockProduct, mockSize, [defaultAdditive])
		const expected = md5(`1-2-${defaultAdditive.additiveId}`)
		expect(key).toBe(expected)
	})

	it('adds items to cart and updates quantities', () => {
		store.addToCart(mockProduct, mockSize, [defaultAdditive], 1)
		expect(store.totalItems).toBe(1)
		store.addToCart(mockProduct, mockSize, [defaultAdditive], 2)
		expect(store.totalItems).toBe(3)
	})

	it('removes item by product and by key', () => {
		store.addToCart(mockProduct, mockSize, [], 1)
		const key = store.generateCartItemKey(mockProduct, mockSize, [])
		store.removeFromCart(mockProduct, mockSize, [])
		expect(store.isEmpty).toBe(true)
		// add again and test removeCartItemByKey
		store.addToCart(mockProduct, mockSize, [], 1)
		store.removeCartItemByKey(key)
		expect(store.isEmpty).toBe(true)
	})

	it('increments and decrements quantity, removes when zero', () => {
		store.addToCart(mockProduct, mockSize, [], 1)
		const key = store.generateCartItemKey(mockProduct, mockSize, [])
		store.incrementQuantity(key, 2)
		expect(store.cartItems[key].quantity).toBe(3)
		store.decrementQuantity(key, 1)
		expect(store.cartItems[key].quantity).toBe(2)
		store.decrementQuantity(key, 2)
		expect(store.cartItems[key]).toBeUndefined()
	})

	it('clears the cart', () => {
		store.addToCart(mockProduct, mockSize, [], 1)
		store.clearCart()
		expect(store.isEmpty).toBe(true)
	})

	it('updates cart items correctly changing size/additives', () => {
		store.addToCart(mockProduct, mockSize, [], 1)
		const originalKey = store.generateCartItemKey(mockProduct, mockSize, [])
		const newSize = { ...mockSize, id: 3, storePrice: 30 }
		store.updateCartItem(originalKey, { size: newSize })
		const newKey = store.generateCartItemKey(mockProduct, newSize, [])
		expect(store.cartItems[newKey].quantity).toBe(1)
	})

	it('merges quantity when update leads to existing key', () => {
		store.addToCart(mockProduct, mockSize, [], 1)
		const altSize = { ...mockSize, id: 4 }
		store.addToCart(mockProduct, altSize, [], 2)
		const key1 = store.generateCartItemKey(mockProduct, mockSize, [])
		const key2 = store.generateCartItemKey(mockProduct, altSize, [])
		store.updateCartItem(key1, { size: altSize, quantity: 3 })
		expect(store.cartItems[key2].quantity).toBe(5)
	})

	it('getter totalPrice calculates correctly including additives', () => {
		const additive1 = { ...defaultAdditive, storePrice: 10, isDefault: false }
		const additive2 = { ...defaultAdditive, storePrice: 5, isDefault: true }
		store.addToCart(mockProduct, mockSize, [additive1, additive2], 2)
		// size=50, only additive1 counts => (50+10)*2 = 120
		expect(store.totalPrice).toBe(120)
	})
})
