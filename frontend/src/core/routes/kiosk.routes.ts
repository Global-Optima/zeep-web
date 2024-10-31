import { type ExtendedRouteRecord } from '@/core/config/routes.config'
import { AppLayouts } from '../types/routes.types'

export const KIOSK_ROUTES_CONFIG = {
	KIOSK_HOME: {
		path: 'kiosk',
		meta: { layout: AppLayouts.KIOSK, title: 'Home', requiresAuth: true },
		component: () => import('@/modules/kiosk/pages/kiosk-home-page.vue'),
	},
} satisfies ExtendedRouteRecord
