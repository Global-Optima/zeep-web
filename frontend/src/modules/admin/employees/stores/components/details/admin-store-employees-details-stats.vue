<template>
	<div class="flex items-center bg-white p-6 border rounded-xl">
		<div class="bg-green-100 p-3 rounded-full">
			<component
				:is="iconComponent"
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
import { CheckCircle, Clock, DollarSign } from 'lucide-vue-next'
import { computed } from 'vue'

const props = defineProps<{
  title: string;
  value: string;
  icon: string;
}>();

const iconComponent = computed(() => {
  switch (props.icon) {
    case 'dollar-sign':
      return DollarSign;
    case 'clock':
      return Clock;
    case 'check-circle':
      return CheckCircle;
    default:
      return null;
  }
});
</script>
