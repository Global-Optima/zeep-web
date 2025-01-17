import { createRouter, createWebHistory } from 'vue-router'

import { getRouteName, ROUTES } from './core/config/routes.config'
import { DEFAULT_TITLE, TITLE_TEMPLATE } from './core/constants/seo.constants'
import { useEmployeeAuthStore } from './modules/auth/store/employee-auth.store'
import { employeesService } from './modules/employees/services/employees.service'

export const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
	scrollBehavior() {
		return { top: 0 }
	},
	routes: ROUTES,
})

router.beforeEach(async (to, _from, next) => {
	const { setCurrentEmployee } = useEmployeeAuthStore()

	// Check for login page access
	if (to.name === getRouteName('LOGIN')) {
		return next()
	}

	// If route requires authentication
	if (to.meta?.requiresAuth) {
		try {
			// Check if the employee is already set in the store
			const employeeFromStore = useEmployeeAuthStore().currentEmployee

			if (!employeeFromStore) {
				// Fetch current employee if not already set
				const currentEmployee = await employeesService.getCurrentEmployee()

				if (!currentEmployee) {
					// Redirect to login if no employee is returned
					return next({ name: getRouteName('LOGIN') })
				}

				// Set current employee in the store
				setCurrentEmployee(currentEmployee)
			}
		} catch (error) {
			// Handle errors and redirect to an appropriate page
			console.error('Error fetching current employee:', error)
			return next({ name: getRouteName('INTERNAL_ERROR') })
		}
	}

	// Set page title
	document.title = to.meta?.title ? TITLE_TEMPLATE(to.meta.title as string) : DEFAULT_TITLE

	// Allow navigation
	return next()
})
