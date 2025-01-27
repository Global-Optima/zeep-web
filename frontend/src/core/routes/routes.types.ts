import type { RouteComponent } from 'vue-router'

type LazyRoute = () => Promise<RouteComponent>

export type ParentRoutePage = {
	path: string
	component: RouteComponent | LazyRoute
	children: AppRouteRecord
}

export type AppRoutePage = {
	path: string
	component: RouteComponent | LazyRoute
	name?: string
	meta: {
		title: string
		pageTitle?: string
		requiresAuth?: boolean
	}
}

export type AppRouteRecord = Record<string, AppRoutePage>
