<script setup lang="ts">
import { useToast } from '@/core/components/ui/toast'
import { computed, onMounted, reactive, watch } from 'vue'

import { Button } from '@/core/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/core/components/ui/card'
import { Input } from '@/core/components/ui/input'
import { Switch } from '@/core/components/ui/switch'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/core/components/ui/table'

import { Label } from '@/core/components/ui/label'
import type { StoreProductDetailsDTO, UpdateStoreProductDTO } from '@/modules/admin/store-products/models/store-products.model'
import type { ProductDetailsDTO } from '@/modules/kiosk/products/models/product.model'
import { ChevronLeft, Trash } from 'lucide-vue-next'

const props = defineProps<{
  initialStoreProduct: StoreProductDetailsDTO
  product: ProductDetailsDTO
}>()

const emits = defineEmits<{
  onSubmit: [dto: UpdateStoreProductDTO]
  onCancel: []
}>()

interface ProductSizeLocal {
  id: number
  name: string
  basePrice: number
  storePrice: number
}

interface StoreProductLocal {
  isAvailable: boolean
  sizes: ProductSizeLocal[]
}

const storeProductLocal = reactive<StoreProductLocal>({
  isAvailable: false,
  sizes: [],
})

const {toast} = useToast()

onMounted(() => {
  initStoreProductLocal()
})

watch(
  () => [props.initialStoreProduct, props.product],
  () => {
    initStoreProductLocal()
  }
)

function initStoreProductLocal() {
  storeProductLocal.isAvailable = props.initialStoreProduct.isAvailable

  storeProductLocal.sizes = props.initialStoreProduct.sizes.map((storeSz) => {
    const productSize = props.product.sizes.find((pSz) => pSz.id === storeSz.id)
    const basePrice = productSize ? productSize.basePrice : storeSz.storePrice

    return {
      id: storeSz.productSizeId,
      name: storeSz.name,
      basePrice,
      storePrice: storeSz.storePrice,
    }
  })
}

const hasSizes = computed(() => storeProductLocal.sizes.length > 0)

function removeSize(index: number) {
  storeProductLocal.sizes.splice(index, 1)
}

function onResetSizes() {
  const existingIDs = storeProductLocal.sizes.map((sz) => sz.id)

  props.product.sizes.forEach((pSz) => {
    if (!existingIDs.includes(pSz.id)) {
      storeProductLocal.sizes.push({
        id: pSz.id,
        name: pSz.name,
        basePrice: pSz.basePrice,
        storePrice: pSz.basePrice,
      })
    }
  })
}

function onSubmit() {
  if (!storeProductLocal.sizes.length) {
    toast({description:'Вы должны добавить хотя бы один размер для сохранения.'})
    return
  }

  const payload: UpdateStoreProductDTO = {
    isAvailable: storeProductLocal.isAvailable,
    productSizes: storeProductLocal.sizes.map((sz) => ({
      productSizeID: sz.id,
      storePrice: sz.storePrice,
    })),
  }
  emits('onSubmit', payload)
}

function onCancel() {
  emits('onCancel')
}
</script>

<template>
	<div class="flex-1 gap-4 grid auto-rows-max mx-auto max-w-4xl">
		<div class="flex items-center gap-4">
			<Button
				variant="outline"
				size="icon"
				@click="onCancel"
			>
				<ChevronLeft class="w-5 h-5" />
				<span class="sr-only">Назад</span>
			</Button>

			<h1 class="font-semibold text-xl tracking-tight shrink-0">Обновить товар</h1>

			<div class="md:flex items-center gap-2 hidden md:ml-auto">
				<Button
					variant="outline"
					@click="onCancel"
					>Отменить</Button
				>
				<Button
					type="submit"
					@click="onSubmit"
					>Сохранить</Button
				>
			</div>
		</div>

		<Card>
			<CardHeader>
				<div class="flex justify-between">
					<div>
						<CardTitle>{{ props.product.name }}</CardTitle>
						<CardDescription class="mt-1">{{ props.product.category.name }}</CardDescription>
					</div>
					<Button
						variant="outline"
						@click="onResetSizes"
						>Сбросить размеры к исходным</Button
					>
				</div>
			</CardHeader>

			<CardContent class="space-y-4">
				<div class="flex flex-row justify-between items-center gap-12 p-4 border rounded-lg">
					<div class="flex flex-col space-y-0.5">
						<Label class="font-medium text-base"> Товар доступен к продаже </Label>
						<p class="text-gray-500 text-sm">
							Укажите доступен ли этот товар для продажи в магазине
						</p>
					</div>

					<Switch
						:checked="storeProductLocal.isAvailable"
						@update:checked="c => storeProductLocal.isAvailable = c"
					/>
				</div>

				<Table class="mt-5">
					<TableHeader>
						<TableRow>
							<TableHead>Размер</TableHead>
							<TableHead>Базовая цена</TableHead>
							<TableHead>Цена в магазине</TableHead>
							<TableHead>Удалить</TableHead>
						</TableRow>
					</TableHeader>
					<TableBody>
						<template v-if="hasSizes">
							<TableRow
								v-for="(size, i) in storeProductLocal.sizes"
								:key="size.id"
							>
								<TableCell>{{ size.name }}</TableCell>
								<TableCell>{{ size.basePrice }}</TableCell>
								<TableCell>
									<Input
										type="number"
										class="w-24"
										v-model.number="size.storePrice"
									/>
								</TableCell>
								<TableCell>
									<Button
										variant="ghost"
										size="icon"
										@click="removeSize(i)"
									>
										<Trash class="w-5 h-5 text-red-500 hover:text-red-700" />
									</Button>
								</TableCell>
							</TableRow>
						</template>
						<template v-else>
							<TableRow>
								<TableCell
									colspan="4"
									class="text-center text-gray-500"
								>
									Нет размеров
								</TableCell>
							</TableRow>
						</template>
					</TableBody>
				</Table>
			</CardContent>
		</Card>

		<div class="flex justify-center items-center gap-2 md:hidden mt-4">
			<Button
				variant="outline"
				@click="onCancel"
				>Отменить</Button
			>
			<Button
				type="submit"
				@click="onSubmit"
				>Сохранить</Button
			>
		</div>
	</div>
</template>
