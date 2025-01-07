import { apiClient } from '@/core/config/axios-instance.config'
import type { CreateUnitDTO, UnitResponse, UpdateUnitDTO } from '../models/units.model'

class UnitsService {
	/**
	 * Fetch all units with optional filters
	 */
	async getAllUnits() {
		try {
			const response = await apiClient.get<UnitResponse[]>(`/units`)
			return response.data
		} catch (error) {
			console.error('Failed to fetch units: ', error)
			throw error
		}
	}

	/**
	 * Fetch unit details by ID
	 */
	async getUnitByID(id: number) {
		try {
			const response = await apiClient.get<UnitResponse>(`/units/${id}`)
			return response.data
		} catch (error) {
			console.error(`Failed to fetch unit details for ID ${id}: `, error)
			throw error
		}
	}

	/**
	 * Create a new unit
	 */
	async createUnit(data: CreateUnitDTO) {
		try {
			const response = await apiClient.post<void>(`/units`, data)
			return response.data
		} catch (error) {
			console.error('Failed to create unit: ', error)
			throw error
		}
	}

	/**
	 * Update an existing unit by ID
	 */
	async updateUnit(id: number, data: UpdateUnitDTO) {
		try {
			const response = await apiClient.put<void>(`/units/${id}`, data)
			return response.data
		} catch (error) {
			console.error(`Failed to update unit with ID ${id}: `, error)
			throw error
		}
	}

	/**
	 * Delete a unit by ID
	 */
	async deleteUnit(id: number) {
		try {
			await apiClient.delete<void>(`/units/${id}`)
		} catch (error) {
			console.error(`Failed to delete unit with ID ${id}: `, error)
			throw error
		}
	}
}

export const unitsService = new UnitsService()
