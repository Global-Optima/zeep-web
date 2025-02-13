<template>
	<AdminWarehouseEmployeesCreateForm
		@on-submit="handleCreate"
		@on-cancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import type { CreateEmployeeDTO } from '@/modules/admin/employees/models/employees.models'
import AdminWarehouseEmployeesCreateForm from '@/modules/admin/employees/warehouses/components/create/admin-warehouse-employees-create-form.vue'
import { warehouseEmployeeService } from '@/modules/admin/employees/warehouses/services/warehouse-employees.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const route = useRoute()
const warehouseId = route.params.warehouseId as string

const createMutation = useMutation({
	mutationFn: ({ dto, warehouseId }: { dto: CreateEmployeeDTO, warehouseId: number }) =>
		warehouseEmployeeService.createWarehouseEmployee(dto, warehouseId),
	onMutate: () => {
		toast({
			title: 'Добавление сотрудника...',
			description: 'Мы создаем нового сотрудника. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['warehouse-employees'] })
		router.back()
		toast({
			title: 'Сотрудник добавлен!',
			description: 'Новый сотрудник успешно зарегистрирован в системе.',
		})
	},
	onError: () => {
		toast({
			title: 'Ошибка при создании',
			description: 'Не удалось создать нового сотрудника. Попробуйте еще раз.',
			variant: 'destructive',
		})
	},
})

function handleCreate(dto: CreateEmployeeDTO) {
	if (!warehouseId) {
		toast({
			title: 'Ошибка',
			description: 'Не удалось определить склад для сотрудника.',
			variant: 'destructive',
		})
		return
	}

	createMutation.mutate({ dto, warehouseId: Number(warehouseId) })
}

function handleCancel() {
	router.back()
}
</script>
