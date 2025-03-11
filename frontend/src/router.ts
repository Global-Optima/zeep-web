import { isAxiosError } from 'axios'
import { createRouter, createWebHistory } from 'vue-router'
import { getRouteName, getRoutes } from './core/config/routes.config'
import { DEFAULT_TITLE, TITLE_TEMPLATE } from './core/constants/seo.constants'
import { employeesService } from './modules/admin/employees/services/employees.service'
import { useEmployeeAuthStore } from './modules/auth/store/employee-auth.store'

export const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
	scrollBehavior() {
		return { top: 0, behavior: 'smooth' }
	},
	routes: getRoutes(),
})

router.beforeEach(async (to, from, next) => {
	const { setCurrentEmployee } = useEmployeeAuthStore()

	// Prevent navigating outside /kiosk
	const isLeavingKiosk = from.path.startsWith('/kiosk') && !to.path.startsWith('/kiosk')
	if (isLeavingKiosk) {
		return next(false) // Block navigation outside kiosk
	}

	// Allow login route without authentication
	if (to.name === getRouteName('LOGIN')) {
		return next()
	}

	// Authentication check for protected routes
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

	// Update document title
	document.title = to.meta?.title ? TITLE_TEMPLATE(to.meta.title as string) : DEFAULT_TITLE

	return next()
})
