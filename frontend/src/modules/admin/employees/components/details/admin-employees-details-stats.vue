<template>
	<div class="flex items-center bg-white p-6 border rounded-xl">
		<div class="bg-green-100 p-3 rounded-full">
			<component
				:is="icon"
				class="w-6 h-6 text-primary"
			/>
		</div>
		<div class="ml-4">
			<p class="text-gray-500 text-sm">{{ title }}</p>
			<p class="font-semibold text-gray-900 text-xl">{{ value }}</p>
		</div>
	</div>
</template>

<script setup lang="ts">
import type { LucideIcon } from 'lucide-vue-next'


defineProps<{
  title: string;
  value: string;
  icon: LucideIcon;
}>();
</script>
