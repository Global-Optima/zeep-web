import type { PaginationParams } from '@/core/utils/pagination.utils'

export enum NotificationPriority {
	HIGH = 'HIGH',
	MEDIUM = 'MEDIUM',
	LOW = 'LOW',
}

export enum NotificationEventType {
	STOCK_REQUEST_STATUS_UPDATED = 'STOCK_REQUEST_STATUS_UPDATED',
	NEW_ORDER = 'NEW_ORDER',
	STORE_WAREHOUSE_RUN_OUT = 'STORE_WAREHOUSE_RUN_OUT',
	CENTRAL_CATALOG_UPDATE = 'CENTRAL_CATALOG_UPDATE',
	STOCK_EXPIRATION = 'STOCK_EXPIRATION',
	OUT_OF_STOCK = 'OUT_OF_STOCK',
	NEW_STOCK_REQUEST = 'NEW_STOCK_REQUEST',
	PRICE_CHANGE = 'PRICE_CHANGE',
}

export interface NotificationDTO {
	id: number
	eventType: NotificationEventType
	priority: NotificationPriority
	messages: {
		en: string
		ru: string
		kk: string
	}
	details: Record<string, unknown>
	isRead: boolean
	createdAt: string
	updatedAt: string
}

export interface GetNotificationsFilter extends PaginationParams {
	priority?: NotificationPriority
	eventType?: NotificationEventType
	search?: string
	isRead?: boolean
	startDate?: string
	endDate?: string
}

export interface MarkNotificationAsReadDTO {
	notificationId: number
}

export interface MarkNotificationsAsReadDTO {
	notificationIds: number[]
}
