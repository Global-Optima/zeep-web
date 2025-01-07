<template>
	<AdminProductSizeCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import AdminProductSizeCreateForm from '@/modules/admin/products/components/create/admin-product-size-create-form.vue'
import type { CreateProductSizeDTO } from '@/modules/kiosk/products/models/product.model'
import { productsService } from '@/modules/kiosk/products/services/products.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const route = useRoute()

const productId = route.query.productId as string

const createMutation = useMutation({
	mutationFn: (dto: CreateProductSizeDTO) => productsService.createProductSize(dto),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-product-sizes', Number(productId)] })
    queryClient.invalidateQueries({ queryKey: ['admin-product-details', productId] })
		router.back()
	},
})

function handleCreate(dto: CreateProductSizeDTO) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
