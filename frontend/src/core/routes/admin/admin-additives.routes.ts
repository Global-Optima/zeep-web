import type { AppRouteRecord } from '../routes.types'

export const ADMIN_ADDITIVES_CHILDREN_ROUTES = {
	ADMIN_ADDITIVES: {
		path: 'additives',
		meta: {
			title: 'Все топпинги',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/additives/pages/admin-additives-page.vue'),
	},
	ADMIN_ADDITIVE_DETAILS: {
		path: 'additives/:id',
		meta: {
			title: 'Детали топпинга',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/additives/pages/admin-additive-details-page.vue'),
	},
	ADMIN_ADDITIVE_CREATE: {
		path: 'additives/create',
		meta: {
			title: 'Создать топпинг',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/additives/pages/admin-additive-create-page.vue'),
	},
	ADMIN_ADDITIVE_CATEGORIES: {
		path: 'additive-categories',
		meta: {
			title: 'Категории топпингов',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/additive-categories/pages/admin-additive-categories-page.vue'),
	},
	ADMIN_ADDITIVE_CATEGORY_CREATE: {
		path: 'additive-categories/create',
		meta: {
			title: 'Создать категорию топпинга',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/additive-categories/pages/admin-additive-category-create-page.vue'),
	},
	ADMIN_ADDITIVE_CATEGORY_DETAILS: {
		path: 'additive-categories/:id',
		meta: {
			title: 'Категория топпинга',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/additive-categories/pages/admin-additive-category-details-page.vue'),
	},
	ADMIN_STORE_ADDITIVES: {
		path: 'store-additives',
		meta: {
			title: 'Топпинги Магазина',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/store-additives/pages/admin-store-additives-page.vue'),
	},
	ADMIN_STORE_ADDITIVE_DETAILS: {
		path: 'store-additives/:id',
		meta: {
			title: 'Детали топпинга магазина',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/store-additives/pages/admin-store-additive-details-page.vue'),
	},
	ADMIN_STORE_ADDITIVE_CREATE: {
		path: 'store-additives/create',
		meta: {
			title: 'Добавить топпинг в магазин',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/store-additives/pages/admin-store-additive-create-page.vue'),
	},
} satisfies AppRouteRecord
