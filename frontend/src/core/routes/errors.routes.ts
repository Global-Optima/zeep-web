import type { ExtendedRouteRecord } from '../config/routes.config'

export const ERRORS_ROUTES_CONFIG = {
	NOT_FOUND_ERROR: {
		path: ':pathMatch(.*)*',
		name: 'not-found',
		meta: { title: 'Not found' },
		component: () => import('@/modules/errors/pages/not-found-page.vue'),
	},
	INTERNAL_ERROR: {
		path: ':pathMatch(.*)*',
		name: 'internal-error',
		meta: { title: 'Internal server error' },
		component: () => import('@/modules/errors/pages/internal-error-page.vue'),
	},
} satisfies ExtendedRouteRecord
