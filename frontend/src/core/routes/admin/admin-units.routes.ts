import type { AppRouteRecord } from '../../config/routes.config'

export const ADMIN_UNITS_CHILDREN_ROUTES = {
	ADMIN_UNITS: {
		path: 'units',
		meta: {
			title: 'Размеры',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/units/pages/admin-units-page.vue'),
	},
	ADMIN_UNIT_CREATE: {
		path: 'units/create',
		meta: {
			title: 'Создать размер',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/units/pages/admin-unit-create-page.vue'),
	},
	ADMIN_UNIT_DETAILS: {
		path: 'units/:id',
		meta: {
			title: 'Детали размера',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/units/pages/admin-unit-details-page.vue'),
	},
} satisfies AppRouteRecord
