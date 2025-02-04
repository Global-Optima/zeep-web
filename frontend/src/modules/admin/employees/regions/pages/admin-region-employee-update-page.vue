<template>
	<AdminRegionsEmployeesUpdateForm
		v-if="employee"
		:employee="employee"
		@on-submit="handleUpdate"
		@on-cancel="handleCancel"
	/>

	<p v-else>Сотрудник не найден</p>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminRegionsEmployeesUpdateForm from '@/modules/admin/employees/regions/components/update/admin-regions-employees-update-form.vue'
import type { UpdateRegionEmployeeDTO } from '@/modules/admin/employees/regions/models/region-employees.model'
import { regionEmployeeService } from '@/modules/admin/employees/regions/services/region-employees.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const route = useRoute()
const employeeId = route.params.id as string

const { data: employee } = useQuery({
	queryKey: ['region-employee', employeeId],
	queryFn: () => regionEmployeeService.getRegionEmployeeById(Number(employeeId)),
	enabled: !!employeeId,
})

const updateMutation = useMutation({
	mutationFn: (newStoreData: UpdateRegionEmployeeDTO) => regionEmployeeService.updateRegionEmployee(Number(employeeId), newStoreData),
	onMutate: () => {
		toast({
			title: 'Обновление...',
			description: 'Обновление данных сотрудника. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['region-employees'] })
		queryClient.invalidateQueries({ queryKey: ['region-employee', employeeId] })
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

function handleUpdate(dto: UpdateRegionEmployeeDTO) {
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
