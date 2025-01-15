<template>
	<p v-if="!storeProductDetails || !productDetails">Товар не найден</p>

	<AdminStoreProductDetailsForm
		v-else
		:initialStoreProduct="storeProductDetails"
		:product="productDetails"
		@onSubmit="handleUpdate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { getRouteName } from '@/core/config/routes.config'
import AdminStoreProductDetailsForm from '@/modules/admin/store-products/components/details/admin-store-product-details-form.vue'
import type { UpdateStoreProductDTO } from '@/modules/admin/store-products/models/store-products.model'
import { storeProductsService } from '@/modules/admin/store-products/services/store-products.service'
import { productsService } from '@/modules/kiosk/products/services/products.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const storeProductId = route.params.id as string

const queryClient = useQueryClient()

const { data: storeProductDetails } = useQuery({
  queryKey: computed(() => ['admin-store-product-details', storeProductId]),
	queryFn: () => storeProductsService.getStoreProduct(Number(storeProductId)),
  enabled: !isNaN(Number(storeProductId)),
})

const { data: productDetails } = useQuery({
  queryKey: computed(() => ['admin-product-details', storeProductDetails.value?.productId]),
	queryFn: () => productsService.getProductDetails(Number(storeProductDetails.value?.productId)),
  enabled: computed(() => Boolean(storeProductDetails.value?.productId)),
})


const updateMutation = useMutation({
	mutationFn: ({id, dto}:{id: number, dto: UpdateStoreProductDTO}) => {
    return storeProductsService.updateStoreProduct(id, dto)
  },
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-store-products'] })
		queryClient.invalidateQueries({ queryKey: ['admin-store-product-details', storeProductId] })
		router.push({ name: getRouteName("ADMIN_STORE_PRODUCTS") })
	},
})

function handleUpdate(updatedData: UpdateStoreProductDTO) {
  if (storeProductDetails.value?.productId) return router.back()

	updateMutation.mutate({id: Number(storeProductId), dto: updatedData})
}

function handleCancel() {
	router.back()
}
</script>
