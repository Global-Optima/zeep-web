<template>
	<Accordion
		type="single"
		collapsible
	>
		<AccordionItem
			value="item-1"
			class="bg-white border-none rounded-3xl"
		>
			<AccordionTrigger class="px-7 py-6 text-2xl hover:no-underline font-normal">
				Ингредиенты
			</AccordionTrigger>
			<AccordionContent class="px-2 pb-4">
				<section>
					<p class="text-base font-normal sm:text-2xl px-5">
						Кофе, молоко, сахар, лед, сливки, ванильный сироп
					</p>
				</section>
				<section class="justify-items-center gap-2 sm:gap-4 grid grid-cols-4 mt-6">
					<div class="flex flex-col items-start">
						<p class="text-gray-600 text-xs sm:text-xl">Энергия</p>
						<p
							class="font-medium text-base sm:text-2xl"
							data-testid="energy-ccal"
						>
							340 ккал
						</p>
					</div>

					<div class="flex flex-col items-start">
						<p class="text-gray-600 text-xs sm:text-lg">Белки</p>
						<p
							class="font-medium text-base sm:text-2xl"
							data-testid="energy-protein"
						>
							22.3 г
						</p>
					</div>

					<div class="flex flex-col items-start">
						<p class="text-gray-600 text-xs sm:text-lg">Жиры</p>
						<p
							class="font-medium text-base sm:text-2xl"
							data-testid="energy-fats"
						>
							12.3 г
						</p>
					</div>

					<div class="flex flex-col items-start">
						<p class="text-gray-600 text-xs sm:text-lg">Углеводы</p>
						<p
							class="font-medium text-base sm:text-2xl"
							data-testid="energy-carbs"
						>
							32.4 г
						</p>
					</div>
				</section>
			</AccordionContent>
		</AccordionItem>
	</Accordion>
</template>

<script setup lang="ts">
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/core/components/ui/accordion'
</script>

<style scoped></style>
