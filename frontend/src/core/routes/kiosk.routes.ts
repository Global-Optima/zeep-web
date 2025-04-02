import AppKioskLayout from '../layouts/kiosk/app-kiosk-layout.vue'
import type { AppRouteRecord, ParentRoutePage } from './routes.types'

export const KIOSK_CHILDREN_ROUTES = {
	KIOSK_HOME: {
		path: '',
		meta: {
			title: 'Главная',
		},
		component: () => import('@/modules/kiosk/products/pages/kiosk-home-page.vue'),
	},
	KIOSK_PRODUCT: {
		path: 'products/:id',
		meta: {
			title: 'Детали продукта',
		},
		component: () => import('@/modules/kiosk/products/pages/kiosk-product-page.vue'),
	},
	KIOSK_LANDING: {
		path: 'landing',
		meta: {
			title: 'Популярное',
		},
		component: () => import('@/modules/kiosk/landing/pages/kiosk-landing-page.vue'),
	},
	KIOSK_CART_PAYMENT: {
		path: 'cart/payment/:orderId',
		meta: {
			title: 'Оплата',
		},
		component: () => import('@/modules/kiosk/cart/pages/kiosk-cart-payment-page.vue'),
	},
	KIOSK_CART_PAYMENT_SUCCESS: {
		path: 'cart/payment/:orderId/success',
		meta: {
			title: 'Успешная оплата',
		},
		component: () => import('@/modules/kiosk/cart/pages/kiosk-cart-payment-success-page.vue'),
	},
	KIOSK_CHECKLIST: {
		path: 'checklist',
		meta: {
			title: 'Проверка',
		},
		component: () => import('@/modules/kiosk/checklist/pages/kiosk-checklist-page.vue'),
	},
} satisfies AppRouteRecord

export const KIOSK_ROUTES_CONFIG: ParentRoutePage = {
	path: '/kiosk',
	component: AppKioskLayout,
	children: KIOSK_CHILDREN_ROUTES,
}
