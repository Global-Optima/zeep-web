<template>
	<section class="bg-white rounded-3xl grid grid-cols-4 p-6 gap-2 sm:gap-4 justify-items-center">
		<div class="flex flex-col items-start">
			<p class="text-xs sm:text-lg text-gray-600">Энергия</p>
			<p
				class="text-base sm:text-2xl font-medium"
				data-testid="energy-ccal"
			>
				{{ energy.ccal }} ккал
			</p>
		</div>

		<div class="flex flex-col items-start">
			<p class="text-xs sm:text-lg text-gray-600">Белки</p>
			<p
				class="text-base sm:text-2xl font-medium"
				data-testid="energy-protein"
			>
				{{ energy.proteins }} г
			</p>
		</div>

		<div class="flex flex-col items-start">
			<p class="text-xs sm:text-lg text-gray-600">Жиры</p>
			<p
				class="text-base sm:text-2xl font-medium"
				data-testid="energy-fats"
			>
				{{ energy.fats }} г
			</p>
		</div>

		<div class="flex flex-col items-start">
			<p class="text-xs sm:text-lg text-gray-600">Углеводы</p>
			<p
				class="text-base sm:text-2xl font-medium"
				data-testid="energy-carbs"
			>
				{{ energy.carbs }} г
			</p>
		</div>
	</section>
</template>

<script setup lang="ts">
const {energy} = defineProps<{energy: {ccal: number, proteins: number, fats: number, carbs: number}}>()
</script>

<style scoped></style>
