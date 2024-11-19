<template>
	<div class="w-full h-screen overflow-hidden">
		<header class="flex justify-center items-center gap-1 bg-white py-4 border-b">
			<button
				v-for="status in statuses"
				:key="status.label"
				:class="cn('flex items-center gap-2 px-5 py-2 rounded-xl text-lg', status.label === selectedStatus?.label ? 'bg-primary text-primary-foreground' :'')"
				@click="selectedStatus = status"
			>
				<p>{{ status.label }}</p>
				<p
					:class="cn('bg-gray-100 px-2 py-1 rounded-sm text-black text-xs', status.label === selectedStatus?.label ? 'bg-green-700 text-primary-foreground' :'')"
				>
					{{ status.count }}
				</p>
			</button>
		</header>

		<main class="grid grid-cols-4 bg-gray-100 w-full h-full">
			<section class="col-span-1 p-2 border-r h-full">
				<p class="font-medium text-center">Заказы</p>
				<div class="flex flex-col gap-2 mt-2 overflow-y-auto no-scrollbar">
					<div class="flex justify-between items-start gap-2 bg-white p-4 rounded-xl">
						<div>
							<div class="flex items-center gap-2 text-lg">
								<p class="font-medium">Mooray</p>
								<p class="text-gray-500">#3302</p>
							</div>

							<p class="mt-1 text-gray-700 text-sm">Доставка, 3 шт</p>
						</div>

						<div>~5 мин</div>
					</div>
				</div>
			</section>
			<section class="col-span-1 p-2 border-r h-full">
				<p class="font-medium text-center">Подзаказы</p>
			</section>
			<section class="col-span-2 p-2 border-r h-full">3</section>
		</main>
	</div>
</template>

<script setup lang="ts">
import { cn } from '@/core/utils/tailwind.utils'
import { ref } from 'vue'

interface Status {
  label: string
  count: number
}

const statuses = ref<Status[]>([
  {label: "Все", count: 12},
  {label: "Активные", count: 4},
  {label: "Завершенные", count: 6},
  {label: "В доставке", count: 2},
])

const selectedStatus = ref<Status | null>(null)
</script>

<style scoped></style>
