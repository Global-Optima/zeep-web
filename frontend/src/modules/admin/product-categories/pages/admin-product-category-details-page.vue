<template>
	<p v-if="!categoryDetails">Категория не найдена</p>

	<AdminProductCategoryDetailsForm
		v-else
		:productCategory="categoryDetails"
		@onSubmit="handleUpdate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminProductCategoryDetailsForm from '@/modules/admin/product-categories/components/details/admin-product-category-details-form.vue'
import type { UpdateProductCategoryDTO } from '@/modules/kiosk/products/models/product.model'
import { productsService } from '@/modules/kiosk/products/services/products.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const categoryId = route.params.id as string

const { data: categoryDetails } = useQuery({
	queryKey: computed(() => ['admin-product-categories-details', categoryId]),
	queryFn: () => productsService.getProductCategoryByID(Number(categoryId)),
	enabled: !isNaN(Number(categoryId)),
})

const updateMutation = useMutation({
	mutationFn: ({ id, dto }: { id: number; dto: UpdateProductCategoryDTO }) =>
		productsService.updateProductCategory(id, dto),
	onMutate: () => {
		toast({
			title: 'Обновление...',
			description: 'Обновление данных категории продукта. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-product-categories'] })
		queryClient.invalidateQueries({ queryKey: ['admin-product-categories-details', categoryId] })
		toast({
			title: 'Успех!',
			description: 'Данные категории успешно обновлены.',
		})
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при обновлении категории продукта.',
			variant: 'destructive',
		})
	},
})

function handleUpdate(updatedData: UpdateProductCategoryDTO) {
	if (isNaN(Number(categoryId))) {
		toast({
			title: 'Ошибка',
			description: 'Неверный идентификатор категории.',
			variant: 'destructive',
		})
		return router.back()
	}

	updateMutation.mutate({ id: Number(categoryId), dto: updatedData })
}

function handleCancel() {
	router.back()
}
</script>
