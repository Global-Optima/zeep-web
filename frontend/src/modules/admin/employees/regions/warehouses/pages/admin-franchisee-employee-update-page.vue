<template>
	<AdminFranchiseesEmployeesUpdateForm
		v-if="employee"
		:employee="employee"
		@on-submit="handleUpdate"
		@on-cancel="handleCancel"
	/>

	<p v-else>Сотрудник не найден</p>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminFranchiseesEmployeesUpdateForm from '@/modules/admin/employees/franchisees/components/update/admin-franchisees-employees-update-form.vue'
import type { UpdateFranchiseeEmployeeDTO } from '@/modules/admin/employees/franchisees/models/franchisees-employees.model'
import { franchiseeEmployeeService } from '@/modules/admin/employees/franchisees/services/franchisee-employees.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const route = useRoute()
const employeeId = route.params.id as string

const { data: employee } = useQuery({
	queryKey: ['franchisee-employee', employeeId],
	queryFn: () => franchiseeEmployeeService.getFranchiseeEmployeeById(Number(employeeId)),
	enabled: !!employeeId,
})

const updateMutation = useMutation({
	mutationFn: (newStoreData: UpdateFranchiseeEmployeeDTO) => franchiseeEmployeeService.updateFranchiseeEmployee(Number(employeeId), newStoreData),
	onMutate: () => {
		toast({
			title: 'Обновление...',
			description: 'Обновление данных сотрудника. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['franchisee-employees'] })
		queryClient.invalidateQueries({ queryKey: ['franchisee-employee', employeeId] })
		toast({
			title: 'Успех!',
variant: 'success',
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

function handleUpdate(dto: UpdateFranchiseeEmployeeDTO) {
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
