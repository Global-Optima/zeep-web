import { apiClient } from '@/core/config/axios-instance.config'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import {
	EmployeeType,
	type CreateEmployeeDto,
	type Employee,
	type EmployeesFilter,
	type UpdateEmployeeDto,
} from '../models/employees.models'

class EmployeesService {
	private readonly baseUrl = '/employees'

	async getStoreEmployees(storeID: number, filter?: EmployeesFilter) {
		try {
			return this.getEmployees({
				type: EmployeeType.Store,
				storeId: storeID,
				...buildRequestFilter(filter),
			})
		} catch (error) {
			console.error(`Failed to fetch employees for store ID ${storeID}:`, error)
			throw error
		}
	}

	async getEmployees(filter?: EmployeesFilter): Promise<Employee[]> {
		try {
			const response = await apiClient.get<Employee[]>(this.baseUrl, {
				params: buildRequestFilter(filter),
			})
			return response.data
		} catch (error) {
			console.error(`Failed to fetch employees:`, error)
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

	async getCurrentEmployee() {
		try {
			return apiClient.get<Employee>(`${this.baseUrl}/current`).then(res => res.data)
		} catch (error) {
			console.error(`Failed to get current employee:`, error)
			throw error
		}
	}

	async getEmployeeById(id: number) {
		try {
			return apiClient.get<Employee>(`${this.baseUrl}/${id}`).then(res => res.data)
		} catch (error) {
			console.error(`Failed to get employee by id:`, error)
			throw error
		}
	}

	async updateEmployee(id: number, dto: UpdateEmployeeDto) {
		try {
			const response = await apiClient.put<void>(`${this.baseUrl}/${id}`, dto)
			return response.data
		} catch (error) {
			console.error(`Failed to update employee by id ${id}:`, error)
			throw error
		}
	}

	async createEmployee(dto: CreateEmployeeDto) {
		try {
			const response = await apiClient.post<void>(`${this.baseUrl}`, dto)
			return response.data
		} catch (error) {
			console.error(`Failed to create employee:`, error)
			throw error
		}
	}
}

export const employeesService = new EmployeesService()
