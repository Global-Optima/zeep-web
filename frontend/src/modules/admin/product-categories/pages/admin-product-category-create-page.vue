<template>
	<AdminProductCategoryCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminProductCategoryCreateForm from '@/modules/admin/product-categories/components/create/admin-product-category-create-form.vue'
import type { CreateProductCategoryDTO } from '@/modules/kiosk/products/models/product.model'
import { productsService } from '@/modules/kiosk/products/services/products.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const createMutation = useMutation({
	mutationFn: (dto: CreateProductCategoryDTO) => productsService.createProductCategory(dto),
	onMutate: () => {
		toast({
			title: 'Создание...',
			description: 'Создание новой категории продукта. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-product-categories'] })
		toast({
			title: 'Успех!',
			description: 'Категория продукта успешно создана.',
		})
		router.back()
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при создании категории продукта.',
			variant: 'destructive',
		})
	},
})

function handleCreate(dto: CreateProductCategoryDTO) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
