import LoginPage from '@/modules/auth/pages/login-page.vue'
import { type ExtendedRouteRecord } from '../config/routes.config'

export const AUTH_ROUTES_CONFIG = {
	LOGIN: {
		path: '',
		meta: { title: 'Login', requiresAuth: false },
		component: LoginPage,
	},
} satisfies ExtendedRouteRecord
