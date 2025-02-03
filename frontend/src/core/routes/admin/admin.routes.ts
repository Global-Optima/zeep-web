import AppAdminLayout from '@/core/layouts/admin/app-admin-layout.vue'
import { ADMIN_ADDITIVES_CHILDREN_ROUTES } from '@/core/routes/admin/admin-additives.routes'
import { ADMIN_INGREDIENTS_CHILDREN_ROUTES } from '@/core/routes/admin/admin-ingredients.routes'
import { ADMIN_PRODUCTS_CHILDREN_ROUTES } from '@/core/routes/admin/admin-products.routes'
import { ADMIN_REGIONS_CHILDREN_ROUTES } from '@/core/routes/admin/admin-regions.routes'
import { ADMIN_STOCK_MATERIALS_CHILDREN_ROUTES } from '@/core/routes/admin/admin-stock-materials.routes'
import { ADMIN_SUPPLIERS_CHILDREN_ROUTES } from '@/core/routes/admin/admin-suppliers.routes'
import { ADMIN_UNITS_CHILDREN_ROUTES } from '@/core/routes/admin/admin-units.routes'
import type { AppRouteRecord, ParentRoutePage } from '../routes.types'
import { ADMIN_EMPLOYEES_CHILDREN_ROUTES } from './admin-employees.routes'

export const ADMIN_CHILDREN_ROUTES = {
	...ADMIN_ADDITIVES_CHILDREN_ROUTES,
	...ADMIN_PRODUCTS_CHILDREN_ROUTES,
	...ADMIN_INGREDIENTS_CHILDREN_ROUTES,
	...ADMIN_UNITS_CHILDREN_ROUTES,
	...ADMIN_STOCK_MATERIALS_CHILDREN_ROUTES,
	...ADMIN_EMPLOYEES_CHILDREN_ROUTES,
	...ADMIN_SUPPLIERS_CHILDREN_ROUTES,
	...ADMIN_REGIONS_CHILDREN_ROUTES,

	ADMIN_NOTIFICATIONS: {
		path: 'notifications',
		meta: {
			title: 'Уведомления',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/notifications/pages/admin-notifications-page.vue'),
	},
	ADMIN_DASHBOARD: {
		path: '',
		meta: {
			title: 'Аналитика',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/dashboard/pages/admin-dashboard-page.vue'),
	},
	ADMIN_STORE_DASHBOARD: {
		path: 'store-analytics',
		meta: {
			title: 'Аналитика магазина',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/dashboard/pages/admin-dashboard-page.vue'),
	},
	ADMIN_STORE_ORDERS: {
		path: 'store-orders',
		meta: {
			title: 'Заказы магазина',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/store-orders/pages/admin-store-orders-page.vue'),
	},
	ADMIN_STORE_STOCKS: {
		path: 'store-stocks',
		meta: {
			title: 'Склад Магазина',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/store-stocks/pages/admin-store-stocks-page.vue'),
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
	ADMIN_STORE_STOCKS_DETAILS: {
		path: 'store-stocks/:id',
		meta: {
			title: 'Детали запаса',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/store-stocks/pages/admin-store-stocks-details-page.vue'),
	},
	ADMIN_CREATE_STORE_STOCKS: {
		path: 'store-stocks/create',
		meta: {
			title: 'Добавить в склад магазина',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/store-stocks/pages/admin-store-stocks-create-page.vue'),
	},
	ADMIN_STORES: {
		path: 'stores',
		meta: {
			title: 'Магазины',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/stores/pages/admin-stores-page.vue'),
	},
	ADMIN_STORE_CREATE: {
		path: 'stores/create',
		meta: {
			title: 'Добавить магазин',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/stores/pages/admin-stores-create-page.vue'),
	},
	ADMIN_STORE_DETAILS: {
		path: 'stores/:id',
		meta: {
			title: 'Детали магазина',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/stores/pages/admin-store-details-page.vue'),
	},
	ADMIN_STORE_STOCK_REQUESTS: {
		path: 'stock-requests',
		meta: {
			title: 'Запросы на склад',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/stock-requests/pages/admin-stock-requests-list-page.vue'),
	},
	ADMIN_STORE_STOCK_REQUESTS_CREATE: {
		path: 'stock-requests/create',
		meta: {
			title: 'Создать запрос на склад',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/stock-requests/pages/admin-stock-requests-create-page.vue'),
	},
	ADMIN_STORE_STOCK_REQUESTS_UPDATE: {
		path: 'stock-requests/:id/update',
		meta: {
			title: 'Обновить запрос на склад',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/stock-requests/pages/admin-stock-requests-update-page.vue'),
	},
	ADMIN_STORE_STOCK_REQUESTS_DETAILS: {
		path: 'stock-requests/:id',
		meta: {
			title: 'Детали запроса на склад',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/stock-requests/pages/admin-stock-requests-details-page.vue'),
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

export const ADMIN_ROUTES_CONFIG = {
	path: '/admin',
	component: AppAdminLayout,
	children: ADMIN_CHILDREN_ROUTES,
} satisfies ParentRoutePage
