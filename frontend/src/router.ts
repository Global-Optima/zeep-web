import { isAxiosError } from 'axios'
import { createRouter, createWebHistory, type RouteLocation } from 'vue-router'
import { getRoutes } from './core/config/routes.config'
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
	const isLeavingKiosk =
		from.path.startsWith('/kiosk') &&
		!to.path.startsWith('/kiosk') &&
		!from.path.startsWith('/kiosk/checklist')
	if (isLeavingKiosk) {
		return next(false) // Block navigation outside kiosk
	}

	// Allow login route without authentication
	if (to.name === 'LOGIN') {
		return next()
	}

	// Authentication check for protected routes
	if (to.meta?.requiresAuth) {
		try {
			const employee = await employeesService.getCurrentEmployee()
			if (!employee) {
				return next({ name: 'LOGIN' })
			}

			setCurrentEmployee(employee)
		} catch (err: unknown) {
			console.error('Unexpected error fetching employee:', err)

			if (isAxiosError(err)) {
				const status = err.response?.status
				if (status === 401) {
					return next({ name: 'LOGIN' })
				}
			}

			return next({ name: 'INTERNAL_ERROR' })
		}
	}

	document.title = to.meta?.title ? TITLE_TEMPLATE(to.meta.title as string) : DEFAULT_TITLE

	return next()
})

router.onError((error, to: RouteLocation) => {
	if (
		error.message.includes('Failed to fetch dynamically imported module') ||
		error.message.includes('Importing a module script failed')
	) {
		window.location.href = to.fullPath
	}

	if (error.response?.status === 401) {
		router.push({ name: 'LOGIN' })
	}

	if (error.response?.status === 404) {
		router.push({ name: 'INTERNAL_ERROR' })
	}

	console.error('Routing error:', error)
})
