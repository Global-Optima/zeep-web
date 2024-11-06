import type { AppRouteRecord, ParentRoutePage } from '../config/routes.config'
import AppAdminLayout from '../layouts/admin/app-admin-layout.vue'

export const ADMIN_CHILDREN_ROUTES = {
	ADMIN_DASHBOARD: {
		path: '',
		meta: {
			title: 'Dashboard',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/pages/admin-dashboard-page.vue'),
	},
	ADMIN_ANALYTICS: {
		path: 'analytics',
		meta: {
			title: 'Analytics',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/pages/admin-analytics-page.vue'),
	},
} satisfies AppRouteRecord

export const ADMIN_ROUTES_CONFIG = {
	path: 'admin',
	component: AppAdminLayout,
	children: ADMIN_CHILDREN_ROUTES,
} satisfies ParentRoutePage
