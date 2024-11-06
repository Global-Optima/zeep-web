<template>
	<div class="relative">
		<button
			@click="goBack"
			class="fixed left-4 top-20 rounded-full bg-slate-800/70 text-white backdrop-blur-md p-2 sm:p-4 z-10"
		>
			<Icon
				icon="mingcute:left-line"
				class="text-3xl"
			/>
		</button>

		<img
			:src="imageUrl"
			alt="Product Image"
			class="w-full h-[500px] sm:h-[600px] object-cover rounded-3xl"
		/>

		<div
			class="absolute inset-0 h-[500px] sm:h-[600px] bg-gradient-to-t to-50% from-[#F5F5F7] to-transparent  pointer-events-none"
		></div>

		<div class="w-full p-4 sm:p-8  sm:-mt-24 relative">
			<h1 class="text-2xl sm:text-4xl font-medium">Какао</h1>
			<p class="tex-base sm:text-xl mt-1 sm:mt-3">Классический какао с молоком на выбор</p>

			<div class="flex items-center  gap-4 justify-between mt-5 sm:mt-10">
				<div class="flex items-center gap-2">
					<KioskDetailsSizes
						v-for="size in sizes"
						:key="size.id"
						:size="size"
						:selected-size="selectedSize"
						@click:size="onSizeClick"
					/>
				</div>

				<div class="flex items-center gap-4 sm:gap-6">
					<p class="text-2xl sm:text-3xl font-medium">{{ formatPrice(999) }}</p>
					<button class="rounded-full bg-primary text-primary-foreground p-2 sm:p-4">
						<Icon
							icon="mingcute:add-line"
							class="text-2xl sm:text-3xl"
						/>
					</button>
				</div>
			</div>

			<div class="mt-6 sm:mt-10">
				<p class="text-lg sm:text-2xl">Добавки</p>
				<div class="flex overflow-x-auto no-scrollbar gap-1 mt-2 sm:mt-6">
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
import { getRouteName } from '@/core/config/routes.config'
import { formatPrice } from '@/core/utils/price.utils'
import type { Additives } from '@/modules/additives/models/additive.model'
import KioskDetailsAdditives from '@/modules/kiosk/components/details/kiosk-details-additives.vue'
import KioskDetailsSizes from '@/modules/kiosk/components/details/kiosk-details-sizes.vue'
import type { ProductSizes } from '@/modules/products/models/product.model'
import { Icon } from '@iconify/vue'
import { ref, type Ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const imageUrl = "https://media.istockphoto.com/id/1409579518/photo/perfect-iced-cocoa-aka-iced-cacao-stills-with-ice-cubes.jpg?s=612x612&w=0&k=20&c=ZUCCIGQaKSRWd_sPv27wT3msZR0yHwN7ot72jAKjbi8="

const sizes: Ref<ProductSizes[]> = ref([
	{ id: 1, label: 'S', volume: 250 },
	{ id: 2, label: 'M', volume: 350 },
	{ id: 3, label: 'L', volume: 450 },
])
const selectedSize = ref(sizes.value[0])

const additives: Ref<Additives[]> = ref([
	{ id: 1, name: 'Сырная пенка', price: 400, imageUrl: 'https://png.pngtree.com/png-vector/20230906/ourmid/pngtree-top-view-of-coffee-foam-art-on-cup-png-image_10008794.png' },
	{ id: 2, name: 'Шоколадная крошка', price: 300, imageUrl: 'https://png.pngtree.com/png-vector/20230906/ourmid/pngtree-top-view-of-coffee-foam-art-on-cup-png-image_10008794.png' },
	{ id: 3, name: 'Взбитые сливки', price: 350, imageUrl: 'https://png.pngtree.com/png-vector/20230906/ourmid/pngtree-top-view-of-coffee-foam-art-on-cup-png-image_10008794.png' },
  { id: 1, name: 'Сырная пенка', price: 400, imageUrl: 'https://png.pngtree.com/png-vector/20230906/ourmid/pngtree-top-view-of-coffee-foam-art-on-cup-png-image_10008794.png' },
	{ id: 2, name: 'Шоколадная крошка', price: 300, imageUrl: 'https://png.pngtree.com/png-vector/20230906/ourmid/pngtree-top-view-of-coffee-foam-art-on-cup-png-image_10008794.png' },
	{ id: 3, name: 'Взбитые сливки', price: 350, imageUrl: 'https://png.pngtree.com/png-vector/20230906/ourmid/pngtree-top-view-of-coffee-foam-art-on-cup-png-image_10008794.png' }
])
const selectedAdditives = ref([] as Additives[])

const onSizeClick = (newSize: ProductSizes) => {
  selectedSize.value = newSize
}

const onAdditiveClick = (additive: Additives) => {
	if (selectedAdditives.value.includes(additive)) {
		selectedAdditives.value = selectedAdditives.value.filter(item => item !== additive)
	} else {
		selectedAdditives.value.push(additive)
	}
}

function goBack() {
	router.push({name: getRouteName("KIOSK_HOME")})
}
</script>

<style lang="scss"></style>
