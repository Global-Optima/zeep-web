<template>
	<AdminProductSizeCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
		:initialProductSize="productSizeDetails"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminProductSizeCreateForm, { type CreateProductSizeFormSchema } from '@/modules/admin/products/components/create/admin-product-size-create-form.vue'
import type {
  CreateProductSizeDTO,
  SelectedAdditiveDTO,
  SelectedIngredientDTO,
  SelectedProvisionDTO
} from '@/modules/kiosk/products/models/product.model'
import { productsService } from '@/modules/kiosk/products/services/products.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const route = useRoute()
const { toast } = useToast()

const templateProductSizeId = route.query.templateProductSizeId as string

const { data: productSizeDetails } = useQuery({
	queryKey: ['admin-additive-details', templateProductSizeId],
	queryFn: () => productsService.getProductSizeById(Number(templateProductSizeId)),
	enabled: computed(() => !isNaN(Number(templateProductSizeId))),
})

const productId = route.query.productId as string

const createMutation = useMutation({
	mutationFn: (dto: CreateProductSizeDTO) => productsService.createProductSize(dto),
	onMutate: () => {
		toast({
			title: 'Создание...',
			description: 'Создание нового размера продукта. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-product-sizes', productId] })
		queryClient.invalidateQueries({ queryKey: ['admin-product-details', productId] })
		toast({
			title: 'Успех!',
      variant: 'success',
			description: 'Размер продукта успешно создан.',
		})
		router.back()
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при создании размера продукта.',
			variant: 'destructive',
		})
	},
})

function handleCreate(data: CreateProductSizeFormSchema) {
	if (isNaN(Number(productId))) {
		toast({
			title: 'Ошибка',
			description: 'Неверный идентификатор продукта.',
			variant: 'destructive',
		})
		return router.back()
	}

  const additives: SelectedAdditiveDTO[] = data.additives.map(a => ({
    additiveId: a.additiveId,
    isDefault: a.isDefault,
    isHidden: a.isHidden,
  }))

  const ingredients: SelectedIngredientDTO[] = data.ingredients.map(a => ({
    ingredientId: a.ingredientId,
    quantity: a.quantity
  }))

  const provisions: SelectedProvisionDTO[] = data.provisions.map(a => ({
    provisionId: a.provisionId,
    volume: a.volume
  }))

	const dto: CreateProductSizeDTO = {
		productId: Number(productId),
		name: data.name,
		unitId: data.unitId,
		basePrice: data.basePrice,
		size: data.size,
    machineId: data.machineId,
		additives,
		ingredients,
    provisions
	}

	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
