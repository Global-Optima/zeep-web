import { apiClient } from '@/core/config/axios-instance.config'
import { EmployeeType, type Employee, type EmployeeLoginDTO } from '../models/employees.models'

class EmployeesService {
	private readonly baseUrl = '/employees'

	async getStoreEmployees(storeID: number): Promise<Employee[]> {
		try {
			const response = await apiClient.get<Employee[]>(this.baseUrl, {
				params: { type: EmployeeType.Store, storeId: storeID },
			})
			return response.data
		} catch (error) {
			console.error(`Failed to fetch employees for store ID ${storeID}:`, error)
			throw error
		}
	}

	async getWarehouseEmployees(warehouseId: number): Promise<Employee[]> {
		try {
			const response = await apiClient.get<Employee[]>(this.baseUrl, {
				params: { type: EmployeeType.Warehouse, warehouseId: warehouseId },
			})
			return response.data
		} catch (error) {
			console.error(`Failed to fetch employees for warehouse ID ${warehouseId}:`, error)
			throw error
		}
	}

	async login(dto: EmployeeLoginDTO) {
		try {
			return apiClient.post(`${this.baseUrl}/login`, dto).then(res => res.data)
		} catch (error) {
			console.error(`Failed to login employee:`, error)
			throw error
		}
	}

	async logout() {
		try {
			return apiClient.post(`${this.baseUrl}/logout`).then(res => res.data)
		} catch (error) {
			console.error(`Failed to logout employee:`, error)
			throw error
		}
	}

	async getCurrentEmployee() {
		try {
			return apiClient.get<Employee>(`${this.baseUrl}/current`).then(res => res.data)
		} catch (error) {
			console.error(`Failed to get current employee:`, error)
			throw error
		}
	}
}

export const employeesService = new EmployeesService()
