<template>
	<div class="flex-1 items-start gap-4 grid mx-auto max-w-6xl">
		<div class="flex justify-between items-center gap-2">
			<div class="flex justify-between items-center gap-4">
				<Button
					variant="outline"
					size="icon"
					type="button"
					@click="onCancel"
				>
					<ChevronLeft class="w-5 h-5" />
					<span class="sr-only">Назад</span>
				</Button>
				<!-- Title -->
				<h1
					class="flex-1 sm:grow-0 font-semibold text-xl tracking-tight whitespace-nowrap shrink-0"
				>
					Размеры товара
				</h1>
			</div>
			<template v-if="!readonly">
				<Button @click="onAddNewVariantClick">
					<span class="sr-only sm:not-sr-only sm:whitespace-nowrap">Добавить</span>
				</Button>
			</template>
		</div>

		<Card>
			<CardContent class="mt-4">
				<Table>
					<TableHeader>
						<TableRow>
							<TableHead>Название</TableHead>
							<TableHead>Размер</TableHead>
							<TableHead>Начальная цена</TableHead>
							<TableHead>По умолчанию</TableHead>
						</TableRow>
					</TableHeader>
					<TableBody>
						<TableRow
							v-for="variant in sortedProductSizes"
							:key="variant.id"
							@click="onVariantClick(variant)"
							class="hover:bg-slate-50 cursor-pointer"
						>
							<TableCell class="py-4 font-medium">{{ variant.name }}</TableCell>
							<TableCell>{{ variant.size }} {{ variant.unit.name }}</TableCell>
							<TableCell>
								{{ formatPrice(variant.basePrice) }}
							</TableCell>
							<TableCell class="text-left">
								<input
									type="radio"
									name="defaultVariant"
									:value="variant.id"
									:checked="variant.isDefault"
									class="size-5"
									:disabled="readonly"
									@click.stop="setDefaultVariant(variant)"
								/>
							</TableCell>
						</TableRow>
					</TableBody>
				</Table>
			</CardContent>
		</Card>
	</div>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { Card, CardContent } from '@/core/components/ui/card'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/core/components/ui/table'
import { useToast } from '@/core/components/ui/toast'
import { formatPrice } from '@/core/utils/price.utils'
import type { ProductDetailsDTO, ProductSizeDTO, UpdateProductSizeDTO } from '@/modules/kiosk/products/models/product.model'
import { productsService } from '@/modules/kiosk/products/services/products.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { ChevronLeft } from 'lucide-vue-next'
import { computed } from 'vue'
import { useRouter } from "vue-router"

const props = defineProps<{
  productDetails: ProductDetailsDTO
  readonly?: boolean
}>()

const emits = defineEmits<{
  onCancel: []
}>()

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const { data: productSizes } = useQuery({
  queryKey: ['admin-product-sizes', props.productDetails.id],
  queryFn: () => productsService.getProductSizesByProductID(props.productDetails.id),
  enabled: Boolean(props.productDetails.id),
})

const sortedProductSizes = computed(() => {
  return productSizes.value ? [...productSizes.value].sort((a, b) => a.basePrice - b.basePrice) : [];
})

const { mutate: updateSize } = useMutation({
  mutationFn: (props: { productSizeId: number, dto: UpdateProductSizeDTO }) =>
    productsService.updateProductSize(props.productSizeId, props.dto),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['admin-product-sizes', props.productDetails.id] })
    toast({
      title: 'Успех!',
      description: 'Размер товара обновлен.',
    })
  },
  onError: () => {
    toast({
      title: 'Ошибка',
      description: 'Ошибка при обновлении размера товара.',
      variant: 'destructive',
    })
  }
})

const setDefaultVariant = (size: ProductSizeDTO) => {
  if (!props.readonly) {
    if (!size.isDefault) {
      updateSize({ productSizeId: size.id, dto: { isDefault: true } })
    }
  }
}

const onVariantClick = (variant: ProductSizeDTO) => {
  router.push(`../product-sizes/${variant.id}?productId=${props.productDetails.id}`)
}

function onCancel() {
  emits('onCancel')
}

const onAddNewVariantClick = () => {
  if (!props.readonly) {
    router.push(`../product-sizes/create?productId=${props.productDetails.id}`)
  }
}
</script>

<style scoped></style>
