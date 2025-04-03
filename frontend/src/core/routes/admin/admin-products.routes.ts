import type { AppRouteRecord } from '../routes.types'

export const ADMIN_PRODUCTS_CHILDREN_ROUTES = {
	ADMIN_PRODUCTS: {
		path: 'products',
		meta: {
			title: 'Все продукты',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/products/pages/admin-products-page.vue'),
	},
	ADMIN_PRODUCT_CREATE: {
		path: 'products/create',
		meta: {
			title: 'Создать продукт',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/products/pages/admin-product-create-page.vue'),
	},
	ADMIN_PRODUCT_DETAILS: {
		path: 'products/:id',
		meta: {
			title: 'Детали продукта',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/products/pages/admin-product-details-page.vue'),
	},
	ADMIN_PRODUCT_SIZE_DETAILS: {
		path: 'product-sizes/:id',
		meta: {
			title: 'Детали размера продукта',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/products/pages/admin-product-size-details-page.vue'),
	},
	ADMIN_PRODUCT_SIZE_CREATE: {
		path: 'product-sizes/create',
		meta: {
			title: 'Детали размера продукта',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/products/pages/admin-product-size-create-page.vue'),
	},
	ADMIN_STORE_PRODUCTS: {
		path: 'store-products',
		meta: {
			title: 'Продукты кафе',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/store-products/pages/admin-store-products-page.vue'),
	},
	ADMIN_STORE_PRODUCT_DETAILS: {
		path: 'store-products/:id',
		meta: {
			title: 'Детали продукта кафе',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/store-products/pages/admin-store-product-details-page.vue'),
	},
	ADMIN_STORE_PRODUCT_CREATE: {
		path: 'store-products/create',
		meta: {
			title: 'Добавить продукт в кафе',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/store-products/pages/admin-store-product-create-page.vue'),
	},
	ADMIN_PRODUCT_CATEGORIES: {
		path: 'product-categories',
		meta: {
			title: 'Категории продуктов',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/product-categories/pages/admin-product-categories-page.vue'),
	},
	ADMIN_PRODUCT_CATEGORY_DETAILS: {
		path: 'product-categories/:id',
		meta: {
			title: 'Детали категории продукта',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/product-categories/pages/admin-product-category-details-page.vue'),
	},
	ADMIN_PRODUCT_CATEGORY_CREATE: {
		path: 'product-categories/create',
		meta: {
			title: 'Cоздать категорию продукта',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/product-categories/pages/admin-product-category-create-page.vue'),
	},
} satisfies AppRouteRecord
