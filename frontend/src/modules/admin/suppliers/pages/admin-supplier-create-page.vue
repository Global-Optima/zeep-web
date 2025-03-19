<template>
	<AdminSupplierCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminSupplierCreateForm from '@/modules/admin/suppliers/components/create/admin-supplier-create-form.vue'
import type { CreateSupplierDTO } from '@/modules/admin/suppliers/models/suppliers.model'
import { suppliersService } from '@/modules/admin/suppliers/services/suppliers.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const createMutation = useMutation({
	mutationFn: (dto: CreateSupplierDTO) => suppliersService.createSupplier(dto),
	onMutate: () => {
		toast({
			title: 'Создание...',
			description: 'Пожалуйста, подождите, создается новый поставщик.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-suppliers'] })
		toast({
			title: 'Успех!',
variant: 'success',
			description: 'Поставщик успешно создан.',
		})
		router.back()
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при создании поставщика.',
			variant: 'destructive',
		})
	},
})

function handleCreate(dto: CreateSupplierDTO) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
