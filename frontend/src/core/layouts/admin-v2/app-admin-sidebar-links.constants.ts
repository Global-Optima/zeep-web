import type { RouteKey } from '@/core/config/routes.config'
import { EmployeeRole } from '@/modules/admin/store-employees/models/employees.models'
import {
	Apple,
	Blocks,
	ChartBar,
	Globe,
	LayoutList,
	ListPlus,
	MapPinned,
	Package,
	PackageCheck,
	Ruler,
	ShoppingCart,
	Store,
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
		accessRoles: [EmployeeRole.ADMIN],
	},
	{
		name: 'Аналитика',
		routeKey: 'ADMIN_STORE_DASHBOARD',
		icon: ChartBar,
		accessRoles: [EmployeeRole.STORE_MANAGER, EmployeeRole.BARISTA],
	},
	{
		name: 'Заказы',
		routeKey: 'ADMIN_STORE_ORDERS',
		icon: ShoppingCart,
		accessRoles: [EmployeeRole.STORE_MANAGER, EmployeeRole.BARISTA],
	},
	{
		name: 'Сотрудники',
		routeKey: 'ADMIN_STORE_EMPLOYEES',
		icon: Users,
		accessRoles: [EmployeeRole.STORE_MANAGER],
	},
	{
		name: 'Склад',
		routeKey: 'ADMIN_STORE_STOCKS',
		icon: Warehouse,
		accessRoles: [EmployeeRole.STORE_MANAGER, EmployeeRole.BARISTA],
	},
	{
		name: 'Товары',
		routeKey: 'ADMIN_STORE_PRODUCTS',
		icon: Package,
		accessRoles: [EmployeeRole.STORE_MANAGER, EmployeeRole.BARISTA],
	},
	{
		name: 'Топпинги',
		routeKey: 'ADMIN_STORE_ADDITIVES',
		icon: ListPlus,
		accessRoles: [EmployeeRole.STORE_MANAGER, EmployeeRole.BARISTA],
	},
	{
		name: 'Магазины',
		routeKey: 'ADMIN_STORES',
		icon: Store,
		accessRoles: [EmployeeRole.ADMIN],
	},
	{
		name: 'Размеры',
		routeKey: 'ADMIN_UNITS',
		icon: Ruler,
		accessRoles: [EmployeeRole.ADMIN],
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
		accessRoles: [EmployeeRole.STORE_MANAGER, EmployeeRole.BARISTA],
	},
	{
		name: 'Поставщики',
		routeKey: 'ADMIN_WAREHOUSE_SUPPLIERS',
		icon: Users,
		accessRoles: [EmployeeRole.WAREHOUSE_MANAGER, EmployeeRole.WAREHOUSE_EMPLOYEE],
	},
	{
		name: 'Запросы на склад',
		routeKey: 'ADMIN_WAREHOUSE_STOCK_REQUESTS',
		icon: Truck,
		accessRoles: [EmployeeRole.WAREHOUSE_MANAGER, EmployeeRole.WAREHOUSE_EMPLOYEE],
	},
	{
		name: 'Запасы',
		routeKey: 'ADMIN_WAREHOUSE_STOCKS',
		icon: Blocks,
		accessRoles: [EmployeeRole.WAREHOUSE_MANAGER, EmployeeRole.WAREHOUSE_EMPLOYEE],
	},
	{
		name: 'Доставки',
		routeKey: 'ADMIN_WAREHOUSE_DELIVERIES',
		icon: PackageCheck,
		accessRoles: [EmployeeRole.WAREHOUSE_MANAGER, EmployeeRole.WAREHOUSE_EMPLOYEE],
	},
	{
		name: 'Регионы',
		routeKey: 'ADMIN_REGIONS',
		icon: MapPinned,
		accessRoles: [EmployeeRole.ADMIN],
	},
	// Collapsible groups
	{
		label: 'Товары',
		icon: Package,
		accessRoles: [EmployeeRole.ADMIN],
		items: [
			{
				name: 'Список',
				routeKey: 'ADMIN_PRODUCTS',
				icon: Package,
				accessRoles: [EmployeeRole.ADMIN],
			},
			{
				name: 'Категории',
				routeKey: 'ADMIN_PRODUCT_CATEGORIES',
				icon: LayoutList,
				accessRoles: [EmployeeRole.ADMIN],
			},
		],
	},

	{
		label: 'Топпинги',
		icon: ListPlus,
		accessRoles: [EmployeeRole.ADMIN],
		items: [
			{
				name: 'Список',
				routeKey: 'ADMIN_ADDITIVES',
				icon: ListPlus,
				accessRoles: [EmployeeRole.ADMIN],
			},
			{
				name: 'Категории',
				routeKey: 'ADMIN_ADDITIVE_CATEGORIES',
				icon: LayoutList,
				accessRoles: [EmployeeRole.ADMIN],
			},
		],
	},
	{
		label: 'Ингредиенты',
		icon: Apple,
		accessRoles: [EmployeeRole.ADMIN],
		items: [
			{
				name: 'Список',
				routeKey: 'ADMIN_INGREDIENTS',
				icon: Apple,
				accessRoles: [EmployeeRole.ADMIN],
			},
			{
				name: 'Категории',
				routeKey: 'ADMIN_INGREDIENT_CATEGORIES',
				icon: LayoutList,
				accessRoles: [EmployeeRole.ADMIN],
			},
		],
	},
	{
		label: 'Складские товары',
		icon: Package,
		accessRoles: [EmployeeRole.ADMIN],
		items: [
			{
				name: 'Список',
				routeKey: 'ADMIN_STOCK_MATERIALS',
				icon: Package,
				accessRoles: [EmployeeRole.ADMIN],
			},
			{
				name: 'Категории',
				routeKey: 'ADMIN_STOCK_MATERIAL_CATEGORIES',
				icon: LayoutList,
				accessRoles: [EmployeeRole.ADMIN],
			},
		],
	},
]
