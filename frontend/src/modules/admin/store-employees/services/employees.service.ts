import { apiClient } from '@/core/config/axios-instance.config'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type {
	CreateStoreEmployeeDTO,
	CreateWarehouseEmployeeDTO,
	EmployeeDTO,
	GetStoreEmployeesFilter,
	GetWarehouseEmployeesFilter,
	StoreEmployeeDTO,
	UpdateStoreEmployeeDTO,
	UpdateWarehouseEmployeeDTO,
	WarehouseEmployeeDTO,
} from '@/modules/admin/store-employees/models/employees.models'
import { type PaginatedResponse } from './../../../../core/utils/pagination.utils'

class EmployeeService {
	// Store Employees

	async getStoreEmployees(filter?: GetStoreEmployeesFilter) {
		const response = await apiClient.get<PaginatedResponse<StoreEmployeeDTO[]>>(
			`/store-employees`,
			{
				params: buildRequestFilter(filter),
			},
		)
		return response.data
	}

	async createStoreEmployee(dto: CreateStoreEmployeeDTO): Promise<StoreEmployeeDTO> {
		const response = await apiClient.post<StoreEmployeeDTO>(`/store-employees`, dto)
		return response.data
	}

	async getStoreEmployeeById(id: number): Promise<StoreEmployeeDTO> {
		const response = await apiClient.get<StoreEmployeeDTO>(`/store-employees/${id}`)
		return response.data
	}

	async updateStoreEmployee(id: number, dto: UpdateStoreEmployeeDTO): Promise<StoreEmployeeDTO> {
		const response = await apiClient.put<StoreEmployeeDTO>(`/store-employees/${id}`, dto)
		return response.data
	}

	// Warehouse Employees

	async getWarehouseEmployees(filter?: GetWarehouseEmployeesFilter) {
		const response = await apiClient.get<PaginatedResponse<WarehouseEmployeeDTO[]>>(
			`/warehouse-employees`,
			{
				params: buildRequestFilter(filter),
			},
		)
		return response.data
	}

	async createWarehouseEmployee(dto: CreateWarehouseEmployeeDTO): Promise<WarehouseEmployeeDTO> {
		const response = await apiClient.post<WarehouseEmployeeDTO>(`/warehouse-employees`, dto)
		return response.data
	}

	async getWarehouseEmployeeById(id: number): Promise<WarehouseEmployeeDTO> {
		const response = await apiClient.get<WarehouseEmployeeDTO>(`/warehouse-employees/${id}`)
		return response.data
	}

	async updateWarehouseEmployee(
		id: number,
		dto: UpdateWarehouseEmployeeDTO,
	): Promise<WarehouseEmployeeDTO> {
		const response = await apiClient.put<WarehouseEmployeeDTO>(`/warehouse-employees/${id}`, dto)
		return response.data
	}

	// General Employee Endpoints

	async getCurrentEmployee(): Promise<EmployeeDTO> {
		const response = await apiClient.get<EmployeeDTO>(`/employees/current`)
		return response.data
	}
}

export const employeesService = new EmployeeService()
