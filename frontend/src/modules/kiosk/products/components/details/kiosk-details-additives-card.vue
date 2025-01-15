<template>
	<div
		@click="handleClick"
		:class="[
      'text-center bg-white rounded-3xl p-5 min-w-44 max-w-44 flex flex-col justify-between border-2',
      isSelected ? 'bg-primary border-primary' : 'border-transparent',
      isDefault  ? 'cursor-not-allowed opacity-60 !border-primary' : 'cursor-pointer',

    ]"
		data-testid="additive-card"
	>
		<img
			:src="additive.imageUrl"
			alt="Additive Image"
			class="w-full h-24 object-contain"
			data-testid="additive-image"
		/>
		<p
			class="flex-grow mt-3 text-sm sm:text-base"
			data-testid="additive-name"
		>
			{{ additive.name }}
		</p>

		<div class="flex justify-between items-center mt-5">
			<p
				:class="[
          'text-gray-400 text-lg sm:text-xl',
          isSelected ? 'text-black' : '',
        ]"
				data-testid="additive-price"
			>
				{{ isDefault ? formatPrice(0) : formatPrice(additive.price) }}
			</p>
			<button
				class="relative rounded-full w-6 sm:w-8 h-6 sm:h-8 focus:outline-none"
				:class="[
          isSelected ? 'bg-primary' : 'bg-gray-200',
          isDefault ? 'cursor-not-allowed opacity-50' : '',
        ]"
				:disabled="isDefault"
				data-testid="additive-button"
			>
				<span
					v-if="isSelected"
					class="absolute inset-0 bg-white m-auto rounded-full w-3 sm:w-4 h-3 sm:h-4"
					data-testid="additive-selected-indicator"
				></span>
			</button>
		</div>
	</div>
</template>

<script setup lang="ts">
import { formatPrice } from '@/core/utils/price.utils'
import type { AdditiveCategoryItemDTO } from '@/modules/admin/additives/models/additives.model'

const {additive, isSelected, isDefault} = defineProps<{
  additive: AdditiveCategoryItemDTO;
  isSelected: boolean;
  isDefault: boolean;
}>()

const emit = defineEmits(['click:additive'])

const handleClick = () => {
  if (isDefault) {
    return
  }
  emit('click:additive', additive)
}
</script>

<style scoped></style>
