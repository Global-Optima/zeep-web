import type { AppRouteRecord, ParentRoutePage } from '../config/routes.config'
import AppKioskLayout from '../layouts/kiosk/app-kiosk-layout.vue'

export const KIOSK_CHILDREN_ROUTES = {
	KIOSK_HOME: {
		path: '',
		meta: {
			title: 'Главная',
		},
		component: () => import('@/modules/kiosk/pages/kiosk-home-page.vue'),
	},
	KIOSK_PRODUCT_DETAILS: {
		path: ':id',
		meta: {
			title: 'Детали продукта',
		},
		component: () => import('@/modules/kiosk/pages/kiosk-product-details-page.vue'),
	},
	KIOSK_CART: {
		path: 'cart',
		meta: {
			title: 'Корзина',
		},
		component: () => import('@/modules/kiosk/pages/kiosk-cart-page.vue'),
	},
} satisfies AppRouteRecord

export const KIOSK_ROUTES_CONFIG: ParentRoutePage = {
	path: 'kiosk',
	component: AppKioskLayout,
	children: KIOSK_CHILDREN_ROUTES,
}
