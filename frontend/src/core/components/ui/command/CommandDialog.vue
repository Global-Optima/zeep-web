<script setup lang="ts">
import { Dialog, DialogContent } from '@/core/components/ui/dialog'
import type { DialogRootEmits, DialogRootProps } from 'radix-vue'
import { useForwardPropsEmits } from 'radix-vue'
import Command from './Command.vue'

const props = defineProps<DialogRootProps>()
const emits = defineEmits<DialogRootEmits>()

const forwarded = useForwardPropsEmits(props, emits)
</script>

<template>
	<Dialog v-bind="forwarded">
		<DialogContent
			:include-close-button="false"
			class="shadow-lg p-0 overflow-hidden"
		>
			<Command
				class="[&_[cmdk-group-heading]]:px-2 [&_[cmdk-group-heading]]:font-medium [&_[cmdk-group-heading]]:text-muted-foreground [&_[cmdk-group]]:px-2 [&_[cmdk-group]:not([hidden])_~[cmdk-group]]:pt-0 [&_[cmdk-item]]:px-2 [&_[cmdk-item]]:py-3 [&_[cmdk-input-wrapper]_svg]:w-5 [&_[cmdk-item]_svg]:w-5 [&_[cmdk-input-wrapper]_svg]:h-5 [&_[cmdk-input]]:h-12 [&_[cmdk-item]_svg]:h-5"
			>
				<slot />
			</Command>
		</DialogContent>
	</Dialog>
</template>
