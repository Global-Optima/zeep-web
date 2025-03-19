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

const {includeCloseButton=true, ...props} = defineProps<DialogContentProps & { class?: HTMLAttributes['class'], includeCloseButton?: boolean }>()
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
			v-bind="forwarded"
			:class="
        cn(
          'fixed left-1/2 top-1/2 z-50 grid w-full max-w-lg -translate-x-1/2 -translate-y-1/2 gap-4 border bg-background p-6 shadow-lg duration-200 data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[state=closed]:slide-out-to-left-1/2 data-[state=closed]:slide-out-to-top-[48%] data-[state=open]:slide-in-from-left-1/2 data-[state=open]:slide-in-from-top-[48%] sm:rounded-lg',
          props.class,
        )"
		>
			<slot />

			<DialogClose
				v-if="includeCloseButton"
				class="top-4 right-4 absolute bg-slate-800/70 data-[state=open]:bg-secondary hover:opacity-100 backdrop-blur-md p-2 rounded-full focus:ring-2 focus:ring-ring ring-offset-background focus:ring-offset-2 text-white transition-opacity focus:outline-none disabled:pointer-events-none"
			>
				<Cross2Icon class="w-5 sm:w-8 h-5 sm:h-8" />
			</DialogClose>
		</DialogContent>
	</DialogPortal>
</template>
