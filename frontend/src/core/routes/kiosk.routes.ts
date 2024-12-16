import type { AppRouteRecord, ParentRoutePage } from '../config/routes.config'
import AppKioskLayout from '../layouts/kiosk/app-kiosk-layout.vue'

export const KIOSK_CHILDREN_ROUTES = {
	KIOSK_HOME: {
		path: '',
		meta: {
			title: 'Главная',
		},
		component: () => import('@/modules/kiosk/products/pages/kiosk-home-page.vue'),
	},
	KIOSK_ORDERS: {
		path: 'orders',
		meta: {
			title: 'Заказы',
		},
		component: () => import('@/modules/kiosk/orders/pages/kiosk-orders-page-v2.vue'),
	},
	KIOSK_CART: {
		path: 'cart',
		meta: {
			title: 'Корзина',
		},
		component: () => import('@/modules/kiosk/cart/pages/kiosk-cart-page.vue'),
	},
	KIOSK_LANDING: {
		path: 'landing',
		meta: {
			title: 'Популряное',
		},
		component: () => import('@/modules/kiosk/landing/pages/kiosk-landing-page.vue'),
	},
} satisfies AppRouteRecord

export const KIOSK_ROUTES_CONFIG: ParentRoutePage = {
	path: '/kiosk',
	component: AppKioskLayout,
	children: KIOSK_CHILDREN_ROUTES,
}
