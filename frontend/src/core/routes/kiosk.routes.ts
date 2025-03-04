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
			title: 'Детали товара',
		},
		component: () => import('@/modules/kiosk/products/pages/kiosk-product-page.vue'),
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
