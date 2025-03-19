<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger
} from '@/core/components/ui/dialog'
import type { ProductTotalNutrition } from '@/modules/kiosk/products/models/product.model'
import { NotebookText } from 'lucide-vue-next'
import { computed } from 'vue'

const {nutrition} = defineProps<{nutrition: ProductTotalNutrition}>()

const ingredientsRecipe = computed(() => nutrition.ingredients.join(", "))
const allergensList = computed(() => nutrition.allergenIngredients.join(", "))
</script>

<template>
	<Dialog>
		<DialogTrigger as-child>
			<Button
				size="icon"
				variant="ghost"
				class="size-14"
			>
				<NotebookText
					class="!size-14 text-gray-400"
					stroke-width="1.6"
				/>
			</Button>
		</DialogTrigger>
		<DialogContent
			:include-close-button="false"
			class="!rounded-3xl"
		>
			<DialogHeader>
				<DialogTitle class="text-2xl">Ингредиенты</DialogTitle>
			</DialogHeader>
			<div class="flex flex-col gap-8">
				<div class="gap-8 grid grid-cols-4 shadow-md p-4 rounded-xl text-lg">
					<div>
						<p class="text-gray-500">Энергии</p>
						<p class="mt-1 text-xl">{{nutrition.calories}} ккал</p>
					</div>

					<div>
						<p class="text-gray-500">Белки</p>
						<p class="text-xl">{{nutrition.proteins}} гр</p>
					</div>

					<div>
						<p class="text-gray-500">Жиры</p>
						<p class="text-xl">{{nutrition.fats}} гр</p>
					</div>

					<div>
						<p class="text-gray-500">Углеводы</p>
						<p class="text-xl">{{nutrition.carbs}} гр</p>
					</div>
				</div>

				<div>
					<p class="font-semibold text-xl">Состав</p>

					<p class="mt-1 text-gray-500 text-lg">
						{{ingredientsRecipe}}
					</p>
				</div>

				<div>
					<p class="font-semibold text-xl">Аллергены</p>

					<p class="mt-1 text-gray-500 text-lg">
						{{allergensList}}
					</p>
				</div>
			</div>
		</DialogContent>
	</Dialog>
</template>
