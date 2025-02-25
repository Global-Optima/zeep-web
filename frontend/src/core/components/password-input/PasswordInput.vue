<script setup lang="ts">
import { cn } from '@/core/utils/tailwind.utils'
import { useVModel } from '@vueuse/core'
import { Eye, EyeOff } from 'lucide-vue-next'
import type { HTMLAttributes } from 'vue'
import { computed, ref } from 'vue'
import Input from '../ui/input/Input.vue'

const props = defineProps<{
  defaultValue?: string
  placeholder?: string
  modelValue?: string
  class?: HTMLAttributes['class']
}>()

const emits = defineEmits<{
  (e: 'update:modelValue', payload: string): void
}>()

const modelValue = useVModel(props, 'modelValue', emits, {
  passive: true,
  defaultValue: props.defaultValue ?? '',
})

const showPassword = ref(false)
const inputType = computed(() => (showPassword.value ? 'text' : 'password'))

const togglePassword = () => {
  showPassword.value = !showPassword.value
}
</script>

<template>
	<div class="relative w-full">
		<Input
			v-model="modelValue"
			:type="inputType"
			:placeholder="placeholder"
			:class="cn(
        props.class
      )"
		/>
		<button
			type="button"
			@click="togglePassword"
			class="right-2 absolute inset-y-0 flex items-center p-2 text-gray-500 hover:text-gray-700"
		>
			<component
				:is="showPassword ? EyeOff : Eye"
				class="w-5 h-5"
			/>
		</button>
	</div>
</template>
