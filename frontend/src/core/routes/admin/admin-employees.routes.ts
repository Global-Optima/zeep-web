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

	ADMIN_EMPLOYEE_AUDIT: {
		path: 'employees/:id/audit',
		meta: {
			title: 'Аудит сотруника',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/employees/pages/admin-employees-audit-page.vue'),
	},
} satisfies AppRouteRecord
