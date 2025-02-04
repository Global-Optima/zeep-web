<template>
	<div>
		<ul>
			<li
				v-for="(shift, index) in shifts"
				:key="index"
				class="flex justify-between py-2 border-b last:border-b-0"
			>
				<span class="text-gray-900">{{ shift.date }}</span>
				<span class="text-gray-600">{{ shift.shift }}</span>
			</li>
		</ul>
	</div>
</template>

<script setup lang="ts">
interface Shift {
  date: string;
  shift: string;
}

defineProps<{
  shifts: Shift[];
}>();
</script>
