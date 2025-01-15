import { apiClient } from '@/core/config/axios-instance.config'
import type { PaginatedResponse } from '@/core/utils/pagination.utils'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import {
	type CreateEmployeeDto,
	type Employee,
	type StoreEmployee,
	type StoreEmployeesFilter,
	type UpdateEmployeeDto,
	type WarehouseEmployee,
	type WarehouseEmployeesFilter,
} from '../models/employees.models'

class EmployeesService {
	private readonly baseUrl = '/employees'

	async getStoreEmployees(filter?: StoreEmployeesFilter) {
		try {
			const response = await apiClient.get<PaginatedResponse<StoreEmployee[]>>(
				`${this.baseUrl}/store`,
				{
					params: buildRequestFilter(filter),
				},
			)
			return response.data
		} catch (error) {
			console.error(`Failed to fetch store employees: `, error)
			throw error
		}
	}

	async getWarehouseEmployees(filter?: WarehouseEmployeesFilter) {
		try {
			const response = await apiClient.get<PaginatedResponse<WarehouseEmployee[]>>(
				`${this.baseUrl}/warehouse`,
				{
					params: buildRequestFilter(filter),
				},
			)
			return response.data
		} catch (error) {
			console.error(`Failed to fetch warehouse employees: `, error)
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
