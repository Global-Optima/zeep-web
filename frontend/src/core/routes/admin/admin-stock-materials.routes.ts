import type { AppRouteRecord } from '../../config/routes.config'

export const ADMIN_STOCK_MATERIALS_CHILDREN_ROUTES = {
	ADMIN_STOCK_MATERIALS: {
		path: 'stock-materials',
		meta: {
			title: 'Складские товары',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/stock-materials/pages/admin-stock-materials-page.vue'),
	},
	ADMIN_STOCK_MATERIAL_DETAILS: {
		path: 'stock-materials/:id',
		meta: {
			title: 'Детали складского товара',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/stock-materials/pages/admin-stock-material-details-page.vue'),
	},
	ADMIN_STOCK_MATERIAL_CREATE: {
		path: 'stock-materials/create',
		meta: {
			title: 'Создать складской товар',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/stock-materials/pages/admin-stock-material-create-page.vue'),
	},
	ADMIN_STOCK_MATERIAL_CATEGORIES: {
		path: 'stock-material-categories',
		meta: {
			title: 'Категории складских товаров',
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
			title: 'Детали категории складских товаров',
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
			title: 'Создать категорию складских товаров',
			requiresAuth: true,
		},
		component: () =>
			import(
				'@/modules/admin/stock-material-categories/pages/admin-stock-material-category-create-page.vue'
			),
	},
} satisfies AppRouteRecord
