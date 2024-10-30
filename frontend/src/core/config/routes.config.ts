import { PRIMARY_ROUTES_CONFIG } from '@/core/routes/kiosk.routes'
import type { RouteRecordRaw } from 'vue-router'
import { ERRORS_ROUTES_CONFIG } from '../routes/errors.routes'

export type ExtendedRoutePage = RouteRecordRaw & {
	meta: {
		layout?: string
		title: string
		pageTitle?: string
		requiresAuth?: boolean
	}
	children?: ExtendedRoutePage[]
}

export type ExtendedRouteRecord = Record<string, ExtendedRoutePage>

const ROUTES_RECORDS: ExtendedRouteRecord[] = addNameToRoutes([
	PRIMARY_ROUTES_CONFIG,

	// Always in the bottom
	ERRORS_ROUTES_CONFIG,
])

function addNameToRoutes(records: ExtendedRouteRecord[]): ExtendedRouteRecord[] {
	return records.map(record => {
		const updatedRecord: ExtendedRouteRecord = {}
		Object.entries(record).forEach(([key, value]) => {
			updatedRecord[key] = { ...value, name: key }
		})
		return updatedRecord
	})
}

export const ROUTES: ExtendedRoutePage[] = ROUTES_RECORDS.flatMap(record => {
	return Object.values(record)
})

type RouteKey = keyof typeof ERRORS_ROUTES_CONFIG | keyof typeof PRIMARY_ROUTES_CONFIG

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
