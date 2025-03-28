<template>
	<div class="flex flex-col gap-2">
		<Button
			v-for="action in actions"
			:key="action.label"
			:variant="action.variant ?? 'default'"
			:disabled="isPending"
			class="w-full"
			@click="onActionClick(action)"
		>
			{{ action.label }}
		</Button>

		<AdminStockRequestsCommentDialog
			v-if="commentDialogOpen"
			:open="commentDialogOpen"
			@close="closeCommentDialog"
			@submit="handleCommentSubmit"
		/>

		<AdminStockRequestsChangesDialog
			v-if="acceptChangesDialogOpen"
			:open="acceptChangesDialogOpen"
			:initial-items="request.stockMaterials"
			@close="closeAcceptDialog"
			@submit="handleAcceptWithChangesSubmit"
		/>
	</div>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

import AdminStockRequestsChangesDialog from '@/modules/admin/stock-requests/components/details/admin-stock-requests-changes-dialog.vue'
import AdminStockRequestsCommentDialog from '@/modules/admin/stock-requests/components/details/admin-stock-requests-comment-dialog.vue'

import {
  getActions,
  type AcceptWithChangesDialogAction,
  type CommentDialogAction,
  type DirectAction,
  type StockRequestAction
} from '@/modules/admin/stock-requests/hooks/use-stock-request-actions.hook'

import {
  StockRequestStatus,
  type AcceptWithChangeRequestStatusDTO,
  type RejectStockRequestStatusDTO,
  type StockRequestResponse
} from '@/modules/admin/stock-requests/models/stock-requests.model'

import { useToast } from '@/core/components/ui/toast'
import { getRouteName } from '@/core/config/routes.config'
import type { LocalizedError } from '@/core/models/errors.model'
import { EmployeeRole } from '@/modules/admin/employees/models/employees.models'
import type { AxiosError } from 'axios'
import { useRouter } from 'vue-router'

/** Props */
const props = defineProps<{
  status: StockRequestStatus
  role: EmployeeRole
  request: StockRequestResponse
}>()

const actions = computed<StockRequestAction[]>(() => {
  return getActions(props.status, props.role)
})

const commentDialogOpen = ref(false)
const acceptChangesDialogOpen = ref(false)
let currentAction: StockRequestAction | null = null

const { toast } = useToast()

function onActionClick(action: StockRequestAction) {
  currentAction = action

  switch (action.type) {
    case 'DIRECT':
      callHandler(action)
      break

    case 'COMMENT_DIALOG':
      commentDialogOpen.value = true
      break

    case 'ACCEPT_WITH_CHANGES_DIALOG':
      acceptChangesDialogOpen.value = true
      break
  }
}

function closeCommentDialog() {
  commentDialogOpen.value = false
  currentAction = null
}

async function handleCommentSubmit(dto: RejectStockRequestStatusDTO) {
  if (currentAction?.type === 'COMMENT_DIALOG') {
    callHandler(currentAction, dto)
  }
  closeCommentDialog()
}

function closeAcceptDialog() {
  acceptChangesDialogOpen.value = false
  currentAction = null
}

async function handleAcceptWithChangesSubmit(dto: AcceptWithChangeRequestStatusDTO) {
  if (currentAction?.type === 'ACCEPT_WITH_CHANGES_DIALOG') {
    callHandler(currentAction, dto)
  }
  closeAcceptDialog()
}

const queryClient = useQueryClient()

type MutationVariables =
  | {
      action: DirectAction
      payload?: undefined
    }
  | {
      action: CommentDialogAction
      payload: RejectStockRequestStatusDTO
    }
  | {
      action: AcceptWithChangesDialogAction
      payload: AcceptWithChangeRequestStatusDTO
    }

const router = useRouter()

// Our single mutation for any action
const { mutate, isPending } = useMutation({
  mutationFn: async (vars: MutationVariables) => {
    const { action } = vars
    switch (action.type) {
      case 'DIRECT':
        return action.handler(props.request.requestId)
      case 'COMMENT_DIALOG':
        return action.handler(props.request.requestId, vars.payload as RejectStockRequestStatusDTO)
      case 'ACCEPT_WITH_CHANGES_DIALOG':
        return action.handler(props.request.requestId, vars.payload as AcceptWithChangeRequestStatusDTO)
    }
  },
  onSuccess: (_, { action }) => {
    queryClient.invalidateQueries({ queryKey: ['stock-requests'] })
    queryClient.invalidateQueries({ queryKey: ['stock-request', props.request.requestId] })

    toast({
      title: 'Успех',
      description: `Действие "${action.label}" выполнено успешно.`,
      variant: 'default',
    })

    if (action.redirectRouteKey) {
      router.push({ name: getRouteName(action.redirectRouteKey) })
    }
  },
  onError: (error: AxiosError<LocalizedError>) => {
    toast({
      title: "Ошибка",
      description: error.response?.data.message.ru ?? 'Не удалось выполнить действие. Попробуйте снова.' ,
      variant: 'destructive',
    });
  }
})

function callHandler(action: StockRequestAction, payload?: RejectStockRequestStatusDTO | AcceptWithChangeRequestStatusDTO): void {
  switch (action.type) {
    case 'DIRECT':
      return mutate({ action })
    case 'COMMENT_DIALOG':
      return mutate({
        action,
        payload: payload as RejectStockRequestStatusDTO,
      })
    case 'ACCEPT_WITH_CHANGES_DIALOG':
      return mutate({
        action,
        payload: payload as AcceptWithChangeRequestStatusDTO,
      })
  }
}
</script>
