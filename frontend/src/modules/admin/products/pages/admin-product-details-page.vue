<template>
	<div>
		<PageLoader v-if="isLoading" />
		<p v-else-if="!productDetails">Продукт не найден</p>

		<Tabs
			v-else
			default-value="details"
		>
			<TabsList class="grid grid-cols-2 mx-auto mb-6 w-full max-w-6xl">
				<TabsTrigger
					class="py-2"
					value="details"
					>Детали</TabsTrigger
				>
				<TabsTrigger
					class="py-2"
					value="variants"
					>Размеры</TabsTrigger
				>
			</TabsList>
			<TabsContent value="details">
				<AdminProductDetailsForm
					ref="formRef"
					:product-details="productDetails"
					:is-submitting="isPending"
					:readonly="!canUpdate"
					@on-submit="onUpdate"
					@on-cancel="onCancel"
				/>
			</TabsContent>

			<TabsContent value="variants">
				<AdminProductsVariants
					:product-details="productDetails"
					:readonly="!canUpdate"
					@on-cancel="onCancel"
				/>
			</TabsContent>
		</Tabs>
	</div>
</template>

<script setup lang="ts">
import PageLoader from "@/core/components/page-loader/PageLoader.vue"
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '@/core/components/ui/tabs'
import { useToast } from '@/core/components/ui/toast/use-toast'
import { useHasRole } from '@/core/hooks/use-has-roles.hook'
import { EmployeeRole } from '@/modules/admin/employees/models/employees.models'
import AdminProductDetailsForm from '@/modules/admin/products/components/details/admin-product-details-form.vue'
import AdminProductsVariants from '@/modules/admin/products/components/details/admin-product-sizes.vue'
import type { UpdateProductDTO } from '@/modules/kiosk/products/models/product.model'
import { productsService } from '@/modules/kiosk/products/services/products.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { useTemplateRef } from "vue"
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()
const canUpdate = useHasRole([EmployeeRole.ADMIN])

const productId = route.params.id as string

const { data: productDetails, isLoading } = useQuery({
	queryKey: ['admin-product-details', productId],
	queryFn: () => productsService.getProductDetails(Number(productId)),
	enabled: !isNaN(Number(productId)),
})

const formRef = useTemplateRef<InstanceType<typeof AdminProductDetailsForm>>('formRef')

const {mutate, isPending} = useMutation({
	mutationFn: ({ id, dto }: { id: number; dto: UpdateProductDTO }) =>
		productsService.updateProduct(id, dto),
	onMutate: () => {
		toast({
			title: 'Обновление...',
			description: 'Обновление данных продукта. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-products'] })
		queryClient.invalidateQueries({ queryKey: ['admin-product-details', productId] })

    formRef.value?.resetFormValues();

    toast({
			title: 'Успех!',
variant: 'success',
			description: 'Данные продукта успешно обновлены.',
		})
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при обновлении продукта.',
			variant: 'destructive',
		})
	},
})

function onUpdate(dto: UpdateProductDTO) {
	if (isNaN(Number(productId))) {
		toast({
			title: 'Ошибка',
			description: 'Неверный идентификатор продукта.',
			variant: 'destructive',
		})
		return router.back()
	}

  mutate({ id: Number(productId), dto })
}

function onCancel() {
	router.back()
}
</script>
