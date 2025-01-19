<template>
	<div class="gap-4 grid">
		<div
			v-for="(day, index) in weekdays"
			:key="index"
			class="flex items-center gap-4"
		>
			<Label class="w-24">{{ day }}</Label>
			<div class="flex items-center gap-2">
				<Input
					v-model="workingHours[day].start"
					type="time"
					:disabled="workingHours[day].isDayOff"
				/>
				<span>-</span>
				<Input
					v-model="workingHours[day].end"
					type="time"
					:disabled="workingHours[day].isDayOff"
				/>
			</div>
			<div class="flex items-center">
				<Checkbox v-model="workingHours[day].isDayOff" />
				<span class="ml-2">Выходной</span>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { Checkbox } from '@/core/components/ui/checkbox'
import { Input } from '@/core/components/ui/input'
import { Label } from '@/core/components/ui/label'
import { reactive, watch } from 'vue'

const props = defineProps<{
  modelValue: Record<string, { start: string; end: string; isDayOff: boolean }>;
}>();

const emits = defineEmits(['update:modelValue']);

const weekdays = [
  'Понедельник',
  'Вторник',
  'Среда',
  'Четверг',
  'Пятница',
  'Суббота',
  'Воскресенье',
];

const workingHours = reactive(props.modelValue || {});

// Initialize working hours if empty
weekdays.forEach((day) => {
  if (!workingHours[day]) {
    workingHours[day] = { start: '', end: '', isDayOff: false };
  }
});

// Watch for changes and emit update
watch(
  () => workingHours,
  () => {
    emits('update:modelValue', { ...workingHours });
  },
  { deep: true }
);
</script>
