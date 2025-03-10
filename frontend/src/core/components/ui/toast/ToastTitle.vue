<script setup lang="ts">
import { cn } from '@/core/utils/tailwind.utils'
import { ToastTitle, type ToastTitleProps } from 'radix-vue'
import { computed, type HTMLAttributes } from 'vue'

const props = defineProps<ToastTitleProps & { class?: HTMLAttributes['class'] }>()

const delegatedProps = computed(() => {
  const { class: _, ...delegated } = props

  return delegated
})
</script>

<template>
	<ToastTitle
		v-bind="delegatedProps"
		:class="cn('lg:text-base text-xl font-semibold', props.class)"
	>
		<slot />
	</ToastTitle>
</template>
