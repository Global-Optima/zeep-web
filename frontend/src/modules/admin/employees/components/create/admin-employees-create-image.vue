<template>
	<div class="flex flex-col items-start">
		<div class="mb-4 w-32 h-32">
			<img
				v-if="previewUrl"
				:src="previewUrl"
				alt="Изображение профиля"
				class="rounded-full w-full h-full object-cover"
			/>
			<div
				v-else
				class="flex justify-center items-center bg-gray-200 rounded-full w-full h-full"
			>
				<UserIcon class="w-16 h-16 text-gray-500" />
			</div>
		</div>
		<Button
			variant="outline"
			@click="triggerFileInput"
		>
			Загрузить изображение
		</Button>
		<input
			ref="fileInput"
			type="file"
			accept="image/*"
			class="hidden"
			@change="handleFileChange"
		/>
	</div>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { User as UserIcon } from 'lucide-vue-next'
import { ref, watch } from 'vue'

const props = defineProps<{
  modelValue: File | null;
}>();

const emits = defineEmits(['update:modelValue']);

const fileInput = ref<HTMLInputElement | null>(null);
const previewUrl = ref<string | null>(null);

// Trigger file input click
const triggerFileInput = () => {
  fileInput.value?.click();
};

// Handle file selection
const handleFileChange = (event: Event) => {
  const files = (event.target as HTMLInputElement).files;
  if (files && files[0]) {
    emits('update:modelValue', files[0]);
  }
};

// Watch for changes to modelValue to update preview
watch(
  () => props.modelValue,
  (newFile) => {
    if (newFile) {
      previewUrl.value = URL.createObjectURL(newFile);
    } else {
      previewUrl.value = null;
    }
  }
);
</script>
