<template>
	<AdminEmployeesUpdateForm
		v-if="employee"
		:employee="employee"
		@on-submit="handleUpdate"
		@on-cancel="handleCancel"
	/>

	<p v-else>Сотрудник не найден</p>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminEmployeesUpdateForm from '@/modules/admin/store-employees/components/update/admin-employees-update-form.vue'
import type { UpdateEmployeeDTO } from '@/modules/admin/store-employees/models/employees.models'
import { employeesService } from '@/modules/admin/store-employees/services/employees.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const route = useRoute()
const employeeId = route.params.id as string

const { data: employee } = useQuery({
	queryKey: ['store-employee', employeeId],
	queryFn: () => employeesService.getStoreEmployeeById(Number(employeeId)),
	enabled: !!employeeId,
})

const updateMutation = useMutation({
	mutationFn: (newStoreData: UpdateEmployeeDTO) =>
		employeesService.updateStoreEmployee(Number(employeeId), newStoreData),
	onMutate: () => {
		toast({
			title: 'Обновление...',
			description: 'Обновление данных сотрудника. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['store-employees'] })
		queryClient.invalidateQueries({ queryKey: ['store-employee', employeeId] })
		toast({
			title: 'Успех!',
			description: 'Данные сотрудника успешно обновлены.',
		})
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при обновлении данных сотрудника.',
			variant: 'destructive',
		})
	},
})

function handleUpdate(dto: UpdateEmployeeDTO) {
	if (!employeeId) {
		toast({
			title: 'Ошибка',
			description: 'Неверный идентификатор сотрудника.',
			variant: 'destructive',
		})
		return
	}

	updateMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
