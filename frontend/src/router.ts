import { createRouter, createWebHistory } from 'vue-router'

import { ROUTES } from './core/config/routes.config'
import { DEFAULT_TITLE, TITLE_TEMPLATE } from './core/constants/seo.constants'
import { AppTranslation } from './core/config/locale.config'

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

	return next()
})
