<template>
	<div
		@click="onAdditiveClick(additive)"
		:class="['text-center bg-white rounded-3xl p-5 min-w-44 max-w-44 flex flex-col justify-between',
                selectedAdditives.includes(additive) ? 'bg-primary border-primary border-2' : 'border-2 border-transparent'
            ]"
	>
		<img
			:src="additive.imageUrl"
			alt="Additive Image"
			class="w-full h-24 object-contain"
		/>
		<p class="mt-3 flex-grow text-sm sm:text-base">{{ additive.name }}</p>

		<div class="flex items-center justify-between mt-5">
			<p
				v-if="additive.price"
				:class="['text-gray-400 text-lg sm:text-xl', selectedAdditives.includes(additive) ? 'text-black' : '']"
			>
				{{ formatPrice(additive.price) }}
			</p>
			<button
				class="relative h-6 w-6 sm:w-8 sm:h-8 rounded-full focus:outline-none"
				:class="selectedAdditives.includes(additive) ? 'bg-primary' : 'bg-gray-200'"
			>
				<span
					v-if="selectedAdditives.includes(additive)"
					class="absolute inset-0 m-auto sm:w-4 sm:h-4 h-3 w-3 rounded-full bg-white"
				></span>
			</button>
		</div>
	</div>
</template>

<script setup lang="ts">
import { formatPrice } from '@/core/utils/price.utils'
import type { Additives } from '@/modules/additives/models/additive.model'

  const {additive, selectedAdditives} = defineProps<{additive: Additives, selectedAdditives: Additives[]}>()
  const emit = defineEmits(["click:additive"])

  const onAdditiveClick = (additive: Additives) => {
    emit("click:additive", additive)
  }
</script>

<style scoped></style>
