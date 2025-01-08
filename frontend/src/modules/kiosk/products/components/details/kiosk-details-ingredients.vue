<template>
	<Accordion
		type="single"
		collapsible
	>
		<AccordionItem
			value="item-1"
			class="bg-white border-none rounded-3xl"
		>
			<AccordionTrigger class="px-7 py-6 font-normal text-2xl hover:no-underline">
				Ингредиенты
			</AccordionTrigger>
			<AccordionContent class="px-2 pb-4">
				<section>
					<p class="px-5 font-normal text-base sm:text-2xl">
						{{ingredientNames}}
					</p>
				</section>
				<section class="justify-items-center gap-2 sm:gap-4 grid grid-cols-4 mt-6">
					<div class="flex flex-col items-start">
						<p class="text-gray-600 text-xs sm:text-xl">Энергия</p>
						<p
							class="font-medium text-base sm:text-2xl"
							data-testid="energy-ccal"
						>
							{{ energy.ccal }} ккал
						</p>
					</div>

					<div class="flex flex-col items-start">
						<p class="text-gray-600 text-xs sm:text-lg">Белки</p>
						<p
							class="font-medium text-base sm:text-2xl"
							data-testid="energy-protein"
						>
							{{ energy.proteins }} г
						</p>
					</div>

					<div class="flex flex-col items-start">
						<p class="text-gray-600 text-xs sm:text-lg">Жиры</p>
						<p
							class="font-medium text-base sm:text-2xl"
							data-testid="energy-fats"
						>
							{{ energy.fats }} г
						</p>
					</div>

					<div class="flex flex-col items-start">
						<p class="text-gray-600 text-xs sm:text-lg">Углеводы</p>
						<p
							class="font-medium text-base sm:text-2xl"
							data-testid="energy-carbs"
						>
							{{ energy.carbs }} г
						</p>
					</div>
				</section>
			</AccordionContent>
		</AccordionItem>
	</Accordion>
</template>

<script setup lang="ts">
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/core/components/ui/accordion'
import type { IngredientsDTO } from '@/modules/admin/ingredients/models/ingredients.model'
import type { KioskDetailsEnergyDTO } from '@/modules/kiosk/products/components/details/kiosk-details-product-info.vue'
import { computed } from 'vue'

const {ingredients} = defineProps<{
  energy: KioskDetailsEnergyDTO
  ingredients: IngredientsDTO[]
}>()

const ingredientNames = computed(() => ingredients.map(i => i.name).join(", "))
</script>

<style scoped></style>
