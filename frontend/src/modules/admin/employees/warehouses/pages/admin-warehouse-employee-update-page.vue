<template>
	<AdminWarehouseEmployeesUpdateForm
		v-if="employee"
		:employee="employee"
		@on-submit="handleUpdate"
		@on-cancel="handleCancel"
	/>

	<p v-else>Сотрудник не найден</p>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminWarehouseEmployeesUpdateForm from '@/modules/admin/employees/warehouses/components/update/admin-warehouse-employees-update-form.vue'
import type { UpdateWarehouseEmployeeDTO } from '@/modules/admin/employees/warehouses/models/warehouse-employees.model'
import { warehouseEmployeeService } from '@/modules/admin/employees/warehouses/services/warehouse-employees.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const route = useRoute()
const employeeId = route.params.id as string

const { data: employee } = useQuery({
	queryKey: ['warehouse-employee', employeeId],
	queryFn: () => warehouseEmployeeService.getWarehouseEmployeeById(Number(employeeId)),
	enabled: !!employeeId,
})

const updateMutation = useMutation({
	mutationFn: (newStoreData: UpdateWarehouseEmployeeDTO) => warehouseEmployeeService.updateWarehouseEmployee(Number(employeeId), newStoreData),
	onMutate: () => {
		toast({
			title: 'Обновление...',
			description: 'Обновление данных сотрудника. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['warehouse-employees'] })
		queryClient.invalidateQueries({ queryKey: ['warehouse-employee', employeeId] })
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

function handleUpdate(dto: UpdateWarehouseEmployeeDTO) {
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
