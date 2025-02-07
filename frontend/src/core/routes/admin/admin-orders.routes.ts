import type { AppRouteRecord } from '../routes.types'

export const ADMIN_STORE_ORDERS_CHILDREN_ROUTES = {
	ADMIN_STORE_ORDERS: {
		path: 'store-orders',
		meta: {
			title: 'Все заказы кафе',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/store-orders/pages/admin-store-orders-page.vue'),
	},
} satisfies AppRouteRecord
