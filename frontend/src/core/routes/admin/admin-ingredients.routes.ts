import type { AppRouteRecord } from '../../config/routes.config'

export const ADMIN_INGREDIENTS_CHILDREN_ROUTES = {
	ADMIN_INGREDIENTS_DETAILS: {
		path: 'ingredients/:id',
		meta: {
			title: 'Детали ингредиента',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/ingredients/pages/admin-ingredient-details-page.vue'),
	},
	ADMIN_INGREDIENT_CREATE: {
		path: 'ingredients/create',
		meta: {
			title: 'Создать ингредиент',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/ingredients/pages/admin-ingredient-create-page.vue'),
	},
	ADMIN_INGREDIENTS: {
		path: 'ingredients',
		meta: {
			title: 'Ингредиенты',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/ingredients/pages/admin-ingredients-page.vue'),
	},
} satisfies AppRouteRecord
