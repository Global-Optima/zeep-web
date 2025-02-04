import type { AppRouteRecord } from '../routes.types'

export const ADMIN_WAREHOUSES_CHILDREN_ROUTES = {
	ADMIN_WAREHOUSES: {
		path: 'warehouses',
		meta: {
			title: 'Склады',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/warehouses/pages/admin-warehouses-page.vue'),
	},
	ADMIN_WAREHOUSE_CREATE: {
		path: 'warehouses/create',
		meta: {
			title: 'Создать склад',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/warehouses/pages/admin-warehouse-create-page.vue'),
	},
	ADMIN_WAREHOUSE_DETAILS: {
		path: 'warehouses/:id',
		meta: {
			title: 'Детали склада',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/warehouses/pages/admin-warehouse-details-page.vue'),
	},
	ADMIN_WAREHOUSE_STOCKS: {
		path: 'warehouse-stocks',
		meta: {
			title: 'Запасы Склада',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/warehouse-stocks/pages/admin-warehouse-stocks-page.vue'),
	},
	ADMIN_WAREHOUSE_STOCK_DETAILS: {
		path: 'warehouse-stocks/:id',
		meta: {
			title: 'Детали материала склада',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/warehouse-stocks/pages/admin-warehouse-stock-details-page.vue'),
	},
	ADMIN_WAREHOUSE_STOCKS_CREATE: {
		path: 'warehouse-stocks/create',
		meta: {
			title: 'Добавить материалы на склад',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/warehouse-stocks/pages/admin-warehouse-stock-create-page.vue'),
	},
	ADMIN_WAREHOUSE_DELIVERIES: {
		path: 'warehouse-deliveries',
		meta: {
			title: 'Доставки на склад',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/warehouse-deliveries/pages/admin-warehouse-deliveries-page.vue'),
	},
	ADMIN_WAREHOUSE_DELIVERIES_DETAILS: {
		path: 'warehouse-deliveries/:id',
		meta: {
			title: 'Детали доставки на склад',
			requiresAuth: true,
		},
		component: () =>
			import(
				'@/modules/admin/warehouse-deliveries/pages/admin-warehouse-deliveries-details-page.vue'
			),
	},
	ADMIN_WAREHOUSE_DELIVERIES_CREATE: {
		path: 'warehouse-deliveries/create',
		meta: {
			title: 'Создать доставку на склад',
			requiresAuth: true,
		},
		component: () =>
			import(
				'@/modules/admin/warehouse-deliveries/pages/admin-warehouse-deliveries-create-page.vue'
			),
	},
	ADMIN_WAREHOUSE_STOCK_REQUESTS: {
		path: 'stock-requests',
		meta: {
			title: 'Запросы на склад',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/stock-requests/pages/admin-stock-requests-list-page.vue'),
	},
	ADMIN_WAREHOUSE_STOCK_REQUESTS_DETAILS: {
		path: 'stock-requests/:id',
		meta: {
			title: 'Детали запроса на склад',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/stock-requests/pages/admin-stock-requests-details-page.vue'),
	},
} satisfies AppRouteRecord
