<template>
	<AdminStoreProductsCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import AdminStoreProductsCreateForm from '@/modules/admin/store-products/components/create/admin-store-products-create-form.vue'
import type { CreateStoreProductDTO } from '@/modules/admin/store-products/models/store-products.model'
import { storeProductsService } from '@/modules/admin/store-products/services/store-products.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()

const createMutation = useMutation({
	mutationFn: (dto: CreateStoreProductDTO[]) => storeProductsService.createMultipleStoreProducts(dto),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-store-products'] })
		router.back()
	},
})

function handleCreate(dto: CreateStoreProductDTO[]) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
