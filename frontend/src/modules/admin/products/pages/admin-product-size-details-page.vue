<template>
	<p v-if="!productSizeDetails">Размер не найден</p>

	<AdminProductSizeUpdateForm
		v-else
		:productSize="productSizeDetails"
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { type CreateProductSizeFormSchema } from '@/modules/admin/products/components/create/admin-product-size-create-form.vue'
import AdminProductSizeUpdateForm from '@/modules/admin/products/components/update/admin-product-size-update-form.vue'
import type { UpdateProductSizeDTO } from '@/modules/kiosk/products/models/product.model'
import { productsService } from '@/modules/kiosk/products/services/products.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const route = useRoute()

const productSizeId = route.params.id as string
const productId = route.query.productId as string

console.log(productSizeId)

const { data: productSizeDetails } = useQuery({
  queryKey: ['admin-product-size-details', productSizeId],
	queryFn: () => productsService.getProductSizeById(Number(productSizeId)),
  enabled: !isNaN(Number(productSizeId)),
})

const updateMutation = useMutation({
	mutationFn: ({sizeId, dto}: {sizeId: number, dto: UpdateProductSizeDTO}) => productsService.updateProductSize(sizeId, dto),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-product-sizes', Number(productId)] })
    queryClient.invalidateQueries({ queryKey: ['admin-product-details', productId] })
    queryClient.invalidateQueries({ queryKey: ['admin-product-size-details', productSizeId] })
		router.back()
	},
})

function handleCreate(data: CreateProductSizeFormSchema) {
  if (isNaN(Number(productSizeId))) return router.back()

  const dto: UpdateProductSizeDTO = {
    name: data.name,
    measure: data.measure,
    basePrice: data.basePrice,
    size: data.size,
    additives: data.additives.map(a => ({additiveId: a.additiveId, isDefault: a.isDefault ?? false})),
  }

	updateMutation.mutate({sizeId: Number(productSizeId), dto})
}

function handleCancel() {
	router.back()
}
</script>
