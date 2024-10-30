import { createRouter, createWebHistory } from 'vue-router'

import { DEFAULT_TITLE, TITLE_TEMPLATE } from './core/constants/seo.constants'
import { ROUTES } from './core/config/routes.config'

export const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
	scrollBehavior() {
		return { top: 0 }
	},
	routes: ROUTES,
})

router.beforeEach(async (to, _from, next) => {
	const metaTitle = to.meta?.title ? TITLE_TEMPLATE(to.meta.title as string) : DEFAULT_TITLE
	document.title = metaTitle

	return next()
})
