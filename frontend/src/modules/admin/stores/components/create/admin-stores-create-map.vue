<template>
	<div class="space-y-4">
		<!-- Address Input -->
		<div>
			<Label for="address">Адрес</Label>
			<Input
				id="address"
				v-model="address"
				placeholder="Введите адрес"
			/>
		</div>
		<!-- Latitude and Longitude Inputs -->
		<div class="gap-4 grid grid-cols-1 sm:grid-cols-2">
			<div>
				<Label for="latitude">Широта</Label>
				<Input
					id="latitude"
					type="number"
					v-model="latitude"
					placeholder="55.7558"
					step="any"
				/>
			</div>
			<div>
				<Label for="longitude">Долгота</Label>
				<Input
					id="longitude"
					type="number"
					v-model="longitude"
					placeholder="37.6173"
					step="any"
				/>
			</div>
		</div>
		<!-- Placeholder for Map -->
		<div class="flex justify-center items-center bg-gray-100 mt-4 w-full h-64">
			<p class="text-gray-500">Здесь будет отображаться карта</p>
		</div>
	</div>
</template>

<script setup lang="ts">
import { Input } from '@/core/components/ui/input'
import { Label } from '@/core/components/ui/label'
import { ref, watch } from 'vue'

const props = defineProps<{
  modelValue: { address: string; lat: number; lng: number } | null;
}>();

const emits = defineEmits(['update:modelValue']);

const address = ref('');
const latitude = ref<number>(0);
const longitude = ref<number>(0);

// Initialize values if modelValue is provided
if (props.modelValue) {
  address.value = props.modelValue.address;
  latitude.value = props.modelValue.lat;
  longitude.value = props.modelValue.lng;
}

// Watch for changes in inputs and emit updated value
watch(
  () => ({ address: address.value, lat: latitude.value, lng: longitude.value }),
  (newValue) => {
    emits('update:modelValue', newValue);
  },
  { deep: true }
);
</script>
