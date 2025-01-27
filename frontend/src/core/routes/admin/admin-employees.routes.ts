import type { AppRouteRecord } from '../../config/routes.config'

export const ADMIN_EMPLOYEES_CHILDREN_ROUTES = {
	ADMIN_STORE_EMPLOYEES: {
		path: 'store-employees',
		meta: {
			title: 'Все сотрудники кафе',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/store-employees/pages/admin-store-employees-page.vue'),
	},
	ADMIN_STORE_EMPLOYEE_DETAILS: {
		path: 'store-employees/:id',
		meta: {
			title: 'Детали сотрудника кафе',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/store-employees/pages/admin-store-employee-details-page.vue'),
	},
	ADMIN_STORE_EMPLOYEE_AUDIT: {
		path: 'store-employees/:id/audit',
		meta: {
			title: 'Аудит сотруника',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/store-employees/pages/admin-store-employee-audit-page.vue'),
	},
	ADMIN_STORE_EMPLOYEE_UPDATE: {
		path: 'store-employees/:id/update',
		meta: {
			title: 'Обновить сотрудника кафе',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/store-employees/pages/admin-store-employee-update-page.vue'),
	},
	ADMIN_STORE_EMPLOYEE_CREATE: {
		path: 'store-employees/create',
		meta: {
			title: 'Создать сотрудника кафе',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/store-employees/pages/admin-store-employee-create-page.vue'),
	},
} satisfies AppRouteRecord
