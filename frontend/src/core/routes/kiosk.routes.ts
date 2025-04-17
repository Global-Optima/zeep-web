import type { AppRouteRecord, ParentRoutePage } from './routes.types'

export const KIOSK_CHILDREN_ROUTES = {
	KIOSK_HOME: {
		path: '',
		meta: {
			title: 'Главная',
			requiresAuth: true,
		},
		component: () => import('@/modules/kiosk/products/pages/kiosk-home-page.vue'),
	},
	KIOSK_PRODUCT: {
		path: 'products/:id',
		meta: {
			title: 'Детали продукта',
			requiresAuth: true,
		},
		component: () => import('@/modules/kiosk/products/pages/kiosk-product-page.vue'),
	},
	KIOSK_LANDING: {
		path: 'landing',
		meta: {
			title: 'Популярное',
			requiresAuth: true,
		},
		component: () => import('@/modules/kiosk/landing/pages/kiosk-landing-page.vue'),
	},
	KIOSK_CART_PAYMENT: {
		path: 'cart/payment/:orderId',
		meta: {
			title: 'Оплата',
			requiresAuth: true,
		},
		component: () => import('@/modules/kiosk/cart/pages/kiosk-cart-payment-page.vue'),
	},
	KIOSK_CART_PAYMENT_SUCCESS: {
		path: 'cart/payment/:orderId/success',
		meta: {
			title: 'Успешная оплата',
			requiresAuth: true,
		},
		component: () => import('@/modules/kiosk/cart/pages/kiosk-cart-payment-success-page.vue'),
	},
	KIOSK_CHECKLIST: {
		path: 'checklist',
		meta: {
			title: 'Проверка',
			requiresAuth: true,
		},
		component: () => import('@/modules/kiosk/checklist/pages/kiosk-checklist-page.vue'),
	},
} satisfies AppRouteRecord

export const KIOSK_ROUTES_CONFIG: ParentRoutePage = {
	path: '/kiosk',
	component: () => import('@/core/layouts/kiosk/app-kiosk-layout.vue'),
	children: KIOSK_CHILDREN_ROUTES,
}
