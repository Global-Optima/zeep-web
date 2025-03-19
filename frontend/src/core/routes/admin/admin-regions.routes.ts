import type { AppRouteRecord } from '../routes.types'

export const ADMIN_REGIONS_CHILDREN_ROUTES = {
	ADMIN_REGIONS: {
		path: 'regions',
		meta: {
			title: 'Регионы',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/regions/pages/admin-regions-page.vue'),
	},
	ADMIN_REGION_DETAILS: {
		path: 'regions/:id',
		meta: {
			title: 'Детали региона',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/regions/pages/admin-region-details-page.vue'),
	},
	ADMIN_REGION_CREATE: {
		path: 'regions/create',
		meta: {
			title: 'Создать регион',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/regions/pages/admin-region-create-page.vue'),
	},
} satisfies AppRouteRecord
