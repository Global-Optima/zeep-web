<template>
	<AdminWarehouseCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import type { CreateWarehouseDTO } from '@/modules/admin/warehouses/models/warehouse.model'
import { warehouseService } from '@/modules/admin/warehouses/services/warehouse.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const createMutation = useMutation({
	mutationFn: (dto: CreateWarehouseDTO) => warehouseService.create(dto),
	onMutate: () => {
		toast({
			title: 'Создание...',
			description: 'Создание нового склада. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-warehouses'] })
		toast({
			title: 'Успех!',
			description: 'Склад успешно создан.',
		})
		router.back()
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при создании склада.',
			variant: 'destructive',
		})
	},
})

function handleCreate(dto: CreateWarehouseDTO) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
