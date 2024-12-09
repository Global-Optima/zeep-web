<template>
	<div class="mx-auto w-full max-w-2xl">
		<!-- Header -->
		<div class="w-full flex items-center justify-between gap-4">
			<div class="flex items-center gap-4">
				<Button
					variant="outline"
					size="icon"
					@click="onBackClick"
				>
					<ChevronLeft class="w-5 h-5" />
					<span class="sr-only">Назад</span>
				</Button>
				<h1
					class="flex-1 sm:grow-0 font-semibold text-xl tracking-tight whitespace-nowrap shrink-0"
				>
					{{ 'Детали варианта' }}
				</h1>
			</div>

			<div class="flex items-center gap-2">
				<Button variant="outline"> Отменить </Button>
				<Button> Сохранить </Button>
			</div>
		</div>

		<!-- Main Content -->
		<div class="items-start gap-4 grid lg:col-span-2 auto-rows-max mt-6">
			<!-- Variant Details -->
			<Card>
				<CardHeader>
					<CardTitle>Детали варианта</CardTitle>
					<CardDescription>Укажите название, размер и цену варианта.</CardDescription>
				</CardHeader>
				<CardContent>
					<div class="gap-6 grid">
						<!-- Name Selector -->
						<div class="gap-3 grid">
							<Label for="variant-name">Название</Label>
							<Select>
								<SelectTrigger id="variant-name">
									<SelectValue placeholder="Выберите название" />
								</SelectTrigger>
								<SelectContent>
									<SelectItem value="S">S</SelectItem>
									<SelectItem value="M">M</SelectItem>
									<SelectItem value="L">L</SelectItem>
									<SelectItem value="XL">XL</SelectItem>
								</SelectContent>
							</Select>
						</div>

						<!-- Measure Input -->
						<div class="flex gap-4 items-center">
							<div class="gap-3 grid flex-1">
								<Label for="measure">Размер (число)</Label>
								<Input
									id="measure"
									type="number"
									class="w-full"
									placeholder="Например, 250"
								/>
							</div>
							<div class="gap-3 grid flex-1">
								<Label for="unit">Единица измерения</Label>
								<Select>
									<SelectTrigger id="unit">
										<SelectValue placeholder="Выберите единицу" />
									</SelectTrigger>
									<SelectContent>
										<SelectItem value="грамм">Грамм</SelectItem>
										<SelectItem value="мл">Миллилитры</SelectItem>
									</SelectContent>
								</Select>
							</div>
						</div>

						<!-- Price Input -->
						<div class="gap-3 grid">
							<Label for="price">Начальная цена</Label>
							<Input
								id="price"
								type="number"
								class="w-full"
								placeholder="Введите начальную цену"
							/>
						</div>
					</div>
				</CardContent>
			</Card>

			<!-- Default Toppings Block -->
			<Card>
				<CardHeader>
					<div class="flex justify-between items-start gap-4">
						<div>
							<CardTitle>Топпинги по умолчанию</CardTitle>
							<CardDescription class="mt-2">
								Добавьте топпинги, которые идут по умолчанию с этим вариантом.
							</CardDescription>
						</div>
						<div>
							<Button variant="outline">Добавить</Button>
						</div>
					</div>
				</CardHeader>
				<CardContent>
					<Table>
						<TableHeader>
							<TableRow>
								<TableHead class="hidden w-[100px] sm:table-cell"> </TableHead>
								<TableHead>Название</TableHead>
								<TableHead>Размер</TableHead>
							</TableRow>
						</TableHeader>
						<TableBody>
							<TableRow
								v-for="additive in defaultAdditives"
								:key="additive.id"
							>
								<TableCell class="hidden sm:table-cell">
									<img
										:src="additive.imageUrl"
										alt="Изображение"
										class="bg-gray-100 rounded-md aspect-square object-contain"
										height="64"
										width="64"
									/>
								</TableCell>
								<TableCell class="font-medium">{{ additive.name }}</TableCell>
								<TableCell>{{ additive.size }}</TableCell>
							</TableRow>
						</TableBody>
					</Table>
				</CardContent>
			</Card>

			<!-- Technical Map Block -->
			<Card>
				<CardHeader>
					<div class="flex justify-between items-start gap-4">
						<div>
							<CardTitle>Техническая карта</CardTitle>
							<CardDescription class="mt-2">
								Выберите ингредиенты и укажите их вес для данного варианта.
							</CardDescription>
						</div>
						<div>
							<Button variant="outline">Добавить</Button>
						</div>
					</div>
				</CardHeader>
				<CardContent>
					<Table>
						<TableHeader>
							<TableRow>
								<TableHead class="hidden w-[100px] sm:table-cell"></TableHead>
								<TableHead>Название</TableHead>
								<TableHead>Единица измерения</TableHead>
								<TableHead>Вес</TableHead>
							</TableRow>
						</TableHeader>
						<TableBody>
							<TableRow
								v-for="ingredient in technicalCardIngredients"
								:key="ingredient.id"
							>
								<TableCell class="hidden sm:table-cell">
									<img
										:src="ingredient.imageUrl"
										alt="Изображение ингредиента"
										class="bg-gray-100 rounded-md aspect-square object-contain"
										height="64"
										width="64"
									/>
								</TableCell>
								<TableCell class="font-medium">{{ ingredient.name }}</TableCell>
								<TableCell>{{ ingredient.unit }}</TableCell>
								<TableCell>
									<Input
										type="number"
										v-model="ingredient.weight"
										placeholder="Введите вес"
										class="w-full"
									/>
								</TableCell>
							</TableRow>
						</TableBody>
					</Table>
				</CardContent>
			</Card>
		</div>
	</div>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import { Input } from '@/core/components/ui/input'
import { Label } from '@/core/components/ui/label'
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/core/components/ui/select'
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from '@/core/components/ui/table'
import { ChevronLeft } from 'lucide-vue-next'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const defaultAdditives = ref([
	{ id: 1, name: 'Карамельный сироп', size: '500 мл', imageUrl: 'caramel.jpg' },
	{ id: 2, name: 'Ванильный сироп', size: '500 мл', imageUrl: 'vanilla.jpg' },
])

const technicalCardIngredients = ref([
	{ id: 1, imageUrl: 'ingredient1.jpg', name: 'Кофе', unit: 'грамм', weight: 30 },
	{ id: 2, imageUrl: 'ingredient2.jpg', name: 'Молоко', unit: 'мл', weight: 100 },
])

const onBackClick = () => {
	router.back()
}
</script>
