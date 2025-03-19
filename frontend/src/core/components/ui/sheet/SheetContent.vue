<script setup lang="ts">
import { cn } from '@/core/utils/tailwind.utils'
import { Cross2Icon } from '@radix-icons/vue'
import {
  DialogClose,
  DialogContent,
  type DialogContentEmits,
  type DialogContentProps,
  DialogOverlay,
  DialogPortal,
  useForwardPropsEmits,
} from 'radix-vue'
import { computed, type HTMLAttributes } from 'vue'
import { type SheetVariants, sheetVariants } from '.'

interface SheetContentProps extends DialogContentProps {
  class?: HTMLAttributes['class']
  side?: SheetVariants['side']
  includeCloseButton?: boolean
}

defineOptions({
  inheritAttrs: false,
})

const {includeCloseButton = true, ...props} = defineProps<SheetContentProps>()

const emits = defineEmits<DialogContentEmits>()

const delegatedProps = computed(() => {
  const { class: _, ...delegated } = props

  return delegated
})

const forwarded = useForwardPropsEmits(delegatedProps, emits)
</script>

<template>
	<DialogPortal>
		<DialogOverlay
			class="z-50 fixed inset-0 bg-black/80 data-[state=closed]:animate-out data-[state=open]:animate-in data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0"
		/>
		<DialogContent
			:class="cn(sheetVariants({ side }), props.class)"
			v-bind="{ ...forwarded, ...$attrs }"
		>
			<slot />

			<DialogClose
				v-if="includeCloseButton"
				class="top-6 right-6 absolute bg-slate-800/70 data-[state=open]:bg-secondary hover:opacity-100 backdrop-blur-md p-2 rounded-full focus:ring-2 focus:ring-ring ring-offset-background focus:ring-offset-2 text-white transition-opacity focus:outline-none disabled:pointer-events-none"
			>
				<Cross2Icon class="w-6 sm:w-8 h-6 sm:h-8" />
			</DialogClose>
		</DialogContent>
	</DialogPortal>
</template>
