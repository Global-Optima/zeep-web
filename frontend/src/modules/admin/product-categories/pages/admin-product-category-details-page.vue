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
import { getRouteName } from '@/core/config/routes.config'
import AdminProductCategoryDetailsForm from '@/modules/admin/product-categories/components/details/admin-product-category-details-form.vue'
import type { UpdateProductCategoryDTO } from '@/modules/kiosk/products/models/product.model'
import { productsService } from '@/modules/kiosk/products/services/products.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const categoryId = route.params.id as string

const queryClient = useQueryClient()

const { data: categoryDetails } = useQuery({
  queryKey: computed(() => ['admin-product-ingredient-details', categoryId]),
	queryFn: () => productsService.getProductCategoryByID(Number(categoryId)),
  enabled: !isNaN(Number(categoryId)),
})

const updateMutation = useMutation({
	mutationFn: ({id, dto}:{id: number, dto: UpdateProductCategoryDTO}) => {
    return productsService.updateProductCategory(id, dto)
  },
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-product-categories'] })
		queryClient.invalidateQueries({ queryKey: ['admin-product-ingredient-details', categoryId] })
		router.push({ name: getRouteName("ADMIN_PRODUCT_CATEGORIES") })
	},
})

function handleUpdate(updatedData: UpdateProductCategoryDTO) {
  if (isNaN(Number(categoryId))) return router.back()

	updateMutation.mutate({id: Number(categoryId), dto: updatedData})
}

function handleCancel() {
	router.back()
}
</script>
