import type { RouteKey } from '@/core/config/routes.config'
import { EmployeeRole } from '@/modules/admin/employees/models/employees.models'
import {
	Apple,
	Blocks,
	BookUser,
	Building2,
	LayoutList,
	ListPlus,
	MapPinned,
	Package,
	PackageCheck,
	Ruler,
	Settings,
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
	// {
	// 	name: 'Аналитика',
	// 	routeKey: 'ADMIN_DASHBOARD',
	// 	icon: ChartBar,
	// 	accessRoles: [
	// 		EmployeeRole.ADMIN,
	// 		EmployeeRole.STORE_MANAGER,
	// 		EmployeeRole.BARISTA,
	// 		EmployeeRole.WAREHOUSE_MANAGER,
	// 		EmployeeRole.WAREHOUSE_EMPLOYEE,
	// 	],
	// },
	{
		name: 'Администраторы',
		routeKey: 'ADMIN_ADMIN_EMPLOYEES',
		icon: Users,
		accessRoles: [EmployeeRole.ADMIN],
	},
	{
		name: 'Сотрудники',
		routeKey: 'ADMIN_FRANCHISEE_EMPLOYEES',
		icon: Users,
		accessRoles: [EmployeeRole.FRANCHISEE_OWNER, EmployeeRole.FRANCHISEE_MANAGER],
	},
	{
		name: 'Сотрудники',
		routeKey: 'ADMIN_REGION_EMPLOYEES',
		icon: Users,
		accessRoles: [EmployeeRole.OWNER],
	},
	{
		name: 'Сотрудники',
		routeKey: 'ADMIN_WAREHOUSE_EMPLOYEES',
		icon: Users,
		accessRoles: [EmployeeRole.WAREHOUSE_EMPLOYEE, EmployeeRole.WAREHOUSE_MANAGER],
	},
	{
		name: 'Сотрудники',
		routeKey: 'ADMIN_STORE_EMPLOYEES',
		icon: Users,
		accessRoles: [EmployeeRole.STORE_MANAGER],
	},
	{
		name: 'Заказы',
		routeKey: 'ADMIN_STORE_ORDERS',
		icon: ShoppingCart,
		accessRoles: [
			EmployeeRole.STORE_MANAGER,
			EmployeeRole.BARISTA,
			EmployeeRole.FRANCHISEE_MANAGER,
			EmployeeRole.FRANCHISEE_OWNER,
		],
	},
	{
		name: 'Склад',
		routeKey: 'ADMIN_STORE_STOCKS',
		icon: Warehouse,
		accessRoles: [
			EmployeeRole.STORE_MANAGER,
			EmployeeRole.BARISTA,
			EmployeeRole.FRANCHISEE_MANAGER,
			EmployeeRole.FRANCHISEE_OWNER,
		],
	},
	{
		name: 'Товары',
		routeKey: 'ADMIN_STORE_PRODUCTS',
		icon: Package,
		accessRoles: [
			EmployeeRole.STORE_MANAGER,
			EmployeeRole.BARISTA,
			EmployeeRole.FRANCHISEE_MANAGER,
			EmployeeRole.FRANCHISEE_OWNER,
		],
	},
	{
		name: 'Топпинги',
		routeKey: 'ADMIN_STORE_ADDITIVES',
		icon: ListPlus,
		accessRoles: [
			EmployeeRole.STORE_MANAGER,
			EmployeeRole.BARISTA,
			EmployeeRole.FRANCHISEE_MANAGER,
			EmployeeRole.FRANCHISEE_OWNER,
		],
	},
	{
		name: 'Кафе',
		routeKey: 'ADMIN_STORES',
		icon: Store,
		accessRoles: [
			EmployeeRole.ADMIN,
			EmployeeRole.FRANCHISEE_MANAGER,
			EmployeeRole.FRANCHISEE_OWNER,
		],
	},
	{
		name: 'Склады',
		routeKey: 'ADMIN_WAREHOUSES',
		icon: Warehouse,
		accessRoles: [EmployeeRole.ADMIN, EmployeeRole.REGION_WAREHOUSE_MANAGER],
	},
	{
		name: 'Регионы',
		routeKey: 'ADMIN_REGIONS',
		icon: MapPinned,
		accessRoles: [EmployeeRole.ADMIN],
	},
	{
		name: 'Франчайзи',
		routeKey: 'ADMIN_FRANCHISEES',
		icon: Building2,
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
		icon: BookUser,
		accessRoles: [EmployeeRole.ADMIN],
	},
	{
		name: 'Запросы на склад',
		routeKey: 'ADMIN_STORE_STOCK_REQUESTS',
		icon: Truck,
		accessRoles: [
			EmployeeRole.STORE_MANAGER,
			EmployeeRole.BARISTA,
			EmployeeRole.REGION_WAREHOUSE_MANAGER,
			EmployeeRole.FRANCHISEE_MANAGER,
			EmployeeRole.FRANCHISEE_OWNER,
		],
	},
	{
		name: 'Поставщики',
		routeKey: 'ADMIN_WAREHOUSE_SUPPLIERS',
		icon: BookUser,
		accessRoles: [
			EmployeeRole.WAREHOUSE_MANAGER,
			EmployeeRole.WAREHOUSE_EMPLOYEE,
			EmployeeRole.REGION_WAREHOUSE_MANAGER,
		],
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
		accessRoles: [
			EmployeeRole.WAREHOUSE_MANAGER,
			EmployeeRole.WAREHOUSE_EMPLOYEE,
			EmployeeRole.REGION_WAREHOUSE_MANAGER,
		],
	},
	{
		name: 'Доставки',
		routeKey: 'ADMIN_WAREHOUSE_DELIVERIES',
		icon: PackageCheck,
		accessRoles: [
			EmployeeRole.WAREHOUSE_MANAGER,
			EmployeeRole.WAREHOUSE_EMPLOYEE,
			EmployeeRole.REGION_WAREHOUSE_MANAGER,
		],
	},
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
		accessRoles: [
			EmployeeRole.ADMIN,
			EmployeeRole.WAREHOUSE_EMPLOYEE,
			EmployeeRole.WAREHOUSE_MANAGER,
			EmployeeRole.REGION_WAREHOUSE_MANAGER,
		],
		items: [
			{
				name: 'Список',
				routeKey: 'ADMIN_STOCK_MATERIALS',
				icon: Package,
				accessRoles: [
					EmployeeRole.ADMIN,
					EmployeeRole.WAREHOUSE_EMPLOYEE,
					EmployeeRole.WAREHOUSE_MANAGER,
				],
			},
			{
				name: 'Категории',
				routeKey: 'ADMIN_STOCK_MATERIAL_CATEGORIES',
				icon: LayoutList,
				accessRoles: [
					EmployeeRole.ADMIN,
					EmployeeRole.WAREHOUSE_EMPLOYEE,
					EmployeeRole.WAREHOUSE_MANAGER,
				],
			},
		],
	},
	{
		name: 'Настройка',
		routeKey: 'ADMIN_STORE_SETTINGS',
		icon: Settings,
		accessRoles: [EmployeeRole.STORE_MANAGER, EmployeeRole.BARISTA],
	},
]
