<script setup lang="ts">
import { cn } from '@/core/utils/tailwind.utils'
import { ChevronDownIcon } from '@radix-icons/vue'
import {
  AccordionHeader,
  AccordionTrigger,
  type AccordionTriggerProps,
} from 'radix-vue'
import { computed, type HTMLAttributes } from 'vue'

const props = defineProps<AccordionTriggerProps & { class?: HTMLAttributes['class'] }>()

const delegatedProps = computed(() => {
  const { class: _, ...delegated } = props

  return delegated
})
</script>

<template>
	<AccordionHeader class="flex">
		<AccordionTrigger
			v-bind="delegatedProps"
			:class="
        cn(
          'flex flex-1 items-center justify-between py-4 text-sm font-medium transition-all hover:underline [&[data-state=open]>svg]:rotate-180',
          props.class,
        )
      "
		>
			<slot />
			<slot name="icon">
				<ChevronDownIcon
					class="w-8 h-8 text-muted-foreground transition-transform duration-200 shrink-0"
				/>
			</slot>
		</AccordionTrigger>
	</AccordionHeader>
</template>
