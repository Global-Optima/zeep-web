<!-- src/components/ProductImage.vue -->
<template>
	<div class="relative pb-52 overflow-y-auto no-scrollbar">
		<img
			:src="imageUrl"
			:alt="altText"
			class="w-full h-[500px] object-cover"
		/>
		<div class="absolute inset-0 bg-gradient-to-t from-gray-100 via-transparent h-[500px]"></div>
	</div>
</template>

<script setup lang="ts">
import { defineProps } from 'vue'

const {imageUrl, altText} = defineProps<{
  imageUrl: string
  altText: string
}>()
</script>

<style scoped>
/* Add any specific styles if needed */
</style>
