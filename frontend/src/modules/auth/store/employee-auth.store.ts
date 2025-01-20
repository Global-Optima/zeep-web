import type { EmployeeDTO } from '@/modules/admin/store-employees/models/employees.models'
import { defineStore } from 'pinia'
import { computed, ref, type Ref } from 'vue'

export const useEmployeeAuthStore = defineStore('EMPLOYEE_AUTH', () => {
	const currentEmployee: Ref<EmployeeDTO | null> = ref(null)
	const isLoggedIn = computed(() => Boolean(true))

	const setCurrentEmployee = (user: EmployeeDTO | null) => {
		currentEmployee.value = user
	}

	return { currentEmployee, isLoggedIn, setCurrentEmployee }
})
