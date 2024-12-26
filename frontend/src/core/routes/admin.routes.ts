import type { AppRouteRecord, ParentRoutePage } from '../config/routes.config'
import AppAdminLayout from '../layouts/admin/app-admin-layout.vue'

export const ADMIN_CHILDREN_ROUTES = {
	ADMIN_DASHBOARD: {
		path: '',
		meta: {
			title: 'Аналитика',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/dashboard/pages/admin-dashboard-page.vue'),
	},
	ADMIN_STORE_DASHBOARD: {
		path: 'store-analytics',
		meta: {
			title: 'Аналитика магазина',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/store-dashboard/pages/admin-store-dashboard-page.vue'),
	},
	ADMIN_STORE_ORDERS: {
		path: 'store-orders',
		meta: {
			title: 'Заказы магазина',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/store-orders/pages/admin-store-orders-page.vue'),
	},
	ADMIN_STORE_STOCKS: {
		path: 'store-stocks',
		meta: {
			title: 'Склад Магазина',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/store-stocks/pages/admin-store-stocks-page.vue'),
	},
	ADMIN_WAREHOUSE_STOCKS: {
		path: 'warehouse-stocks',
		meta: {
			title: 'Склад Магазина',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/warehouse-stocks/pages/admin-warehouse-stocks-page.vue'),
	},
	ADMIN_STORE_STOCKS_DETAILS: {
		path: 'store-stocks/:id',
		meta: {
			title: 'Детали запаса',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/store-stocks/pages/admin-store-stocks-details-page.vue'),
	},
	ADMIN_CREATE_STORE_STOCKS: {
		path: 'store-stocks/create',
		meta: {
			title: 'Добавить в склад магазина',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/store-stocks/pages/admin-store-stocks-create-page.vue'),
	},
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
		component: () => import('@/modules/admin/suppliers/pages/admin-suppliers-create-page.vue'),
	},
	ADMIN_SUPPLIERS_DETAILS: {
		path: 'suppliers/:id',
		meta: {
			title: 'Детали постащика',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/suppliers/pages/admin-suppliers-details-page.vue'),
	},
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
	ADMIN_PRODUCTS: {
		path: 'products',
		meta: {
			title: 'Все товары',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/products/pages/admin-products-page.vue'),
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
		component: () => import('@/modules/admin/products/pages/admin-product-sizes-details-page.vue'),
	},
	ADMIN_INGREDIENTS_DETAILS: {
		path: 'ingredients/:id',
		meta: {
			title: 'Детали ингредиента',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/ingredients/pages/admin-ingredients-details-page.vue'),
	},
	ADMIN_INGREDIENTS: {
		path: 'ingredients',
		meta: {
			title: 'Ингредиенты',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/ingredients/pages/admin-ingredients-page.vue'),
	},
	ADMIN_STORE_PRODUCTS: {
		path: 'store-products',
		meta: {
			title: 'Товары Магазина',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/store-products/pages/admin-store-products-page.vue'),
	},
	ADMIN_CREATE_EMPLOYEE: {
		path: 'employees/create',
		meta: {
			title: 'Добавить сотрудника',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/employees/pages/admin-store-employees-create-page.vue'),
	},
	ADMIN_EMPLOYEES: {
		path: 'employees',
		meta: {
			title: 'Сотрудники',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/employees/pages/admin-store-employees-page.vue'),
	},
	ADMIN_EMPLOYEES_DETAILS: {
		path: 'employees/:id',
		meta: {
			title: 'Сотрудник',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/employees/pages/admin-store-employees-details-page.vue'),
	},
	ADMIN_STORES: {
		path: 'stores',
		meta: {
			title: 'Магазины',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/stores/pages/admin-stores-page.vue'),
	},
	ADMIN_STORE_CREATE: {
		path: 'stores/create',
		meta: {
			title: 'Добавить магазин',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/stores/pages/admin-stores-create-page.vue'),
	},
	ADMIN_STORE_DETAILS: {
		path: 'stores/:id',
		meta: {
			title: 'Детали магазина',
			requiresAuth: true,
		},
		component: () => import('@/modules/admin/stores/pages/admin-store-details-page.vue'),
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
	ADMIN_ADDITIVE_CATEGORIES: {
		path: 'additive-categories',
		meta: {
			title: 'Категории топпингов',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/additive-categories/pages/admin-additive-categories-page.vue'),
	},
	ADMIN_STORE_STOCK_REQUESTS: {
		path: 'store-stock-requests',
		meta: {
			title: 'Запросы на склад',
			requiresAuth: true,
		},
		component: () =>
			import('@/modules/admin/store-stock-requests/pages/admin-store-stock-requests-list-page.vue'),
	},
	ADMIN_STORE_STOCK_REQUESTS_CREATE: {
		path: 'store-stock-requests/create',
		meta: {
			title: 'Создать запрос на склад',
			requiresAuth: true,
		},
		component: () =>
			import(
				'@/modules/admin/store-stock-requests/pages/admin-store-stock-requests-create-page.vue'
			),
	},
	ADMIN_STORE_STOCK_REQUESTS_UPDATE: {
		path: 'store-stock-requests/:id/update',
		meta: {
			title: 'Обновить запрос на склад',
			requiresAuth: true,
		},
		component: () =>
			import(
				'@/modules/admin/store-stock-requests/pages/admin-store-stock-requests-update-page.vue'
			),
	},
	ADMIN_STORE_STOCK_REQUESTS_DETAILS: {
		path: 'store-stock-requests/:id',
		meta: {
			title: 'Детали запроса на склад',
			requiresAuth: true,
		},
		component: () =>
			import(
				'@/modules/admin/store-stock-requests/pages/admin-store-stock-requests-details-page.vue'
			),
	},
	ADMIN_WAREHOUSE_STOCK_REQUESTS: {
		path: 'warehouse-stock-requests',
		meta: {
			title: 'Запросы на склад',
			requiresAuth: true,
		},
		component: () =>
			import(
				'@/modules/admin/warehouse-stock-requests/pages/admin-warehouse-stock-requests-list-page.vue'
			),
	},
	ADMIN_WAREHOUSE_STOCK_REQUESTS_DETAILS: {
		path: 'warehouse-stock-requests/:id',
		meta: {
			title: 'Детали запроса на склад',
			requiresAuth: true,
		},
		component: () =>
			import(
				'@/modules/admin/warehouse-stock-requests/pages/admin-warehouse-stock-requests-details-page.vue'
			),
	},
} satisfies AppRouteRecord

export const ADMIN_ROUTES_CONFIG = {
	path: '/admin',
	component: AppAdminLayout,
	children: ADMIN_CHILDREN_ROUTES,
} satisfies ParentRoutePage
