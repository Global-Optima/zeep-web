<template>
	<AdminProductCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
		:isSubmitting="isPending"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import { getRouteName } from '@/core/config/routes.config'
import AdminProductCreateForm from '@/modules/admin/products/components/create/admin-product-create-form.vue'
import type { CreateProductDTO } from '@/modules/kiosk/products/models/product.model'
import { productsService } from '@/modules/kiosk/products/services/products.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const {mutate, isPending} = useMutation({
	mutationFn: (dto: CreateProductDTO) => productsService.createProduct(dto),
	onMutate: () => {
		toast({
			title: 'Создание...',
			description: 'Пожалуйста, подождите, создается новый продукт.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-products'] })
		toast({
			title: 'Успех!',
      variant: 'success',
			description: 'Продукт успешно создан.',
		})
		router.push({ name: getRouteName('ADMIN_PRODUCTS') })
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при создании продукта.',
			variant: 'destructive',
		})
	},
})

function handleCreate(dto: CreateProductDTO) {
	mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
