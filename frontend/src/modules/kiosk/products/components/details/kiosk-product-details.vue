<template>
	<SwipeModal
		v-model="isOpen"
		snapPoint="95dvh"
		class="relative"
	>
		<div class="relative">
			<button
				@click="closeModal"
				class="absolute top-0 right-6 sm:right-8 rounded-full p-2 sm:p-3 bg-gray-600 z-10"
			>
				<X class="w-6 h-6 sm:w-8 sm:h-8 text-white" />
			</button>

			<!-- Loading Skeletons -->
			<SkeletonLoader v-if="isLoading || isAdditivesLoading || !productDetails" />

			<!-- Error Handling -->
			<ErrorDisplay
				v-else-if="isError || isAdditivesError"
				@retry="retryFetch"
			/>

			<!-- Product Content -->
			<template v-else>
				<div class="relative">
					<img
						:src="productDetails.imageUrl"
						:alt="productDetails.name"
						class="product-image"
					/>
					<div class="overlay"></div>

					<div class="content-container">
						<h1 class="text-2xl sm:text-4xl font-medium">{{ productDetails.name }}</h1>
						<p class="text-base sm:text-xl mt-1 sm:mt-3">{{ productDetails.description }}</p>

						<!-- Size Selection -->
						<div
							class="flex flex-col sm:flex-row items-start sm:items-center gap-4 justify-between mt-5 sm:mt-10"
						>
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
								<p class="text-2xl sm:text-3xl font-medium">{{ formatPrice(totalPrice) }}</p>
								<button
									@click="handleAddToCart"
									class="rounded-full bg-primary text-primary-foreground p-2 sm:p-4"
									:aria-label="'Add ' + productDetails.name + ' to cart'"
								>
									<Icon
										icon="mingcute:add-line"
										class="text-2xl sm:text-3xl"
									/>
								</button>
							</div>
						</div>

						<!-- Energy Info -->
						<div class="mt-6 sm:mt-8">
							<KioskDetailsEnergy :energy="calculatedEnergy" />
						</div>

						<!-- Additives Selection -->
						<div
							v-for="category in sortedAdditiveCategories"
							:key="category.id"
							class="mt-6 sm:mt-8"
						>
							<p class="text-lg sm:text-2xl font-medium">{{ category.name }}</p>
							<div class="flex overflow-x-auto no-scrollbar gap-1 mt-2 sm:mt-4">
								<KioskDetailsAdditives
									v-for="additive in category.additives"
									:key="additive.id"
									:default-additives="productDetails.defaultAdditives"
									:additive="additive"
									:selected-additives="selectedAdditives"
									@click:additive="onAdditiveClick"
								/>
							</div>
						</div>
					</div>
				</div>
			</template>
		</div>
	</SwipeModal>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue';
import { SwipeModal } from '@takuma-ru/vue-swipe-modal';
import { computed, ref, watch } from 'vue';
import { useQuery } from '@tanstack/vue-query';
import { formatPrice } from '@/core/utils/price.utils';
import { useCartStore } from '@/modules/kiosk/cart/stores/cart.store';
import KioskDetailsEnergy from '@/modules/kiosk/products/components/details/kiosk-details-energy.vue';
import KioskDetailsSizes from '@/modules/kiosk/products/components/details/kiosk-details-sizes.vue';
import type { Additive, ProductSize } from '@/modules/kiosk/products/models/product.model';
import { productService } from '@/modules/kiosk/products/services/products.service';
import { useProductStore } from '@/modules/kiosk/products/stores/current-product.store';
import { storeToRefs } from "pinia";
import { X } from "lucide-vue-next";
import KioskDetailsAdditives from "@/modules/kiosk/products/components/details/kiosk-details-additives.vue";

const isOpen = ref(false);
const ss = useProductStore();
const { selectedProductId } = storeToRefs(ss);
const cartStore = useCartStore();

const selectedSize = ref<ProductSize | null>(null);
const selectedAdditives = ref<Additive[]>([]);
const quantity = ref(1);

// Fetch Product Details
const { data: productDetails, isLoading, isError, refetch: refetchProductDetails } = useQuery({
	queryKey: ['product-details', selectedProductId],
	queryFn: () => productService.getStoreProductDetails(selectedProductId.value!),
	enabled: computed(() => Boolean(selectedProductId.value)),
});

// Fetch Additive Categories
const { data: additiveCategories, isLoading: isAdditivesLoading, isError: isAdditivesError, refetch: refetchAdditiveCategories } = useQuery({
	queryKey: ['additive-categories', selectedProductId],
	queryFn: () => productService.getAdditiveCategoriesByStoreAndProduct(selectedProductId.value!),
	enabled: computed(() => Boolean(selectedProductId.value)),
});

const totalPrice = computed(() => {
	if (!selectedSize.value) return 0;

	const basePrice = selectedSize.value.basePrice;
	const additivesPrice = selectedAdditives.value.reduce((sum, additive) => sum + additive.price, 0);
	return (basePrice + additivesPrice) * quantity.value;
});

const calculatedEnergy = computed(() => {
	if (!productDetails.value) return { ccal: 0, proteins: 0, carbs: 0, fats: 0 };
	return { ccal: 400, proteins: 20, carbs: 13, fats: 10 };
});

watch(
  () => productDetails.value,
  (newProductDetails) => {
    if (newProductDetails) {
      selectedSize.value = newProductDetails.sizes[0]
      selectedAdditives.value = [];
      quantity.value = 1;
    }
  },
  {flush: "sync"}
);

// Ensure default additives appear first
const sortedAdditiveCategories = computed(() => {
	return additiveCategories.value?.map(category => ({
		...category,
		additives: [
			...category.additives.filter(additive => productDetails.value?.defaultAdditives.some(defaultAdd => defaultAdd.id === additive.id)),
			...category.additives.filter(additive => !productDetails.value?.defaultAdditives.some(defaultAdd => defaultAdd.id === additive.id))
		]
	}));
});

watch(
	() => selectedProductId.value,
	(newProductId, prevProductId) => {
	if (newProductId && newProductId !== prevProductId) {
		isOpen.value = true;
	}
	}
 );

watch(isOpen, (open) => {
	if (!open) resetSelections();
});

const resetSelections = () => {
	selectedSize.value = null;
	selectedAdditives.value = [];
	quantity.value = 1;

	console.log("Here")
};

const onSizeClick = (newSize: ProductSize) => selectedSize.value = newSize;

const onAdditiveClick = (additive: Additive) => {
	const index = selectedAdditives.value.findIndex(item => item.id === additive.id);
	if (index !== -1) {
		selectedAdditives.value.splice(index, 1)
	} else {
		selectedAdditives.value.push(additive)
	}
};

const handleAddToCart = () => {
	if (!productDetails.value || !selectedSize.value) return;
	cartStore.addToCart(productDetails.value, selectedSize.value, selectedAdditives.value, quantity.value);
	isOpen.value = false;
};

const closeModal = () => {
	isOpen.value = false;
	ss.closeBottomSheet();
};

const retryFetch = () => {
	if (isError) refetchProductDetails();
	if (isAdditivesError) refetchAdditiveCategories();
};
</script>

<style scoped>
.product-image {
	width: 100%;
	height: 500px;
	object-fit: cover;
}

.overlay {
	position: absolute;
	inset: 0;
	height: 500px;
	background: linear-gradient(to top, #F5F5F7 5%, transparent);
}


.content-container {
	width: 100%;
	padding: 4rem;
	position: relative;
}

.modal-style {
  background-color: #F5F5F7 !important;
  color:#000 !important;
}
</style>
