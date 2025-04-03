import type { AppRouteRecord } from '../routes.types'

export const ADMIN_ADDITIVES_CHILDREN_ROUTES = {
	ADMIN_ADDITIVES: {
		path: 'additives',
		meta: {
			title: 'Все модификаторы',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/additives/pages/admin-additives-page.vue'),
	},
	ADMIN_ADDITIVE_DETAILS: {
		path: 'additives/:id',
		meta: {
			title: 'Детали модификатора',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/additives/pages/admin-additive-details-page.vue'),
	},
	ADMIN_ADDITIVE_CREATE: {
		path: 'additives/create',
		meta: {
			title: 'Создать модификатор',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/additives/pages/admin-additive-create-page.vue'),
	},
	ADMIN_ADDITIVE_CATEGORIES: {
		path: 'additive-categories',
		meta: {
			title: 'Категории модификаторов',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/additive-categories/pages/admin-additive-categories-page.vue'),
	},
	ADMIN_ADDITIVE_CATEGORY_CREATE: {
		path: 'additive-categories/create',
		meta: {
			title: 'Создать категорию модификатора',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/additive-categories/pages/admin-additive-category-create-page.vue'),
	},
	ADMIN_ADDITIVE_CATEGORY_DETAILS: {
		path: 'additive-categories/:id',
		meta: {
			title: 'Категория модификатора',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/additive-categories/pages/admin-additive-category-details-page.vue'),
	},
	ADMIN_STORE_ADDITIVES: {
		path: 'store-additives',
		meta: {
			title: 'Модификаторы кафе',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/store-additives/pages/admin-store-additives-page.vue'),
	},
	ADMIN_STORE_ADDITIVE_DETAILS: {
		path: 'store-additives/:id',
		meta: {
			title: 'Детали модификатора кафе',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/store-additives/pages/admin-store-additive-details-page.vue'),
	},
	ADMIN_STORE_ADDITIVE_CREATE: {
		path: 'store-additives/create',
		meta: {
			title: 'Добавить модификатор в кафе',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/store-additives/pages/admin-store-additive-create-page.vue'),
	},
} satisfies AppRouteRecord
