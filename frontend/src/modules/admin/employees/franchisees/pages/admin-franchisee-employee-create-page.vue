<template>
	<AdminFranchiseeEmployeesCreateForm
		@on-submit="handleCreate"
		@on-cancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminFranchiseeEmployeesCreateForm from '@/modules/admin/employees/franchisees/components/create/admin-franchisee-employees-create-form.vue'
import { franchiseeEmployeeService } from '@/modules/admin/employees/franchisees/services/franchisee-employees.service'
import type { CreateEmployeeDTO } from '@/modules/admin/employees/models/employees.models'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRoute, useRouter } from 'vue-router'
import {type AxiosLocalizedError, useAxiosLocaleToast} from "@/core/hooks/use-axios-locale-toast.hooks";

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()
const { toastLocalizedError } = useAxiosLocaleToast()
const route = useRoute()
const franchiseeId = route.params.franchiseeId as string

const createMutation = useMutation({
	mutationFn: ({ dto, franchiseeId }: { dto: CreateEmployeeDTO, franchiseeId: number }) =>
		franchiseeEmployeeService.createFranchiseeEmployee(dto, franchiseeId),
	onMutate: () => {
		toast({
			title: 'Добавление сотрудника...',
			description: 'Мы создаем нового сотрудника. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['franchisee-employees'] })
		router.back()
		toast({
			title: 'Сотрудник добавлен!',
			description: 'Новый сотрудник успешно зарегистрирован в системе.',
		})
	},
  onError: (error: AxiosLocalizedError) => {
    toastLocalizedError(error, "Произошла ошибка при обновлении материалов запроса на склад.")
  }
})

function handleCreate(dto: CreateEmployeeDTO) {
	if (!franchiseeId) {
		toast({
			title: 'Ошибка',
			description: 'Не удалось определить франчайзи для сотрудника.',
			variant: 'destructive',
		})
		return
	}

	createMutation.mutate({ dto, franchiseeId: Number(franchiseeId) })
}

function handleCancel() {
	router.back()
}
</script>
