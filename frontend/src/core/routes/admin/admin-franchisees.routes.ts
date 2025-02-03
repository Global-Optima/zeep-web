import type { AppRouteRecord } from '../routes.types'

export const ADMIN_FRANCHISEES_CHILDREN_ROUTES = {
	ADMIN_FRANCHISEES: {
		path: 'franchisees',
		meta: {
			title: 'Франчайзи',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/franchisees/pages/admin-franchisees-page.vue'),
	},
	ADMIN_FRANCHISEE_DETAILS: {
		path: 'franchisees/:id',
		meta: {
			title: 'Детали франчайзи',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/franchisees/pages/admin-franchisee-details-page.vue'),
	},
	ADMIN_FRANCHISEE_CREATE: {
		path: 'franchisees/create',
		meta: {
			title: 'Создать франчайзи',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/franchisees/pages/admin-franchisee-create-page.vue'),
	},
} satisfies AppRouteRecord
