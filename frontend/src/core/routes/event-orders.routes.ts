import type { AppRouteRecord, ParentRoutePage } from './routes.types'

export const EVENT_ORDERS_CHILDREN_ROUTES = {
	EVENT_ORDERS_BARISTA: {
		path: 'barista-orders',
		meta: {
			title: 'Дэшборд заказов кафе',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/store-orders/barista/pages/admin-barista-orders-page.vue'),
	},
	EVENT_ORDERS_DISPLAY: {
		path: 'orders-display',
		meta: {
			title: 'Заказы',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/store-orders/display/pages/admin-orders-display-page.vue'),
	},
} satisfies AppRouteRecord

export const EVENT_ORDERS_ROUTES_CONFIG: ParentRoutePage = {
	path: '/',
	component: () => import('@/core/layouts/default/app-default-layout.vue'),
	children: EVENT_ORDERS_CHILDREN_ROUTES,
}
