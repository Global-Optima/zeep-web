<template>
	<p v-if="!productSizeDetails">Размер не найден</p>

	<AdminProductSizeUpdateForm
		v-else
		:productSize="productSizeDetails"
		@onSubmit="handleUpdate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminProductSizeUpdateForm, { type UpdateProductSizeFormSchema } from '@/modules/admin/products/components/details/admin-product-size-details-form.vue'
import type { UpdateProductSizeDTO } from '@/modules/kiosk/products/models/product.model'
import { productsService } from '@/modules/kiosk/products/services/products.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const route = useRoute()
const { toast } = useToast()

const productSizeId = route.params.id as string
const productId = route.query.productId as string

const { data: productSizeDetails } = useQuery({
	queryKey: ['admin-product-size-details', productSizeId],
	queryFn: () => productsService.getProductSizeById(Number(productSizeId)),
	enabled: !isNaN(Number(productSizeId)),
})

const updateMutation = useMutation({
	mutationFn: ({ sizeId, dto }: { sizeId: number; dto: UpdateProductSizeDTO }) =>
		productsService.updateProductSize(sizeId, dto),
	onMutate: () => {
		toast({
			title: 'Обновление...',
			description: 'Обновление данных размера продукта. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-product-sizes', Number(productId)] })
		queryClient.invalidateQueries({ queryKey: ['admin-product-details', productId] })
		queryClient.invalidateQueries({ queryKey: ['admin-product-size-details', productSizeId] })
		toast({
			title: 'Успех!',
variant: 'success',
			description: 'Данные размера успешно обновлены.',
		})
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при обновлении размера продукта.',
			variant: 'destructive',
		})
	},
})

function handleUpdate(data: UpdateProductSizeFormSchema) {
	if (isNaN(Number(productSizeId))) {
		toast({
			title: 'Ошибка',
			description: 'Неверный идентификатор размера.',
			variant: 'destructive',
		})
		return router.back()
	}

	const dto: UpdateProductSizeDTO = {
		name: data.name,
		basePrice: data.basePrice,
		size: data.size,
		unitId: data.unitId,
    machineId: data.machineId,
		additives: data.additives.map(a => ({ additiveId: a.additiveId, isDefault: a.isDefault ?? false })),
		ingredients: data.ingredients.map(a => ({ ingredientId: a.ingredientId, quantity: a.quantity })),
	}

	updateMutation.mutate({ sizeId: Number(productSizeId), dto })
}

function handleCancel() {
	router.back()
}
</script>
