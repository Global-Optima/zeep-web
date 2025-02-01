import LoginPage from '@/modules/auth/pages/login-page-v2.vue'
import AppDefaultLayout from '../layouts/default/app-default-layout.vue'
import type { AppRouteRecord, ParentRoutePage } from './routes.types'

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
