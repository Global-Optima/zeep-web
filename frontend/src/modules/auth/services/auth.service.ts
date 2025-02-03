import { apiClient } from '@/core/config/axios-instance.config'
import type {
	EmployeeAccount,
	EmployeeLoginDTO,
} from '@/modules/admin/store-employees/models/employees.models'

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

	async getStoreAccounts(storeId: number) {
		try {
			return apiClient
				.get<EmployeeAccount[]>(`${this.baseUrl}/employees/store/${storeId}`)
				.then(res => res.data)
		} catch (error) {
			console.error(`Failed to get store accounts:`, error)
			throw error
		}
	}

	async getWarehouseAccounts(warehouseId: number) {
		try {
			return apiClient
				.get<EmployeeAccount[]>(`${this.baseUrl}/employees/warehouse/${warehouseId}`)
				.then(res => res.data)
		} catch (error) {
			console.error(`Failed to get warehouse accounts:`, error)
			throw error
		}
	}

	async getAdminsAccounts() {
		try {
			return apiClient
				.get<EmployeeAccount[]>(`${this.baseUrl}/employees/admins`)
				.then(res => res.data)
		} catch (error) {
			console.error(`Failed to get admin accounts:`, error)
			throw error
		}
	}

	async getRegionAccounts(regionId: number) {
		try {
			return apiClient
				.get<EmployeeAccount[]>(`${this.baseUrl}/employees/region/${regionId}`)
				.then(res => res.data)
		} catch (error) {
			console.error(`Failed to get region accounts:`, error)
			throw error
		}
	}

	async getFranchiseeAccounts(franchiseeId: number) {
		try {
			return apiClient
				.get<EmployeeAccount[]>(`${this.baseUrl}/employees/franchisee/${franchiseeId}`)
				.then(res => res.data)
		} catch (error) {
			console.error(`Failed to get franchisee accounts:`, error)
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
