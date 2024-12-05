import type { Employee } from '@/modules/employees/models/employees.models'
import { defineStore } from 'pinia'
import { computed, ref, type Ref } from 'vue'

export const useEmployeeAuthStore = defineStore('EMPLOYEE_AUTH', () => {
	const currentEmployee: Ref<Employee | null> = ref(null)
	const isLoggedIn = computed(() => Boolean(true))

	const setCurrentEmployee = (user: Employee | null) => {
		currentEmployee.value = user
	}

	return { currentEmployee, isLoggedIn, setCurrentEmployee }
})
