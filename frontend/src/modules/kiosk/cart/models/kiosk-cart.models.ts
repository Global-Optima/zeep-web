import type { LucideIcon } from 'lucide-vue-next'

// checkout.types.ts
export enum CheckoutStep {
	CUSTOMER = 'customer',
	PAYMENT = 'payment',
	CONFIRMATION = 'confirmation',
}

export enum PaymentMethod {
	KASPI = 'kaspi',
	CARD = 'card',
}

export interface PaymentOption {
	id: PaymentMethod
	label: string
	icon: LucideIcon
}
