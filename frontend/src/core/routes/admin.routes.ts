import { type ExtendedRouteRecord } from '../config/routes.config'
import { AppLayouts } from '../types/routes.types'

export const ADMIN_ROUTES_CONFIG = {
	ADMIN_DASHBOARD: {
		path: '/admin',
		meta: { layout: AppLayouts.ADMIN, title: 'Dashboard', requiresAuth: true },
		component: () => import('@/modules/admin/pages/admin-dashboard-page.vue'),
	},
	ADMIN_ANALYTICS: {
		path: '/admin/analytics',
		meta: { layout: AppLayouts.ADMIN, title: 'Analytics', requiresAuth: true },
		component: () => import('@/modules/admin/pages/admin-analytics-page.vue'),
	},
} satisfies ExtendedRouteRecord
