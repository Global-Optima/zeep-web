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
					Размеры продукта
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
							<TableHead v-if="!readonly"></TableHead>
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

							<TableCell
								v-if="!readonly"
								class="text-right"
							>
								<DropdownMenu>
									<DropdownMenuTrigger asChild>
										<Button
											variant="ghost"
											size="icon"
											@click.stop
										>
											<EllipsisVertical class="w-6 h-6" />
										</Button>
									</DropdownMenuTrigger>
									<DropdownMenuContent align="end">
										<DropdownMenuItem
											@click="(e) => { e.stopPropagation(); onDuplicateClick(variant.id); }"
										>
											Дублировать
										</DropdownMenuItem>
									</DropdownMenuContent>
								</DropdownMenu>
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
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/core/components/ui/dropdown-menu'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/core/components/ui/table'
import { getRouteName } from '@/core/config/routes.config'
import { formatPrice } from '@/core/utils/price.utils'
import type { ProductDetailsDTO, ProductSizeDTO } from '@/modules/kiosk/products/models/product.model'
import { productsService } from '@/modules/kiosk/products/services/products.service'
import { useQuery } from '@tanstack/vue-query'
import { ChevronLeft, EllipsisVertical } from 'lucide-vue-next'
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

const { data: productSizes } = useQuery({
  queryKey: ['admin-product-sizes', props.productDetails.id],
  queryFn: () => productsService.getProductSizesByProductID(props.productDetails.id),
  enabled: Boolean(props.productDetails.id),
})

const sortedProductSizes = computed(() => {
  return productSizes.value ? [...productSizes.value].sort((a, b) => a.basePrice - b.basePrice) : [];
})


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

const onDuplicateClick = (id: number) => {
 router.push({name: getRouteName('ADMIN_PRODUCT_SIZE_CREATE'), query: {templateProductSizeId: id, productId: props.productDetails.id}})
}
</script>

<style scoped></style>
