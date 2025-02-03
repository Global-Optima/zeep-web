import { apiClient } from '@/core/config/axios-instance.config'
import { type PaginatedResponse } from '@/core/utils/pagination.utils'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type {
	GetNotificationsFilter,
	MarkNotificationAsReadDTO,
	MarkNotificationsAsReadDTO,
	NotificationDTO,
} from '../models/notifications.model'

class NotificationsService {
	async getNotifications(filter?: GetNotificationsFilter) {
		try {
			const response = await apiClient.get<PaginatedResponse<NotificationDTO[]>>('/notifications', {
				params: buildRequestFilter(filter),
			})
			return response.data
		} catch (error) {
			console.error('Failed to fetch notifications:', error)
			throw error
		}
	}

	async getNotification(id: number): Promise<NotificationDTO> {
		try {
			const response = await apiClient.get<NotificationDTO>(`/notifications/${id}`)
			return response.data
		} catch (error) {
			console.error(`Failed to fetch notification by ID ${id}:`, error)
			throw error
		}
	}

	async markNotificationAsRead(dto: MarkNotificationAsReadDTO) {
		try {
			const response = await apiClient.post<void>(
				`/notifications/${dto.notificationId}/mark-as-read`,
			)
			return response.data
		} catch (error) {
			console.error(`Failed to mark notification ${dto.notificationId} as read:`, error)
			throw error
		}
	}

	async markNotificationsAsRead(dto: MarkNotificationsAsReadDTO) {
		try {
			const response = await apiClient.post<void>('/notifications/mark-multiple-as-read', dto)
			return response.data
		} catch (error) {
			console.error('Failed to mark multiple notifications as read:', error)
			throw error
		}
	}

	async deleteNotification(id: number) {
		try {
			const response = await apiClient.delete<void>(`/notifications/${id}`)
			return response.data
		} catch (error) {
			console.error(`Failed to delete notification ${id}:`, error)
			throw error
		}
	}
}

export const notificationsService = new NotificationsService()
