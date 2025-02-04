import type { EmployeeRole } from '@/modules/admin/employees/models/employees.models'
import { useEmployeeAuthStore } from '@/modules/auth/store/employee-auth.store'
import { computed } from 'vue'

export function useHasRole(allowedRoles: EmployeeRole | EmployeeRole[]) {
	const authStore = useEmployeeAuthStore()
	const currentRole = computed(() => authStore.currentEmployee?.role)

	return computed(() => {
		if (!currentRole.value) return false

		const rolesArray = Array.isArray(allowedRoles) ? allowedRoles : [allowedRoles]

		return rolesArray.includes(currentRole.value)
	})
}
