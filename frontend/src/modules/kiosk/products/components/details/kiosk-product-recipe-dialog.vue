<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger
} from '@/core/components/ui/dialog'
import { NotebookText } from 'lucide-vue-next'
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
						<p class="mt-1 text-xl">502 ккал</p>
					</div>

					<div>
						<p class="text-gray-500">Белки</p>
						<p class="text-xl">2.4 гр</p>
					</div>

					<div>
						<p class="text-gray-500">Жиры</p>
						<p class="text-xl">5.2 гр</p>
					</div>

					<div>
						<p class="text-gray-500">Углеводы</p>
						<p class="text-xl">10.6 гр</p>
					</div>
				</div>

				<div>
					<p class="font-semibold text-xl">Состав</p>

					<p class="mt-1 text-gray-500 text-lg">
						хлеб на ржаной закваске, специальный соус (майонез, соус шрирача, зерновая горчица,
						трюфельное масло, халапеньо), сыр чеддер, сыр моцарелла, говяжий пастрами, квашеная
						капуста, маринованный огурец
					</p>
				</div>

				<div>
					<p class="font-semibold text-xl">Аллергены</p>

					<p class="mt-1 text-gray-500 text-lg">
						молоко (лактоза), яйца, глютен, горчицакапуста, маринованный огурец
					</p>
				</div>
			</div>
		</DialogContent>
	</Dialog>
</template>
