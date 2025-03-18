<template>
	<AdminAdminEmployeesCreateForm
		@on-submit="handleCreate"
		@on-cancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminAdminEmployeesCreateForm from '@/modules/admin/employees/admins/components/create/admin-admin-employees-create-form.vue'
import { adminEmployeeService } from '@/modules/admin/employees/admins/services/admin-employees.service'
import type { CreateEmployeeDTO } from '@/modules/admin/employees/models/employees.models'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'
import {type AxiosLocalizedError, useAxiosLocaleToast} from "@/core/hooks/use-axios-locale-toast.hooks";

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()
const { toastLocalizedError } = useAxiosLocaleToast()

const createMutation = useMutation({
	mutationFn: ({ dto }: { dto: CreateEmployeeDTO }) =>
		adminEmployeeService.createAdminEmployee(dto),
	onMutate: () => {
		toast({
			title: 'Добавление сотрудника...',
			description: 'Мы создаем нового сотрудника. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-employees'] })
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
	createMutation.mutate({ dto})
}

function handleCancel() {
	router.back()
}
</script>
