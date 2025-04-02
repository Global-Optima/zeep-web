<template>
	<p v-if="!storeProductDetails || !productDetails">Продукт не найден</p>

	<AdminStoreProductDetailsForm
		v-else
		:initialStoreProduct="storeProductDetails"
		:product="productDetails"
		@onSubmit="handleUpdate"
		@onCancel="handleCancel"
		:readonly="!canUpdate"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import { useHasRole } from '@/core/hooks/use-has-roles.hook'
import { EmployeeRole } from '@/modules/admin/employees/models/employees.models'
import AdminStoreProductDetailsForm from '@/modules/admin/store-products/components/details/admin-store-product-details-form.vue'
import type { UpdateStoreProductDTO } from '@/modules/admin/store-products/models/store-products.model'
import { storeProductsService } from '@/modules/admin/store-products/services/store-products.service'
import { productsService } from '@/modules/kiosk/products/services/products.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const storeProductId = route.params.id as string

const canUpdate = useHasRole([EmployeeRole.STORE_MANAGER])

const { data: storeProductDetails } = useQuery({
	queryKey: computed(() => ['admin-store-product-details', storeProductId]),
	queryFn: () => storeProductsService.getStoreProduct(Number(storeProductId)),
	enabled: !isNaN(Number(storeProductId)),
})

const { data: productDetails } = useQuery({
	queryKey: computed(() => ['admin-product-details', storeProductDetails.value?.productId]),
	queryFn: () => productsService.getProductDetails(storeProductDetails.value!.productId),
	enabled: computed(() => !!storeProductDetails.value),
})

const updateMutation = useMutation({
	mutationFn: ({ id, dto }: { id: number; dto: UpdateStoreProductDTO }) =>
		storeProductsService.updateStoreProduct(id, dto),
	onMutate: () => {
		toast({
			title: 'Обновление...',
			description: 'Обновление данных продукта кафе. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-store-products'] })
		queryClient.invalidateQueries({ queryKey: ['admin-store-product-details', storeProductId] })
		toast({
			title: 'Успех!',
variant: 'success',
			description: 'Данные продукта кафе успешно обновлены.',
		})
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при обновлении данных продукта кафе.',
			variant: 'destructive',
		})
	},
})

function handleUpdate(updatedData: UpdateStoreProductDTO) {
	if (!storeProductDetails.value?.id) {
		toast({
			title: 'Ошибка',
			description: 'Неверный идентификатор продукта.',
			variant: 'destructive',
		})
		return router.back()
	}

	updateMutation.mutate({ id: Number(storeProductId), dto: updatedData })
}

function handleCancel() {
	router.back()
}
</script>
