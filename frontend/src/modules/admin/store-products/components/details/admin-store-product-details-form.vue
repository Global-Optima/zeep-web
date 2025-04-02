<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { Card, CardContent } from '@/core/components/ui/card'
import { Input } from '@/core/components/ui/input'
import { Label } from '@/core/components/ui/label'
import { Switch } from '@/core/components/ui/switch'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/core/components/ui/table'
import { useToast } from '@/core/components/ui/toast'
import type { StoreProductDetailsDTO, UpdateStoreProductDTO } from '@/modules/admin/store-products/models/store-products.model'
import type { ProductDetailsDTO } from '@/modules/kiosk/products/models/product.model'
import { ChevronLeft, Trash } from 'lucide-vue-next'
import { computed, onMounted, reactive, watch } from 'vue'

const props = defineProps<{
  initialStoreProduct: StoreProductDetailsDTO;
  product: ProductDetailsDTO;
  readonly?: boolean;
}>();

const emits = defineEmits<{
  onSubmit: [dto: UpdateStoreProductDTO];
  onCancel: [];
}>();

interface ProductSizeLocal {
  id: number;
  name: string;
  basePrice: number;
  storePrice: number;
  size: number;
}

interface StoreProductLocal {
  isAvailable: boolean;
  sizes: ProductSizeLocal[];
}

const storeProductLocal = reactive<StoreProductLocal>({
  isAvailable: false,
  sizes: [],
});

const { toast } = useToast();

const sortedSizes = computed(() => [...storeProductLocal.sizes].sort((a, b) => a.size - b.size));

onMounted(() => {
  initStoreProductLocal();
});

watch(
  () => [props.initialStoreProduct, props.product],
  () => {
    initStoreProductLocal();
  }
);

function initStoreProductLocal() {
  storeProductLocal.isAvailable = props.initialStoreProduct.isAvailable;

  storeProductLocal.sizes = props.initialStoreProduct.sizes.map((storeSz) => {
    return {
      id: storeSz.productSizeId,
      name: storeSz.name,
      basePrice: storeSz.basePrice,
      storePrice: storeSz.storePrice,
      size: storeSz.size
    };
  });
}

const hasSizes = computed(() => storeProductLocal.sizes.length > 0);

function removeSize(index: number) {
  if (props.readonly) return;
  storeProductLocal.sizes.splice(index, 1);
}

function onResetSizes() {
  if (props.readonly) return;
  const existingIDs = storeProductLocal.sizes.map((sz) => sz.id);

  props.product.sizes.forEach((pSz) => {
    if (!existingIDs.includes(pSz.id)) {
      storeProductLocal.sizes.push({
        id: pSz.id,
        name: pSz.name,
        basePrice: pSz.basePrice,
        storePrice: pSz.basePrice,
        size: pSz.size
      });
    }
  });
}

function onSubmit() {
  if (props.readonly) return;
  if (!storeProductLocal.sizes.length) {
    toast({ description: 'Вы должны добавить хотя бы один размер для сохранения.' });
    return;
  }

  const payload: UpdateStoreProductDTO = {
    isAvailable: storeProductLocal.isAvailable,
    productSizes: storeProductLocal.sizes.map((sz) => ({
      productSizeID: sz.id,
      storePrice: sz.storePrice,
    })),
  };
  emits('onSubmit', payload);
}

function onCancel() {
  emits('onCancel');
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

			<h1 class="font-semibold text-xl tracking-tight shrink-0">{{ props.product.name }}</h1>

			<div
				v-if="!props.readonly"
				class="hidden md:flex items-center gap-2 md:ml-auto"
			>
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
			<CardContent class="space-y-4 mt-5">
				<div class="flex flex-row justify-between items-center gap-12 p-4 border rounded-lg">
					<div class="flex flex-col space-y-1">
						<Label class="font-medium"> Продукт доступен к продаже </Label>
						<p class="text-gray-500 text-sm">Укажите доступен ли этот продукт для продажи в кафе</p>
					</div>

					<Switch
						:checked="storeProductLocal.isAvailable"
						:disabled="props.readonly"
						@update:checked="c => storeProductLocal.isAvailable = c"
					/>
				</div>

				<Table class="mt-5">
					<TableHeader>
						<TableRow>
							<TableHead>Размер</TableHead>
							<TableHead>Базовая цена</TableHead>
							<TableHead>Цена в кафе</TableHead>
							<TableHead v-if="!props.readonly">Удалить</TableHead>
						</TableRow>
					</TableHeader>
					<TableBody>
						<template v-if="hasSizes">
							<TableRow
								v-for="(size, i) in sortedSizes"
								:key="size.id"
							>
								<TableCell>{{ size.name }}</TableCell>
								<TableCell>{{ size.basePrice }}</TableCell>
								<TableCell>
									<Input
										type="number"
										class="w-24"
										v-model.number="size.storePrice"
										:readonly="props.readonly"
									/>
								</TableCell>
								<TableCell v-if="!props.readonly">
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
									class="text-gray-500 text-center"
									>Нет размеров</TableCell
								>
							</TableRow>
						</template>
					</TableBody>
				</Table>

				<div class="flex justify-center !mt-10">
					<Button
						variant="outline"
						@click="onResetSizes"
						:disabled="readonly"
						>Восстановить размеры</Button
					>
				</div>
			</CardContent>
		</Card>
	</div>
</template>
