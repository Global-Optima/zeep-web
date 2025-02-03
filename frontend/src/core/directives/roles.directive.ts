import type { EmployeeRole } from '@/modules/admin/store-employees/models/employees.models'
import { useEmployeeAuthStore } from '@/modules/auth/store/employee-auth.store'
import type { Directive, DirectiveBinding } from 'vue'

export const roleDirective: Directive = {
	mounted(el, binding: DirectiveBinding<EmployeeRole | EmployeeRole[]>) {
		const authStore = useEmployeeAuthStore()
		const userRole = authStore.currentEmployee?.role

		if (!userRole) {
			el.remove()
			return
		}

		const allowedRoles = Array.isArray(binding.value) ? binding.value : [binding.value]

		if (!allowedRoles.includes(userRole)) {
			el.remove()
		}
	},
}
