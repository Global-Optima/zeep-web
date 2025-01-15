<template>
	<div>
		<p v-if="!productDetails">Товар не найден</p>

		<Tabs
			v-else
			default-value="details"
		>
			<TabsList class="grid grid-cols-2 mx-auto mb-6 w-full max-w-6xl">
				<TabsTrigger
					class="py-2"
					value="details"
				>
					Детали
				</TabsTrigger>
				<TabsTrigger
					class="py-2"
					value="variants"
				>
					Варианты
				</TabsTrigger>
			</TabsList>
			<TabsContent value="details">
				<AdminProductDetailsForm
					:product-details="productDetails"
					@on-submit="onUpdate"
					@on-cancel="onCancel"
				/>
			</TabsContent>

			<TabsContent value="variants">
				<AdminProductsVariants
					:product-details="productDetails"
					@on-cancel="onCancel"
				/>
			</TabsContent>
		</Tabs>
	</div>
</template>

<script setup lang="ts">
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '@/core/components/ui/tabs'
import { getRouteName } from '@/core/config/routes.config'
import AdminProductDetailsForm from "@/modules/admin/products/components/details/admin-product-details-form.vue"
import AdminProductsVariants from "@/modules/admin/products/components/details/admin-products-variants.vue"
import type { UpdateProductDTO } from '@/modules/kiosk/products/models/product.model'
import { productsService } from '@/modules/kiosk/products/services/products.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const queryClient = useQueryClient()

const productId = route.params.id as string

const { data: productDetails } = useQuery({
  queryKey: ['admin-product-details', productId],
	queryFn: () => productsService.getProductDetails(Number(productId)),
  enabled: !isNaN(Number(productId)),
})

const productUpdateMutation = useMutation({
	mutationFn: ({id, dto} : {id: number, dto: UpdateProductDTO}) => productsService.updateProduct(id, dto),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-products'] })
    queryClient.invalidateQueries({ queryKey: ['admin-product-details', productId] })
		router.push({ name: getRouteName("ADMIN_PRODUCTS") })
	},
})

function onUpdate(dto: UpdateProductDTO) {
	productUpdateMutation.mutate({id: Number(productId), dto})
}

function onCancel() {
	router.back()
}
</script>
