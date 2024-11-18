import { useCartStore } from '@/modules/kiosk/cart/stores/cart.store'

export const useResetKioskState = () => {
	const cartState = useCartStore()

	const resetAll = () => {
		cartState.clearCart()
	}

	return { resetAll }
}
