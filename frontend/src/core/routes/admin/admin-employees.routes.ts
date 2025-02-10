import type { AppRouteRecord } from '../routes.types'

export const ADMIN_EMPLOYEES_CHILDREN_ROUTES = {
	ADMIN_STORE_EMPLOYEES: {
		path: 'employees/store',
		meta: {
			title: 'Все сотрудники кафе',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/employees/stores/pages/admin-store-employees-page.vue'),
	},
	ADMIN_STORE_EMPLOYEE_DETAILS: {
		path: 'employees/store/:id',
		meta: {
			title: 'Детали сотрудника кафе',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/employees/stores/pages/admin-store-employee-details-page.vue'),
	},
	ADMIN_STORE_EMPLOYEE_UPDATE: {
		path: 'employees/store/:id/update',
		meta: {
			title: 'Обновить сотрудника кафе',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/employees/stores/pages/admin-store-employee-update-page.vue'),
	},
	ADMIN_STORE_EMPLOYEE_CREATE: {
		path: 'stores/:storeId/employees/create',
		meta: {
			title: 'Создать сотрудника кафе',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/employees/stores/pages/admin-store-employee-create-page.vue'),
	},

	// Franchisee
	ADMIN_FRANCHISEE_EMPLOYEES: {
		path: 'employees/franchisee',
		meta: {
			title: 'Все сотрудники франчайзи',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/employees/franchisees/pages/admin-franchisee-employees-page.vue'),
	},
	ADMIN_FRANCHISEE_EMPLOYEE_DETAILS: {
		path: 'employees/franchisee/:id',
		meta: {
			title: 'Детали сотрудника франчайзи',
			requiresAuth: true,
		},
		component: () =>
			import(
				'@/modules/admin/employees/franchisees/pages/admin-franchisee-employee-details-page.vue'
			),
	},
	ADMIN_FRANCHISEE_EMPLOYEE_UPDATE: {
		path: 'employees/franchisee/:id/update',
		meta: {
			title: 'Обновить сотрудника франчайзи',
			requiresAuth: true,
		},
		component: () =>
			import(
				'@/modules/admin/employees/franchisees/pages/admin-franchisee-employee-update-page.vue'
			),
	},
	ADMIN_FRANCHISEE_EMPLOYEE_CREATE: {
		path: 'franchisees/:franchiseeId/employees/create',
		meta: {
			title: 'Создать сотрудника франчайзи',
			requiresAuth: true,
		},
		component: () =>
			import(
				'@/modules/admin/employees/franchisees/pages/admin-franchisee-employee-create-page.vue'
			),
	},

	// Regions
	ADMIN_REGION_EMPLOYEES: {
		path: 'employees/region',
		meta: {
			title: 'Все региональные сотрудники',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/employees/regions/pages/admin-region-employees-page.vue'),
	},
	ADMIN_REGION_EMPLOYEE_DETAILS: {
		path: 'employees/region/:id',
		meta: {
			title: 'Детали регионального сотрудника',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/employees/regions/pages/admin-region-employee-details-page.vue'),
	},
	ADMIN_REGION_EMPLOYEE_UPDATE: {
		path: 'employees/region/:id/update',
		meta: {
			title: 'Обновить регионального сотрудника',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/employees/regions/pages/admin-region-employee-update-page.vue'),
	},
	ADMIN_REGION_EMPLOYEE_CREATE: {
		path: 'regions/:regionId/employees/create',
		meta: {
			title: 'Создать сотрудника региона',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/employees/regions/pages/admin-region-employee-create-page.vue'),
	},

	// Warehouse
	ADMIN_WAREHOUSE_EMPLOYEES: {
		path: 'employees/warehouse',
		meta: {
			title: 'Все сотрудники складов',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/employees/warehouses/pages/admin-warehouse-employees-page.vue'),
	},
	ADMIN_WAREHOUSE_EMPLOYEE_DETAILS: {
		path: 'employees/warehouse/:id',
		meta: {
			title: 'Детали сотрудника склада',
			requiresAuth: true,
		},
		component: () =>
			import(
				'@/modules/admin/employees/warehouses/pages/admin-warehouse-employee-details-page.vue'
			),
	},
	ADMIN_WAREHOUSE_EMPLOYEE_UPDATE: {
		path: 'employees/warehouse/:id/update',
		meta: {
			title: 'Обновить сотрудника склада',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/employees/warehouses/pages/admin-warehouse-employee-update-page.vue'),
	},
	ADMIN_WAREHOUSE_EMPLOYEE_CREATE: {
		path: 'warehouses/:warehouseId/employees/create',
		meta: {
			title: 'Создать сотрудника склада',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/employees/warehouses/pages/admin-warehouse-employee-create-page.vue'),
	},

	ADMIN_EMPLOYEE_AUDIT: {
		path: 'employees/:id/audit',
		meta: {
			title: 'Аудит сотруника',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/employees/pages/admin-employees-audit-page.vue'),
	},
	ADMIN_EMPLOYEE_REASSIGN_TYPE: {
		path: 'employees/:id/reassign',
		meta: {
			title: 'Сменить тип сотрудник',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/employees/pages/admin-employee-reassign-type-page.vue'),
	},
} satisfies AppRouteRecord
