import type { AppRouteRecord } from '../routes.types'

export const ADMIN_SUPPLIERS_CHILDREN_ROUTES = {
	ADMIN_SUPPLIERS: {
		path: 'suppliers',
		meta: {
			title: 'Постащики',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/suppliers/pages/admin-suppliers-page.vue'),
	},
	ADMIN_SUPPLIER_CREATE: {
		path: 'suppliers/create',
		meta: {
			title: 'Создать постащика',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/suppliers/pages/admin-supplier-create-page.vue'),
	},
	ADMIN_SUPPLIER_DETAILS: {
		path: 'suppliers/:id',
		meta: {
			title: 'Детали постащика',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/suppliers/pages/admin-supplier-details-page.vue'),
	},
	ADMIN_WAREHOUSE_SUPPLIERS: {
		path: 'suppliers',
		meta: {
			title: 'Постащики',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/suppliers/pages/admin-suppliers-page.vue'),
	},
	ADMIN_WAREHOUSE_SUPPLIER_DETAILS: {
		path: 'suppliers/:id',
		meta: {
			title: 'Детали постащика',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/suppliers/pages/admin-supplier-details-page.vue'),
	},
} satisfies AppRouteRecord
