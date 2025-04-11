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
	ADMIN_STORE_ORDER_DETEAILS: {
		path: 'store-orders/:id',
		meta: {
			title: 'Заказ кафе',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/store-orders/pages/admin-store-order-details-page.vue'),
	},
} satisfies AppRouteRecord
