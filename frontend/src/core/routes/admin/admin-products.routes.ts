import type { AppRouteRecord } from '../../config/routes.config'

export const ADMIN_PRODUCTS_CHILDREN_ROUTES = {
	ADMIN_PRODUCTS: {
		path: 'products',
		meta: {
			title: 'Все товары',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/products/pages/admin-products-page.vue'),
	},
	ADMIN_PRODUCT_CREATE: {
		path: 'products/create',
		meta: {
			title: 'Создать товар',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/products/pages/admin-product-create-page.vue'),
	},
	ADMIN_PRODUCT_DETAILS: {
		path: 'products/:id',
		meta: {
			title: 'Детали товара',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/products/pages/admin-product-details-page.vue'),
	},
	ADMIN_PRODUCT_SIZE_DETAILS: {
		path: 'product-sizes/:id',
		meta: {
			title: 'Детали размера товара',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/products/pages/admin-product-size-details-page.vue'),
	},
	ADMIN_PRODUCT_SIZE_CREATE: {
		path: 'product-sizes/create',
		meta: {
			title: 'Детали размера товара',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/products/pages/admin-product-size-create-page.vue'),
	},
	ADMIN_STORE_PRODUCTS: {
		path: 'store-products',
		meta: {
			title: 'Товары Магазина',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/store-products/pages/admin-store-products-page.vue'),
	},
	ADMIN_PRODUCT_CATEGORIES: {
		path: 'product-categories',
		meta: {
			title: 'Категории товаров',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/product-categories/pages/admin-product-categories-page.vue'),
	},
} satisfies AppRouteRecord
