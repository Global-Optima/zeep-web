import AppDefaultLayout from '@/core/layouts/default/app-default-layout.vue'
import type { AppRouteRecord, ParentRoutePage } from './routes.types'

export const ERRORS_CHILDREN_ROUTES = {
	NOT_FOUND_ERROR: {
		path: 'not-found',
		name: 'not-found',
		meta: {
			title: 'Not found',
		},
		component: () => import('@/modules/errors/pages/not-found-page.vue'),
	},
	INTERNAL_ERROR: {
		path: 'internal-error',
		name: 'internal-error',
		meta: {
			title: 'Internal server error',
		},
		component: () => import('@/modules/errors/pages/internal-error-page.vue'),
	},
} satisfies AppRouteRecord

export const ERRORS_ROUTES_CONFIG: ParentRoutePage = {
	path: '/errors',
	component: AppDefaultLayout,
	children: ERRORS_CHILDREN_ROUTES,
}
