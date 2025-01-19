import type { ButtonVariants } from '@/core/components/ui/button'
import { EmployeeRole } from '@/modules/admin/employees/models/employees.models'
import {
	StockRequestStatus,
	type AcceptWithChangeRequestStatusDTO,
	type RejectStockRequestStatusDTO,
} from '@/modules/admin/store-stock-requests/models/stock-requests.model'
import { stockRequestsService } from '@/modules/admin/store-stock-requests/services/stock-requests.service'

interface Action {
	label: string
	variant?: ButtonVariants['variant']
}

export interface DirectAction extends Action {
	type: 'DIRECT'
	handler: (requestId: number) => Promise<void>
}

export interface CommentDialogAction extends Action {
	type: 'COMMENT_DIALOG'
	handler: (requestId: number, payload: RejectStockRequestStatusDTO) => Promise<void>
}

export interface AcceptWithChangesDialogAction extends Action {
	type: 'ACCEPT_WITH_CHANGES_DIALOG'
	handler: (requestId: number, payload: AcceptWithChangeRequestStatusDTO) => Promise<void>
}

export type StockRequestAction = DirectAction | CommentDialogAction | AcceptWithChangesDialogAction

export function getActions(status: StockRequestStatus, role: EmployeeRole): StockRequestAction[] {
	if (role === EmployeeRole.BARISTA) {
		switch (status) {
			case StockRequestStatus.CREATED:
				return [
					{
						type: 'DIRECT',
						label: 'Удалить',
						variant: 'destructive', // highlight as destructive
						handler: id => stockRequestsService.deleteStockRequest(id),
					},
					{
						type: 'DIRECT',
						label: 'Отправить заявку на склад',
						variant: 'default',
						handler: id => stockRequestsService.setProcessedStatus(id),
					},
				]

			case StockRequestStatus.REJECTED_BY_WAREHOUSE:
				return [
					{
						type: 'DIRECT',
						label: 'Исправить и повторно отправить',
						variant: 'default',
						handler: id => stockRequestsService.setProcessedStatus(id),
					},
				]

			case StockRequestStatus.IN_DELIVERY:
				return [
					{
						type: 'ACCEPT_WITH_CHANGES_DIALOG',
						label: 'Принять с изменениями',
						variant: 'outline',
						handler: (id, dto) => stockRequestsService.acceptWithChanges(id, dto),
					},
					{
						type: 'COMMENT_DIALOG',
						label: 'Отклонить заявку',
						variant: 'destructive',
						handler: (id, dto) => stockRequestsService.rejectStore(id, dto),
					},
					{
						type: 'DIRECT',
						label: 'Завершить заявку',
						variant: 'default',
						handler: id => stockRequestsService.setCompletedStatus(id),
					},
				]

			default:
				return []
		}
	}

	if (role === EmployeeRole.WAREHOUSE) {
		switch (status) {
			case StockRequestStatus.PROCESSED:
				return [
					{
						type: 'COMMENT_DIALOG',
						label: 'Отклонить заявку',
						variant: 'destructive',
						handler: (id, dto) => stockRequestsService.rejectWarehouse(id, dto),
					},
					{
						type: 'DIRECT',
						label: 'Перевести заявку в доставку',
						variant: 'default',
						handler: id => stockRequestsService.setInDeliveryStatus(id),
					},
				]

			case StockRequestStatus.REJECTED_BY_STORE:
				return [
					{
						type: 'DIRECT',
						label: 'Исправить и повторно отправить',
						variant: 'default',
						handler: id => stockRequestsService.setInDeliveryStatus(id),
					},
				]

			default:
				return []
		}
	}

	return []
}
