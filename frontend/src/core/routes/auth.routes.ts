import LoginPage from '@/modules/auth/pages/login-page.vue'
import type { AppRouteRecord, ParentRoutePage } from './routes.types'
import AppAuthLayout from '../layouts/auth/app-auth-layout.vue'

export const AUTH_CHILDREN_ROUTES = {
	LOGIN: {
		path: '',
		meta: {
			title: 'Login',
			requiresAuth: false,
		},
		component: LoginPage,
	},
} satisfies AppRouteRecord

export const AUTH_ROUTES_CONFIG: ParentRoutePage = {
	path: '/',
	component: AppAuthLayout,
	children: AUTH_CHILDREN_ROUTES,
}
