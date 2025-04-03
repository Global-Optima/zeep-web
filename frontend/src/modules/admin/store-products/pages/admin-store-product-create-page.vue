<template>
	<AdminStoreProductsCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminStoreProductsCreateForm from '@/modules/admin/store-products/components/create/admin-store-products-create-form.vue'
import type { CreateStoreProductDTO } from '@/modules/admin/store-products/models/store-products.model'
import { storeProductsService } from '@/modules/admin/store-products/services/store-products.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const createMutation = useMutation({
	mutationFn: (dto: CreateStoreProductDTO[]) =>
		storeProductsService.createMultipleStoreProducts(dto),
	onMutate: () => {
		toast({
			title: 'Создание...',
			description: 'Добавление новых продуктов кафе. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-store-products'] })
		toast({
			title: 'Успех!',
variant: 'success',
			description: 'Продукты кафе успешно добавлены.',
		})
		router.back()
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при добавлении продуктов кафе.',
			variant: 'destructive',
		})
	},
})

function handleCreate(dto: CreateStoreProductDTO[]) {
	if (dto.length === 0) {
		toast({
			title: 'Ошибка',
			description: 'Список продуктов пуст. Пожалуйста, добавьте продукты перед сохранением.',
			variant: 'destructive',
		})
		return
	}

	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
