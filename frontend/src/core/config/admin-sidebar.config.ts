import { type RouteKey } from './routes.config'

interface MenuItem {
	routeName: RouteKey
	icon: string
}

export const adminSidebarMenu: MenuItem[] = [
	{ routeName: 'ADMIN_DASHBOARD', icon: 'mingcute:home-4-line' },
	{ routeName: 'ADMIN_ANALYTICS', icon: 'mingcute:chart-pie-2-line' },
]
