import type { AppRouteRecord } from '../routes.types'

export const ADMIN_INGREDIENTS_CHILDREN_ROUTES = {
	ADMIN_INGREDIENTS_DETAILS: {
		path: 'ingredients/:id',
		meta: {
			title: 'Детали сырья',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/ingredients/pages/admin-ingredient-details-page.vue'),
	},
	ADMIN_INGREDIENT_CREATE: {
		path: 'ingredients/create',
		meta: {
			title: 'Создать сырье',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/ingredients/pages/admin-ingredient-create-page.vue'),
	},
	ADMIN_INGREDIENTS: {
		path: 'ingredients',
		meta: {
			title: 'Сырье',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/ingredients/pages/admin-ingredients-page.vue'),
	},
	ADMIN_INGREDIENT_CATEGORIES: {
		path: 'ingredient-categories',
		meta: {
			title: 'Категории продуктов',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/ingredient-categories/pages/admin-ingredient-categories-page.vue'),
	},
	ADMIN_INGREDIENT_CATEGORY_DETAILS: {
		path: 'ingredient-categories/:id',
		meta: {
			title: 'Детали категории продукта',
			requiresAuth: true,
		},
		component: () =>
			import(
				'@/modules/admin/ingredient-categories/pages/admin-ingredient-category-details-page.vue'
			),
	},
	ADMIN_INGREDIENT_CATEGORY_CREATE: {
		path: 'ingredient-categories/create',
		meta: {
			title: 'Cоздать категорию продукта',
			requiresAuth: true,
		},
		component: () =>
			import(
				'@/modules/admin/ingredient-categories/pages/admin-ingredient-category-create-page.vue'
			),
	},
} satisfies AppRouteRecord
