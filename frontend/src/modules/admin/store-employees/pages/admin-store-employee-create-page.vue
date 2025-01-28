<template>
	<AdminStoreEmployeesCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import { getRouteName } from '@/core/config/routes.config'
import AdminStoreEmployeesCreateForm from '@/modules/admin/store-employees/components/create/admin-store-employees-create-form.vue'
import { type CreateStoreEmployeeDTO } from '@/modules/admin/store-employees/models/employees.models'
import { employeesService } from '@/modules/admin/store-employees/services/employees.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const createMutation = useMutation({
	mutationFn: (newStoreData: CreateStoreEmployeeDTO) => employeesService.createStoreEmployee(newStoreData),
	onMutate: () => {
		toast({
			title: 'Создание...',
			description: 'Создание нового сотрудника. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['store-employees'] })
		toast({
			title: 'Успех!',
			description: 'Сотрудник успешно добавлен.',
		})
		router.push({ name: getRouteName('ADMIN_STORE_EMPLOYEES') })
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при добавлении сотрудника.',
			variant: 'destructive',
		})
	},
})

function handleCreate(dto: CreateStoreEmployeeDTO) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
