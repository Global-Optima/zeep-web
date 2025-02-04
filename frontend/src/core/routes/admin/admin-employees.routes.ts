import type { AppRouteRecord } from '../routes.types'

export const ADMIN_EMPLOYEES_CHILDREN_ROUTES = {
	ADMIN_STORE_EMPLOYEES: {
		path: 'store-employees',
		meta: {
			title: 'Все сотрудники кафе',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/employees/stores/pages/admin-store-employees-page.vue'),
	},
	ADMIN_STORE_EMPLOYEE_DETAILS: {
		path: 'store-employees/:id',
		meta: {
			title: 'Детали сотрудника кафе',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/employees/stores/pages/admin-store-employee-details-page.vue'),
	},
	ADMIN_STORE_EMPLOYEE_UPDATE: {
		path: 'store-employees/:id/update',
		meta: {
			title: 'Обновить сотрудника кафе',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/employees/stores/pages/admin-store-employee-update-page.vue'),
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
