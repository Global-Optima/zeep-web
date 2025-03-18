<template>
	<AdminStoreEmployeesCreateForm
		@on-submit="handleCreate"
		@on-cancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import type { CreateEmployeeDTO } from '@/modules/admin/employees/models/employees.models'
import AdminStoreEmployeesCreateForm from '@/modules/admin/employees/stores/components/create/admin-store-employees-create-form.vue'
import { storeEmployeeService } from '@/modules/admin/employees/stores/services/store-employees.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRoute, useRouter } from 'vue-router'
import {type AxiosLocalizedError, useAxiosLocaleToast} from "@/core/hooks/use-axios-locale-toast.hooks";

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()
const { toastLocalizedError } = useAxiosLocaleToast()

const route = useRoute()
const storeId = route.params.storeId as string

const createMutation = useMutation({
	mutationFn: ({ dto, storeId }: { dto: CreateEmployeeDTO, storeId: number }) =>
		storeEmployeeService.createStoreEmployee(dto, storeId),
	onMutate: () => {
		toast({
			title: 'Добавление сотрудника...',
			description: 'Мы создаем нового сотрудника. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['store-employees'] })
		router.back()
		toast({
			title: 'Сотрудник добавлен!',
			description: 'Новый сотрудник успешно зарегистрирован в системе.',
		})
	},
  onError: (error: AxiosLocalizedError) => {
    toastLocalizedError(error, "Не удалось создать нового сотрудника. Попробуйте еще раз.")
  }
})

function handleCreate(dto: CreateEmployeeDTO) {
	if (!storeId) {
		toast({
			title: 'Ошибка',
			description: 'Не удалось определить магазин для сотрудника.',
			variant: 'destructive',
		})
		return
	}

	createMutation.mutate({ dto, storeId: Number(storeId) })
}

function handleCancel() {
	router.back()
}
</script>
