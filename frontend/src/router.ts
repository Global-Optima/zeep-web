import { createRouter, createWebHistory } from 'vue-router'

import { AppTranslation } from './core/config/locale.config'
import { getRouteName, ROUTES } from './core/config/routes.config'
import { DEFAULT_TITLE, TITLE_TEMPLATE } from './core/constants/seo.constants'
import { CURRENT_STORE_COOKIES_CONFIG } from './modules/stores/constants/store-cookies.constant'

export const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
	scrollBehavior() {
		return { top: 0 }
	},
	routes: [
		{
			path: '/:locale?',
			beforeEnter: AppTranslation.routeMiddleware,
			children: ROUTES,
		},
	],
})

router.beforeEach(async (to, _from, next) => {
	const metaTitle = to.meta?.title ? TITLE_TEMPLATE(to.meta.title as string) : DEFAULT_TITLE
	document.title = metaTitle

	if (to.name !== getRouteName('LOGIN')) {
		const storeId = localStorage.getItem(CURRENT_STORE_COOKIES_CONFIG.key)

		if (!storeId) {
			return next({ name: getRouteName('LOGIN') })
		}
	}

	return next()
})
