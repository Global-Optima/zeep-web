<script setup lang="ts">
import { cn } from '@/core/utils/tailwind.utils'
import type { HTMLAttributes } from 'vue'

const props = defineProps<{
  class?: HTMLAttributes['class']
}>()
</script>

<template>
	<div
		:class="
      cn(
        'rounded-xl border bg-card text-card-foreground',
        props.class,
      )
    "
	>
		<slot />
	</div>
</template>
