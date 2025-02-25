<!-- src/components/ProductImage.vue -->
<template>
	<div class="relative overflow-y-auto no-scrollbar">
		<LazyImage
			:src="imageUrl"
			:alt="altText"
			:class="cn('w-full object-cover h-[500px]', imageClass)"
		/>
		<div
			:class="cn('absolute inset-0 bg-gradient-to-t from-gray-100 via-transparent h-[500px]', overlayClass)"
		></div>
	</div>
</template>

<script setup lang="ts">
import LazyImage from '@/core/components/lazy-image/LazyImage.vue'
import { cn } from '@/core/utils/tailwind.utils'

const {imageUrl, altText, imageClass, overlayClass} = defineProps<{
  imageUrl: string
  altText: string,
  imageClass?: string,
  overlayClass?: string
}>()
</script>

<style scoped>
/* Add any specific styles if needed */
</style>
