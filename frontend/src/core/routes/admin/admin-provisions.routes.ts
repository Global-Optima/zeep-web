import type { AppRouteRecord } from '../routes.types'

export const ADMIN_PROVISIONS_CHILDREN_ROUTES = {
	ADMIN_PROVISIONS: {
		path: 'provisions',
		meta: {
			title: 'Заготовки',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/provisions/pages/admin-provisions-page.vue'),
	},
	ADMIN_PROVISION_DETAILS: {
		path: 'provisions/:id',
		meta: {
			title: 'Детали заготовки',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/provisions/pages/admin-provision-details-page.vue'),
	},
	ADMIN_PROVISION_CREATE: {
		path: 'provisions/create',
		meta: {
			title: 'Создать заготовку',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/provisions/pages/admin-provision-create-page.vue'),
	},

	ADMIN_STORE_PROVISIONS: {
		path: 'store-provisions',
		meta: {
			title: 'Заготовки',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/store-provisions/pages/admin-store-provisions-page.vue'),
	},
	ADMIN_STORE_PROVISION_DETAILS: {
		path: 'store-provisions/:id',
		meta: {
			title: 'Детали заготовки',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/store-provisions/pages/admin-store-provision-details-page.vue'),
	},
	ADMIN_STORE_PROVISION_CREATE: {
		path: 'store-provisions/create',
		meta: {
			title: 'Создать заготовку',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/store-provisions/pages/admin-store-provision-create-page.vue'),
	},
} satisfies AppRouteRecord
