<template>
	<AdminProductCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { getRouteName } from '@/core/config/routes.config'
import AdminProductCreateForm from '@/modules/admin/products/components/create/admin-product-create-form.vue'
import type { CreateProductDTO } from '@/modules/kiosk/products/models/product.model'
import { productsService } from '@/modules/kiosk/products/services/products.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()

const createMutation = useMutation({
	mutationFn: (dto: CreateProductDTO) => productsService.createProduct(dto),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-products'] })
		router.push({ name: getRouteName("ADMIN_PRODUCTS") })
	},
})

function handleCreate(dto: CreateProductDTO) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
