<template>
	<div class="relative">
		<!-- Product Image -->
		<img
			src="https://img.pikbest.com/wp/202342/cappuccino-3d-rendering-of-a-professional-cup-with-chocolate-and-cinnamon-perfect-for-an-autumn-concept-or-background_9862878.jpg!bw700"
			alt="Product Image"
			class="w-full h-[500px] sm:h-[600px] object-cover rounded-3xl"
		/>

		<!-- Gradient Overlay -->
		<div
			class="absolute inset-0 h-[500px] sm:h-[600px] bg-gradient-to-t to-50% from-[#F5F5F7] to-transparent pointer-events-none"
		></div>

		<!-- Product Details -->
		<div class="w-full p-4 sm:p-8 sm:-mt-24 relative">
			<h1 class="text-2xl sm:text-4xl font-medium">{{ product.title }}</h1>
			<p class="text-base sm:text-xl mt-1 sm:mt-3">{{ product.description }}</p>

			<!-- Size Selection and Add to Cart -->
			<div class="flex items-center gap-4 justify-between mt-5 sm:mt-10">
				<!-- Size Options -->
				<div class="flex items-center gap-2">
					<KioskDetailsSizes
						v-for="size in sizes"
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
			<div class="mt-6 sm:mt-8">
				<p class="text-lg sm:text-2xl font-medium">Добавки</p>
				<div class="flex overflow-x-auto no-scrollbar gap-1 mt-2 sm:mt-4">
					<KioskDetailsAdditives
						v-for="additive in additives"
						:key="additive.id"
						:selected-additives="selectedAdditives"
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
import type { Additives } from '@/modules/additives/models/additive.model'
import { useCartStore } from '@/modules/kiosk/cart/stores/cart.store'
import KioskDetailsAdditives from '@/modules/kiosk/products/components/details/kiosk-details-additives.vue'
import KioskDetailsEnergy from '@/modules/kiosk/products/components/details/kiosk-details-energy.vue'
import KioskDetailsSizes from '@/modules/kiosk/products/components/details/kiosk-details-sizes.vue'
import type { Products, ProductSizes } from '@/modules/products/models/product.model'
import { Icon } from '@iconify/vue'
import { computed, ref, watch } from 'vue'

// Define props
const props = defineProps<{ product: Products }>()
const emit = defineEmits({"close": null })

const cartStore = useCartStore()

// Available Sizes (if dynamic, otherwise pass via props or fetch)
const sizes = ref<ProductSizes[]>([
  { id: 1, label: 'S', volume: 250, price: 1000 },
  { id: 2, label: 'M', volume: 350, price: 1100 },
  { id: 3, label: 'L', volume: 450, price: 1200 },
])

// Available Additives (if dynamic, otherwise pass via props or fetch)
const additives = ref<Additives[]>([
  {
    id: 1,
    name: 'Сырная пенка',
    price: 400,
    imageUrl: 'https://example.com/images/additive1.png',
  },
  {
    id: 2,
    name: 'Шоколадная крошка',
    price: 300,
    imageUrl: 'https://example.com/images/additive2.png',
  },
  {
    id: 3,
    name: 'Взбитые сливки',
    price: 350,
    imageUrl: 'https://example.com/images/additive3.png',
  },
])

// Selected Size (default to first size)
const selectedSize = ref<ProductSizes>(sizes.value[0])

// Selected Additives
const selectedAdditives = ref<Additives[]>([])

// Quantity (default to 1)
const quantity = ref<number>(1)

// Computed Total Price
const totalPrice = computed<number>(() => {
  const basePrice = selectedSize.value.price
  const additivesPrice = selectedAdditives.value.reduce((sum, additive) => sum + additive.price, 0)
  return (basePrice + additivesPrice) * quantity.value
})

const calculatedEnergy = computed(() => {
  return {
    ccal: 400,
    proteins: 200,
    carbs: 120,
    fats: 30,
  };
});

// Handlers
const onSizeClick = (newSize: ProductSizes) => {
  selectedSize.value = newSize
}

const onAdditiveClick = (additive: Additives) => {
  const index = selectedAdditives.value.findIndex((item) => item.id === additive.id)
  if (index !== -1) {
    selectedAdditives.value.splice(index, 1)
  } else {
    selectedAdditives.value.push(additive)
  }
}


const handleAddToCart = () => {
  cartStore.addToCart(props.product, selectedSize.value, selectedAdditives.value, quantity.value)
  toastSuccess('Добавлено в корзину');
}

// Watch for product prop changes to reset selections
watch(
  () => props.product,
  () => {
    selectedSize.value = sizes.value[0]
    selectedAdditives.value = []
    quantity.value = 1
  },
  { immediate: true }
)
</script>

<style scoped>
/* Add your styles here */
</style>
