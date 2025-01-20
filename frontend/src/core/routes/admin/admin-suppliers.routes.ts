import type { AppRouteRecord } from '../../config/routes.config'

export const ADMIN_SUPPLIERS_CHILDREN_ROUTES = {
	ADMIN_SUPPLIERS: {
		path: 'suppliers',
		meta: {
			title: 'Постащики',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/suppliers/pages/admin-suppliers-page.vue'),
	},
	ADMIN_CREATE_SUPPLIER: {
		path: 'suppliers/create',
		meta: {
			title: 'Создать постащика',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/suppliers/pages/admin-supplier-create-page.vue'),
	},
} satisfies AppRouteRecord
