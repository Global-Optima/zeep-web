<template>
	<AdminProductSizeCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminProductSizeCreateForm, { type CreateProductSizeFormSchema } from '@/modules/admin/products/components/create/admin-product-size-create-form.vue'
import type { CreateProductSizeDTO } from '@/modules/kiosk/products/models/product.model'
import { productsService } from '@/modules/kiosk/products/services/products.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const route = useRoute()
const { toast } = useToast()

const productId = route.query.productId as string

const createMutation = useMutation({
	mutationFn: (dto: CreateProductSizeDTO) => productsService.createProductSize(dto),
	onMutate: () => {
		toast({
			title: 'Создание...',
			description: 'Создание нового размера продукта. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-product-sizes', Number(productId)] })
		queryClient.invalidateQueries({ queryKey: ['admin-product-details', productId] })
		toast({
			title: 'Успех!',
variant: 'success',
			description: 'Размер продукта успешно создан.',
		})
		router.back()
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при создании размера продукта.',
			variant: 'destructive',
		})
	},
})

function handleCreate(data: CreateProductSizeFormSchema) {
	if (isNaN(Number(productId))) {
		toast({
			title: 'Ошибка',
			description: 'Неверный идентификатор продукта.',
			variant: 'destructive',
		})
		return router.back()
	}

	const dto: CreateProductSizeDTO = {
		productId: Number(productId),
		name: data.name,
		unitId: data.unitId,
		basePrice: data.basePrice,
		size: data.size,
    machineId: data.machineId,
		isDefault: false,
		additives: data.additives.map(a => ({ additiveId: a.additiveId, isDefault: a.isDefault })),
		ingredients: data.ingredients.map(a => ({ ingredientId: a.ingredientId, quantity: a.quantity })),
	}

	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
