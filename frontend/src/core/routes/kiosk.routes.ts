import type { ExtendedRouteRecord } from '../config/routes.config'

export const PRIMARY_ROUTES_CONFIG = {
	KIOSK_NAME: {
		path: '/kiosk',
		meta: { layout: 'KIOSK', title: 'Home', requiresAuth: true },
		component: () => import('@/modules/kiosk/pages/kiosk-home-page.vue'),
	},
} satisfies ExtendedRouteRecord
