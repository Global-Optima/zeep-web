import { createRouter, createWebHistory } from 'vue-router'

import { isAxiosError } from 'axios'
import { getRouteName, ROUTES } from './core/config/routes.config'
import { DEFAULT_TITLE, TITLE_TEMPLATE } from './core/constants/seo.constants'
import { employeesService } from './modules/admin/store-employees/services/employees.service'
import { useEmployeeAuthStore } from './modules/auth/store/employee-auth.store'

export const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
	scrollBehavior() {
		return { top: 0 }
	},
	routes: ROUTES,
})

router.beforeEach(async (to, _from, next) => {
	const { setCurrentEmployee } = useEmployeeAuthStore()

	if (to.name === getRouteName('LOGIN')) {
		return next()
	}

	if (to.meta?.requiresAuth) {
		try {
			const currentEmployee = await employeesService.getCurrentEmployee()

			if (!currentEmployee) {
				return next({ name: getRouteName('LOGIN') })
			}

			setCurrentEmployee(currentEmployee)
		} catch (error) {
			console.error('Error fetching current employee:', error)

			if (isAxiosError(error) && error.status === 401) {
				return next({ name: getRouteName('LOGIN') })
			}

			return next({ name: getRouteName('INTERNAL_ERROR') })
		}
	}

	document.title = to.meta?.title ? TITLE_TEMPLATE(to.meta.title as string) : DEFAULT_TITLE

	return next()
})
