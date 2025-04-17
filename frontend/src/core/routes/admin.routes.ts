import { ADMIN_ADDITIVES_CHILDREN_ROUTES } from '@/core/routes/admin/admin-additives.routes'
import { ADMIN_FRANCHISEES_CHILDREN_ROUTES } from '@/core/routes/admin/admin-franchisees.routes'
import { ADMIN_INGREDIENTS_CHILDREN_ROUTES } from '@/core/routes/admin/admin-ingredients.routes'
import { ADMIN_PRODUCTS_CHILDREN_ROUTES } from '@/core/routes/admin/admin-products.routes'
import { ADMIN_REGIONS_CHILDREN_ROUTES } from '@/core/routes/admin/admin-regions.routes'
import { ADMIN_STOCK_MATERIALS_CHILDREN_ROUTES } from '@/core/routes/admin/admin-stock-materials.routes'
import { ADMIN_SUPPLIERS_CHILDREN_ROUTES } from '@/core/routes/admin/admin-suppliers.routes'
import { ADMIN_UNITS_CHILDREN_ROUTES } from '@/core/routes/admin/admin-units.routes'
import { ADMIN_WAREHOUSES_CHILDREN_ROUTES } from '@/core/routes/admin/admin-warehouses.routes'
import { ADMIN_EMPLOYEES_CHILDREN_ROUTES } from './admin/admin-employees.routes'
import { ADMIN_STORE_ORDERS_CHILDREN_ROUTES } from './admin/admin-orders.routes'
import { ADMIN_PROVISIONS_CHILDREN_ROUTES } from './admin/admin-provisions.routes'
import type { AppRouteRecord, ParentRoutePage } from './routes.types'

export const ADMIN_CHILDREN_ROUTES = {
	...ADMIN_ADDITIVES_CHILDREN_ROUTES,
	...ADMIN_PRODUCTS_CHILDREN_ROUTES,
	...ADMIN_INGREDIENTS_CHILDREN_ROUTES,
	...ADMIN_UNITS_CHILDREN_ROUTES,
	...ADMIN_STOCK_MATERIALS_CHILDREN_ROUTES,
	...ADMIN_EMPLOYEES_CHILDREN_ROUTES,
	...ADMIN_SUPPLIERS_CHILDREN_ROUTES,
	...ADMIN_REGIONS_CHILDREN_ROUTES,
	...ADMIN_FRANCHISEES_CHILDREN_ROUTES,
	...ADMIN_WAREHOUSES_CHILDREN_ROUTES,
	...ADMIN_STORE_ORDERS_CHILDREN_ROUTES,
	...ADMIN_PROVISIONS_CHILDREN_ROUTES,

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
			title: 'Аналитика кафе',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/dashboard/pages/admin-dashboard-page.vue'),
	},

	ADMIN_STORE_STOCKS: {
		path: 'store-stocks',
		meta: {
			title: 'Склад кафе',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/store-stocks/pages/admin-store-stocks-page.vue'),
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
			title: 'Добавить в склад кафе',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/store-stocks/pages/admin-store-stocks-create-page.vue'),
	},
	ADMIN_STORES: {
		path: 'stores',
		meta: {
			title: 'Кафе',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/stores/pages/admin-stores-page.vue'),
	},
	ADMIN_STORE_SETTINGS: {
		path: 'stores/settings',
		meta: {
			title: 'Настройки',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/settings/stores/pages/admin-store-settings-page.vue'),
	},

	ADMIN_STORE_CREATE: {
		path: 'stores/create',
		meta: {
			title: 'Добавить кафе',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/stores/pages/admin-stores-create-page.vue'),
	},
	ADMIN_STORE_DETAILS: {
		path: 'stores/:id',
		meta: {
			title: 'Детали кафе',
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
} satisfies AppRouteRecord

export const ADMIN_ROUTES_CONFIG = {
	path: '/admin',
	component: () => import('@/core/layouts/admin/app-admin-layout.vue'),
	children: ADMIN_CHILDREN_ROUTES,
} satisfies ParentRoutePage
