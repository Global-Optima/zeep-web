<template>
	<!-- Product  -->
	<div v-if="productDetails">
		<!-- Product Image -->
		<img
			:src="productDetails.imageUrl"
			alt="Product Image"
			class="w-full h-[500px] sm:h-[600px] object-cover rounded-3xl"
		/>

		<!-- Gradient Overlay -->
		<div
			class="absolute inset-0 h-[500px] sm:h-[600px] bg-gradient-to-t to-50% from-[#F5F5F7] to-transparent pointer-events-none"
		></div>

		<!-- Product Details -->
		<div class="w-full p-4 sm:p-8 sm:-mt-24 relative">
			<h1 class="text-2xl sm:text-4xl font-medium">{{ productDetails.name }}</h1>
			<p class="text-base sm:text-xl mt-1 sm:mt-3">{{ productDetails.description }}</p>

			<!-- Size Selection and Add to Cart -->
			<div class="flex items-center gap-4 justify-between mt-5 sm:mt-10">
				<!-- Size Options -->
				<div class="flex items-center gap-2">
					<KioskDetailsSizes
						v-for="size in productDetails.sizes"
						:key="size.id"
						:size="size"
						:selected-size="selectedSize"
						@click:size="onSizeClick"
					/>
				</div>

				<!-- Price and Add to Cart Button -->
				<div class="flex items-center gap-4 sm:gap-6">
					<p class="text-2xl sm:text-3xl font-medium">
						{{ formatPrice(totalPrice) }}
					</p>
					<button
						@click="handleAddToCart"
						class="rounded-full bg-primary text-primary-foreground p-2 sm:p-4"
					>
						<Icon
							icon="mingcute:add-line"
							class="text-2xl sm:text-3xl"
						/>
					</button>
				</div>
			</div>

			<div class="mt-6 sm:mt-8">
				<KioskDetailsEnergy :energy="calculatedEnergy" />
			</div>

			<!-- Additives Selection -->
			<div
				class="mt-6 sm:mt-8"
				v-for="additiveCategory in additiveCategories"
				:key="additiveCategory.id"
			>
				<p class="text-lg sm:text-2xl font-medium">{{ additiveCategory.name }}</p>
				<div class="flex overflow-x-auto no-scrollbar gap-1 mt-2 sm:mt-4">
					<KioskDetailsAdditives
						v-for="additive in additiveCategory.additives"
						:key="additive.id"
						:selected-additives="selectedAdditives"
						:default-additives="productDetails.defaultAdditives"
						:additive="additive"
						@click:additive="onAdditiveClick"
					/>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { toastSuccess } from '@/core/config/toast.config'
import { formatPrice } from '@/core/utils/price.utils'
import { useCartStore } from '@/modules/kiosk/cart/stores/cart.store'
import KioskDetailsAdditives from '@/modules/kiosk/products/components/details/kiosk-details-additives.vue'
import KioskDetailsEnergy from '@/modules/kiosk/products/components/details/kiosk-details-energy.vue'
import KioskDetailsSizes from '@/modules/kiosk/products/components/details/kiosk-details-sizes.vue'
import type { Additive, ProductSize } from '@/modules/kiosk/products/models/product.model'
import { productService } from '@/modules/kiosk/products/services/products.service'
import { Icon } from '@iconify/vue'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref, watch } from 'vue'

const {selectedProductId} = defineProps<{ selectedProductId: number | null }>();
const cartStore = useCartStore();


const { data: productDetails, isLoading, isError } = useQuery({
  queryKey:['product-details', selectedProductId],
  queryFn: () => productService.getStoreProductDetails(selectedProductId!),
  enabled: Boolean(selectedProductId),
});

const { data: additiveCategories } = useQuery({
  queryKey:['additive-categories', selectedProductId],
  queryFn: () => productService.getAdditiveCategoriesByStoreAndProduct(selectedProductId!),
  enabled: Boolean(selectedProductId),
});

const selectedSize = ref<ProductSize | null>(null);
const selectedAdditives = ref<Additive[]>([]);
const quantity = ref<number>(1);

const totalPrice = computed<number>(() => {
  if (!selectedSize.value) return 0;
  const basePrice = selectedSize.value.basePrice;
  const additivesPrice = selectedAdditives.value.reduce((sum, additive) => sum + additive.price, 0);
  return (basePrice + additivesPrice) * quantity.value;
});

const calculatedEnergy = computed(() => {
  if (!productDetails.value) return { ccal: 0, proteins: 0, carbs: 0, fats: 0 };
  return {
    ccal: 400,
    proteins: 200,
    carbs: 120,
    fats: 30,
  };
});

watch(
  () => productDetails.value,
  (newProductDetails) => {
    if (newProductDetails) {
      selectedSize.value = newProductDetails.sizes && newProductDetails.sizes[0] ? newProductDetails.sizes[0] : null;
      selectedAdditives.value = [];
      quantity.value = 1;
    }
  },
  { immediate: true },
);

const onSizeClick = (newSize: ProductSize) => {
  selectedSize.value = newSize;
};

const onAdditiveClick = (additive: Additive) => {
  const index = selectedAdditives.value.findIndex((item) => item.id === additive.id);
  if (index !== -1) {
    selectedAdditives.value.splice(index, 1);
  } else {
    selectedAdditives.value.push(additive);
  }
};

const handleAddToCart = () => {
  if (!productDetails.value || !selectedSize.value) return

  cartStore.addToCart(
    productDetails.value,
    selectedSize.value,
    selectedAdditives.value,
    quantity.value,
  );
  toastSuccess('Добавлено в корзину');
};
</script>

<style scoped></style>
