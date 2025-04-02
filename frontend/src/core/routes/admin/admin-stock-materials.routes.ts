import type { AppRouteRecord } from '../routes.types'

export const ADMIN_STOCK_MATERIALS_CHILDREN_ROUTES = {
	ADMIN_STOCK_MATERIALS: {
		path: 'stock-materials',
		meta: {
			title: 'Складские продукты',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/stock-materials/pages/admin-stock-materials-page.vue'),
	},
	ADMIN_STOCK_MATERIAL_DETAILS: {
		path: 'stock-materials/:id',
		meta: {
			title: 'Детали складского продукта',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/stock-materials/pages/admin-stock-material-details-page.vue'),
	},
	ADMIN_STOCK_MATERIAL_CREATE: {
		path: 'stock-materials/create',
		meta: {
			title: 'Создать складской продукт',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/stock-materials/pages/admin-stock-material-create-page.vue'),
	},
	ADMIN_STOCK_MATERIAL_CATEGORIES: {
		path: 'stock-material-categories',
		meta: {
			title: 'Категории складских продуктов',
			requiresAuth: true,
		},
		component: () =>
			import(
				'@/modules/admin/stock-material-categories/pages/admin-stock-material-categories-page.vue'
			),
	},
	ADMIN_STOCK_MATERIAL_CATEGORY_DETAILS: {
		path: 'stock-material-categories/:id',
		meta: {
			title: 'Детали категории складских продуктов',
			requiresAuth: true,
		},
		component: () =>
			import(
				'@/modules/admin/stock-material-categories/pages/admin-stock-material-category-details-page.vue'
			),
	},
	ADMIN_STOCK_MATERIAL_CATEGORY_CREATE: {
		path: 'stock-material-categories/create',
		meta: {
			title: 'Создать категорию складских продуктов',
			requiresAuth: true,
		},
		component: () =>
			import(
				'@/modules/admin/stock-material-categories/pages/admin-stock-material-category-create-page.vue'
			),
	},
} satisfies AppRouteRecord
