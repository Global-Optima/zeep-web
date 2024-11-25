import type { AppRouteRecord, ParentRoutePage } from '../config/routes.config'
import AppAdminLayout from '../layouts/admin/app-admin-layout.vue'

export const ADMIN_CHILDREN_ROUTES = {
	ADMIN_DASHBOARD: {
		path: '',
		meta: {
			title: 'Главная',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/pages/admin-dashboard-page.vue'),
	},
	ADMIN_ANALYTICS: {
		path: 'analytics',
		meta: {
			title: 'Аналитика',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/pages/admin-analytics-page.vue'),
	},
	ADMIN_ORDERS: {
		path: 'orders',
		meta: {
			title: 'Заказы',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/pages/admin-orders-page.vue'),
	},
	ADMIN_WAREHOUSE: {
		path: 'warehouse',
		meta: {
			title: 'Склад',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/pages/admin-warehouse-page.vue'),
	},
	ADMIN_PRODUCTS: {
		path: 'products',
		meta: {
			title: 'Товары',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/pages/admin-products-page.vue'),
	},
	ADMIN_EMPLOYEES: {
		path: 'employees',
		meta: {
			title: 'Сотрудники',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/pages/admin-employees-page.vue'),
	},
	ADMIN_STORES: {
		path: 'stores',
		meta: {
			title: 'Магазины',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/pages/admin-stores-page.vue'),
	},
} satisfies AppRouteRecord

export const ADMIN_ROUTES_CONFIG = {
	path: 'admin',
	component: AppAdminLayout,
	children: ADMIN_CHILDREN_ROUTES,
} satisfies ParentRoutePage
