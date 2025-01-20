import { apiClient } from '@/core/config/axios-instance.config'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type {
	CreateStoreEmployeeDTO,
	CreateWarehouseEmployeeDTO,
	EmployeeDTO,
	GetStoreEmployeesFilter,
	GetWarehouseEmployeesFilter,
	RoleDTO,
	StoreEmployeeDTO,
	UpdatePasswordDTO,
	UpdateStoreEmployeeDTO,
	UpdateWarehouseEmployeeDTO,
	WarehouseEmployeeDTO,
} from '@/modules/admin/store-employees/models/employees.models'
import { type PaginatedResponse } from './../../../../core/utils/pagination.utils'

class EmployeeService {
	private readonly baseUrl: string = '/employees'

	// Store Employees

	async getStoreEmployees(filter?: GetStoreEmployeesFilter) {
		const response = await apiClient.get<PaginatedResponse<StoreEmployeeDTO[]>>(
			`${this.baseUrl}/stores`,
			{
				params: buildRequestFilter(filter),
			},
		)
		return response.data
	}

	async createStoreEmployee(dto: CreateStoreEmployeeDTO): Promise<StoreEmployeeDTO> {
		const response = await apiClient.post<StoreEmployeeDTO>(`${this.baseUrl}/stores`, dto)
		return response.data
	}

	async getStoreEmployeeById(id: number): Promise<StoreEmployeeDTO> {
		const response = await apiClient.get<StoreEmployeeDTO>(`${this.baseUrl}/stores/${id}`)
		return response.data
	}

	async updateStoreEmployee(id: number, dto: UpdateStoreEmployeeDTO): Promise<StoreEmployeeDTO> {
		const response = await apiClient.put<StoreEmployeeDTO>(`${this.baseUrl}/stores/${id}`, dto)
		return response.data
	}

	// Warehouse Employees

	async getWarehouseEmployees(filter?: GetWarehouseEmployeesFilter) {
		const response = await apiClient.get<PaginatedResponse<WarehouseEmployeeDTO[]>>(
			`${this.baseUrl}/warehouses`,
			{
				params: buildRequestFilter(filter),
			},
		)
		return response.data
	}

	async createWarehouseEmployee(dto: CreateWarehouseEmployeeDTO): Promise<WarehouseEmployeeDTO> {
		const response = await apiClient.post<WarehouseEmployeeDTO>(`${this.baseUrl}/warehouses`, dto)
		return response.data
	}

	async getWarehouseEmployeeById(id: number): Promise<WarehouseEmployeeDTO> {
		const response = await apiClient.get<WarehouseEmployeeDTO>(`${this.baseUrl}/warehouses/${id}`)
		return response.data
	}

	async updateWarehouseEmployee(
		id: number,
		dto: UpdateWarehouseEmployeeDTO,
	): Promise<WarehouseEmployeeDTO> {
		const response = await apiClient.put<WarehouseEmployeeDTO>(
			`${this.baseUrl}/warehouses/${id}`,
			dto,
		)
		return response.data
	}

	// General Employee Endpoints

	async getCurrentEmployee(): Promise<EmployeeDTO> {
		const response = await apiClient.get<EmployeeDTO>(`${this.baseUrl}/current`)
		return response.data
	}

	async deleteEmployee(id: number): Promise<void> {
		await apiClient.delete(`${this.baseUrl}/${id}`)
	}

	async getAllRoles(): Promise<RoleDTO[]> {
		const response = await apiClient.get<RoleDTO[]>(`${this.baseUrl}/roles`)
		return response.data
	}

	async updatePassword(id: number, dto: UpdatePasswordDTO): Promise<void> {
		await apiClient.put(`${this.baseUrl}/${id}/password`, dto)
	}
}

export const employeesService = new EmployeeService()
