<template>
	<AdminProductCategoryCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import AdminProductCategoryCreateForm from '@/modules/admin/product-categories/components/create/admin-product-category-create-form.vue'
import type { CreateProductCategoryDTO } from '@/modules/kiosk/products/models/product.model'
import { productsService } from '@/modules/kiosk/products/services/products.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()

const createMutation = useMutation({
	mutationFn: (dto: CreateProductCategoryDTO) => productsService.createProductCategory(dto),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-product-categories'] })
		router.back()
	},
})

function handleCreate(dto: CreateProductCategoryDTO) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
