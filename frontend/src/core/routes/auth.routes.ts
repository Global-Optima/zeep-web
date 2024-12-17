import LoginPage from '@/modules/auth/pages/login-page.vue'
import type { AppRouteRecord, ParentRoutePage } from '../config/routes.config'
import AppDefaultLayout from '../layouts/default/app-default-layout.vue'

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
	component: AppDefaultLayout,
	children: AUTH_CHILDREN_ROUTES,
}
