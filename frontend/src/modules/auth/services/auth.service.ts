import { apiClient } from '@/core/config/axios-instance.config'
import type { EmployeeLoginDTO } from '@/modules/employees/models/employees.models'

class AuthService {
	private readonly baseUrl = '/auth'

	async loginEmployee(dto: EmployeeLoginDTO) {
		try {
			return apiClient.post(`${this.baseUrl}/employees/login`, dto).then(res => res.data)
		} catch (error) {
			console.error(`Failed to login employee:`, error)
			throw error
		}
	}

	async logoutEmployee() {
		try {
			return apiClient.post(`${this.baseUrl}/employees/logout`).then(res => res.data)
		} catch (error) {
			console.error(`Failed to logout employee:`, error)
			throw error
		}
	}
}

export const authService = new AuthService()
