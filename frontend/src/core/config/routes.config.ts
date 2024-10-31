import type { RouteRecordNameGeneric, RouteRecordRaw } from 'vue-router'
import { ERRORS_ROUTES_CONFIG } from '../routes/errors.routes'
import { KIOSK_ROUTES_CONFIG } from '@/core/routes/kiosk.routes'
import { ADMIN_ROUTES_CONFIG } from '@/core/routes/admin.routes'
import type { AppLayouts } from '../types/routes.types'
import { AUTH_ROUTES_CONFIG } from '../routes/auth.routes'

export type ExtendedRoutePage = RouteRecordRaw & {
	meta: {
		layout?: AppLayouts
		title: string
		pageTitle?: string
		requiresAuth?: boolean
	}
	children?: ExtendedRoutePage[]
}

export type ExtendedRouteRecord = Record<string, ExtendedRoutePage>

const ROUTES_RECORDS = addNameToRoutes([
	KIOSK_ROUTES_CONFIG,
	ADMIN_ROUTES_CONFIG,
	AUTH_ROUTES_CONFIG,

	// Always in the bottom
	ERRORS_ROUTES_CONFIG,
] as const)

function addNameToRoutes<T extends readonly ExtendedRouteRecord[]>(records: T): T {
	return records.map(record => {
		const updatedRecord: ExtendedRouteRecord = {}
		Object.entries(record).forEach(([key, value]) => {
			updatedRecord[key] = { ...value, name: key }
		})
		return updatedRecord
	}) as unknown as T
}

type ExtractRouteKeys<T extends readonly ExtendedRouteRecord[]> = T[number] extends infer R
	? R extends ExtendedRouteRecord
		? keyof R
		: never
	: never

export type RouteKey = ExtractRouteKeys<typeof ROUTES_RECORDS>

export const ROUTES: ExtendedRoutePage[] = ROUTES_RECORDS.flatMap(record => {
	return Object.values(record)
})

const routeLookupMap: Record<RouteKey, ExtendedRoutePage | undefined> = ROUTES_RECORDS.reduce(
	(acc, record) => {
		Object.entries(record).forEach(([key, value]) => {
			acc[key as RouteKey] = value
		})
		return acc
	},
	{} as Record<RouteKey, ExtendedRoutePage | undefined>,
)

export function getRoute(key: RouteKey): ExtendedRoutePage | undefined {
	return routeLookupMap[key]
}

export function getRouteName(key: RouteKey): RouteRecordNameGeneric | undefined {
	return routeLookupMap[key]?.name
}
