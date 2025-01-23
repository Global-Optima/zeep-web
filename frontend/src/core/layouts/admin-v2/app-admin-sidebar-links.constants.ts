import type { RouteKey } from '@/core/config/routes.config'
import { EmployeeRole } from '@/modules/admin/store-employees/models/employees.models'
import {
	Blocks,
	ChartBar,
	FolderKanban,
	LayoutList,
	Package,
	PackageCheck,
	Ruler,
	ShoppingBasket,
	Truck,
	Users,
	Warehouse,
	type LucideIcon,
} from 'lucide-vue-next'

// Base interface for a single navigation link
export interface NavItem {
	name: string
	routeKey: RouteKey
	icon: LucideIcon
	accessRoles: EmployeeRole[]
}

// Interface for collapsible navigation groups
export interface CollapsibleNavItem {
	label: string
	icon: LucideIcon
	accessRoles: EmployeeRole[]
	items: NavItem[]
}

// Union type to allow normal or collapsible items
export type SidebarNavItem = NavItem | CollapsibleNavItem

// Utility to check if an item is collapsible
export const isCollapsibleNavItem = (item: SidebarNavItem): item is CollapsibleNavItem => {
	return 'items' in item
}

export const adminNavItems: SidebarNavItem[] = [
	// Normal links
	{
		name: 'Аналитика',
		routeKey: 'ADMIN_DASHBOARD',
		icon: ChartBar,
		accessRoles: [EmployeeRole.ADMIN, EmployeeRole.DIRECTOR],
	},
	{
		name: 'Аналитика',
		routeKey: 'ADMIN_STORE_DASHBOARD',
		icon: ChartBar,
		accessRoles: [EmployeeRole.MANAGER, EmployeeRole.BARISTA],
	},
	{
		name: 'Сотрудники',
		routeKey: 'ADMIN_STORE_EMPLOYEES',
		icon: Users,
		accessRoles: [EmployeeRole.MANAGER],
	},
	{
		name: 'Запасы',
		routeKey: 'ADMIN_STORE_STOCKS',
		icon: Warehouse,
		accessRoles: [EmployeeRole.MANAGER, EmployeeRole.BARISTA],
	},
	{
		name: 'Проекты',
		routeKey: 'ADMIN_PROJECTS',
		icon: FolderKanban,
		accessRoles: [EmployeeRole.MANAGER, EmployeeRole.BARISTA],
	},
	{
		name: 'Размеры',
		routeKey: 'ADMIN_UNITS',
		icon: Ruler,
		accessRoles: [EmployeeRole.ADMIN, EmployeeRole.DIRECTOR],
	},
	{
		name: 'Поставщики',
		routeKey: 'ADMIN_SUPPLIERS',
		icon: Truck,
		accessRoles: [EmployeeRole.ADMIN],
	},
	{
		name: 'Запросы на склад',
		routeKey: 'ADMIN_STORE_STOCK_REQUESTS',
		icon: Truck,
		accessRoles: [EmployeeRole.MANAGER, EmployeeRole.BARISTA],
	},
	{
		name: 'Запросы на склад',
		routeKey: 'ADMIN_WAREHOUSE_STOCK_REQUESTS',
		icon: Truck,
		accessRoles: [EmployeeRole.WAREHOUSE],
	},
	{
		name: 'Запасы',
		routeKey: 'ADMIN_WAREHOUSE_STOCKS',
		icon: Blocks,
		accessRoles: [EmployeeRole.WAREHOUSE],
	},
	{
		name: 'Доставки',
		routeKey: 'ADMIN_WAREHOUSE_DELIVERIES',
		icon: PackageCheck,
		accessRoles: [EmployeeRole.WAREHOUSE],
	},
	// Collapsible groups
	{
		label: 'Каталог',
		icon: ShoppingBasket,
		accessRoles: [EmployeeRole.ADMIN, EmployeeRole.DIRECTOR],
		items: [
			{
				name: 'Список',
				routeKey: 'ADMIN_PRODUCTS',
				icon: Package,
				accessRoles: [EmployeeRole.ADMIN, EmployeeRole.DIRECTOR],
			},
			{
				name: 'Категории',
				routeKey: 'ADMIN_PRODUCT_CATEGORIES',
				icon: LayoutList,
				accessRoles: [EmployeeRole.ADMIN, EmployeeRole.DIRECTOR],
			},
		],
	},

	{
		label: 'Складские товары',
		icon: Package,
		accessRoles: [EmployeeRole.ADMIN, EmployeeRole.DIRECTOR],
		items: [
			{
				name: 'Список',
				routeKey: 'ADMIN_STOCK_MATERIALS',
				icon: Package,
				accessRoles: [EmployeeRole.ADMIN, EmployeeRole.DIRECTOR],
			},
			{
				name: 'Категории',
				routeKey: 'ADMIN_STOCK_MATERIAL_CATEGORIES',
				icon: LayoutList,
				accessRoles: [EmployeeRole.ADMIN, EmployeeRole.DIRECTOR],
			},
		],
	},
]
