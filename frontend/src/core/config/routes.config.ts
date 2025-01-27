import { ADMIN_CHILDREN_ROUTES, ADMIN_ROUTES_CONFIG } from '@/core/routes/admin/admin.routes'
import { KIOSK_CHILDREN_ROUTES, KIOSK_ROUTES_CONFIG } from '@/core/routes/kiosk.routes'
import type { RouteRecordRaw } from 'vue-router'
import { AUTH_CHILDREN_ROUTES, AUTH_ROUTES_CONFIG } from '../routes/auth.routes'
import { ERRORS_CHILDREN_ROUTES, ERRORS_ROUTES_CONFIG } from '../routes/errors.routes'
import type { AppRoutePage, AppRouteRecord, ParentRoutePage } from '../routes/routes.types'

const PARENT_ROUTES_RECORDS = addNameToChildrenInParent([
	KIOSK_ROUTES_CONFIG,
	AUTH_ROUTES_CONFIG,
	ERRORS_ROUTES_CONFIG,
	ADMIN_ROUTES_CONFIG,
] as const)

const CHILDREN_ROUTES_RECORDS = addNameToChildren([
	KIOSK_CHILDREN_ROUTES,
	AUTH_CHILDREN_ROUTES,
	ERRORS_CHILDREN_ROUTES,
	ADMIN_CHILDREN_ROUTES,
])

function addNameToChildrenInParent(parentRoutes: ParentRoutePage[]): ParentRoutePage[] {
	return parentRoutes.map(parent => {
		const clonedParent = { ...parent }

		if (clonedParent.children) {
			const updatedChildren = Object.entries(clonedParent.children).reduce((acc, [key, value]) => {
				acc[key] = { ...value, name: key }
				return acc
			}, {} as AppRouteRecord)

			clonedParent.children = updatedChildren
		}

		return clonedParent
	})
}

function addNameToChildren<T extends readonly AppRouteRecord[]>(records: T): T {
	return records.map(record => {
		const updatedRecord: AppRouteRecord = {}
		Object.entries(record).forEach(([key, value]) => {
			updatedRecord[key] = { ...value, name: key }
		})
		return updatedRecord
	}) as unknown as T
}

type ExtractRouteKeys<T extends readonly AppRouteRecord[]> = T[number] extends infer R
	? R extends AppRouteRecord
		? keyof R
		: never
	: never

export type RouteKey = ExtractRouteKeys<typeof CHILDREN_ROUTES_RECORDS>

export const ROUTES: RouteRecordRaw[] = PARENT_ROUTES_RECORDS.map(record => {
	const route: RouteRecordRaw = {
		path: record.path,
		component: record.component,
		children: Object.values(record.children),
	}

	return route
})

const routeLookupMap: Record<RouteKey, AppRoutePage> = CHILDREN_ROUTES_RECORDS.reduce(
	(acc, record) => {
		Object.entries(record).forEach(([key, value]) => {
			acc[key as RouteKey] = value
		})
		return acc
	},
	{} as Record<RouteKey, AppRoutePage>,
)

export function getRoute(key: RouteKey): AppRoutePage {
	return routeLookupMap[key]
}

export function getRouteName(key: RouteKey): string {
	return routeLookupMap[key].name!
}

export function getRoutePath(key: RouteKey): string {
	return routeLookupMap[key].path
}
