<template>
	<div class="flex gap-4 sm:gap-8 p-4 sm:p-6 rounded-3xl bg-white">
		<img
			:src="item.product.image"
			alt="Product Image"
			class="w-16 h-16 sm:w-fit sm:h-20 object-contain rounded"
		/>

		<div class="flex-1">
			<p class="text-base sm:text-2xl">{{ item.product.title }}</p>

			<div class="flex justify-between items-center mt-2">
				<p class="text-lg sm:text-2xl font-medium">{{ formatPrice(item.product.price) }}</p>

				<div class="flex items-center gap-2">
					<button
						@click="decrement(item.product.id)"
						class="p-1 sm:p-2 bg-gray-200 rounded-xl"
					>
						<Icon
							icon="mingcute:minimize-line"
							class="text-lg sm:text-xl"
						/>
					</button>

					<span class="text-base sm:text-2xl mx-1 sm:mx-2">{{ item.quantity }}</span>

					<button
						@click="increment(item.product.id)"
						class="p-1 sm:p-2 bg-gray-200 rounded-xl"
					>
						<Icon
							icon="mingcute:add-line"
							class="text-lg sm:text-xl"
						/>
					</button>
				</div>
			</div>

			<div class="mt-3">
				<div
					v-if="item.additives.length > 0"
					class="flex flex-col gap-1 sm:gap-2 text-gray-600"
				>
					<div
						v-for="additive in item.additives"
						:key="additive.id"
						class="inline-flex items-center gap-1 sm:gap-2"
					>
						<Icon
							icon="mingcute:add-line"
							class="text-lg sm:text-xl text-primary"
						/>
						<p class="text-xs sm:text-lg">
							{{ additive.name }}
						</p>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { formatPrice } from '@/core/utils/price.utils'
import type { CartItems } from '@/modules/kiosk/models/cart.models'
import { Icon } from '@iconify/vue'

const {item} = defineProps<{item: CartItems}>()

const removeItem = (id: number) => {
  console.log(id)
}

const increment = (id: number) => {
  console.log(id)

}

const decrement = (id: number) => {
  console.log(id)
  removeItem(id)

}
</script>

<style scoped></style>
