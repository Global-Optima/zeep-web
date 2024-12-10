<template>
	<div class="flex-1 gap-4 grid auto-rows-max mx-auto max-w-6xl">
		<!-- Header -->
		<div class="flex items-center gap-4">
			<Button
				variant="outline"
				size="icon"
				@click="onBackClick"
			>
				<ChevronLeft class="w-5 h-5" />
				<span class="sr-only">Назад</span>
			</Button>
			<h1 class="flex-1 sm:grow-0 font-semibold text-xl tracking-tight whitespace-nowrap shrink-0">
				{{ 'Детали ингредиента' }}
			</h1>

			<div class="md:flex items-center gap-2 hidden md:ml-auto">
				<Button variant="outline"> Отменить </Button>
				<Button> Сохранить </Button>
			</div>
		</div>

		<!-- Main Content -->
		<div class="gap-4 grid md:grid-cols-[1fr_250px] lg:grid-cols-3">
			<div class="items-start gap-4 grid lg:col-span-2 auto-rows-max">
				<!-- Ingredient Details -->
				<Card>
					<CardHeader>
						<CardTitle>Основная информация</CardTitle>
						<CardDescription>Введите данные ингредиента.</CardDescription>
					</CardHeader>
					<CardContent>
						<div class="gap-6 grid">
							<!-- Name -->
							<div class="gap-3 grid">
								<Label for="name">Название</Label>
								<Input
									id="name"
									type="text"
									class="w-full"
									v-model="ingredientName"
									placeholder="Введите название ингредиента"
								/>
							</div>

							<!-- Nutrition Facts -->
							<div class="gap-4 grid grid-cols-4">
								<!-- Calories -->
								<div class="w-full">
									<Label for="calories">Калории</Label>
									<Input
										id="calories"
										type="number"
										class="w-full"
										v-model="ingredientCalories"
										placeholder="Введите калории (ккал)"
									/>
									<p class="mt-1 text-gray-500 text-sm">Укажите в ккал</p>
								</div>
								<!-- Proteins -->
								<div class="w-full">
									<Label for="proteins">Белки</Label>
									<Input
										id="proteins"
										type="number"
										class="w-full"
										v-model="ingredientProteins"
										placeholder="Введите белки (г)"
									/>
									<p class="mt-1 text-gray-500 text-sm">Укажите в граммах</p>
								</div>
								<!-- Carbs -->
								<div class="w-full">
									<Label for="carbs">Углеводы</Label>
									<Input
										id="carbs"
										type="number"
										class="w-full"
										v-model="ingredientCarbs"
										placeholder="Введите углеводы (г)"
									/>
									<p class="mt-1 text-gray-500 text-sm">Укажите в граммах</p>
								</div>
								<!-- Fats -->
								<div class="w-full">
									<Label for="fats">Жиры</Label>
									<Input
										id="fats"
										type="number"
										class="w-full"
										v-model="ingredientFats"
										placeholder="Введите жиры (г)"
									/>
									<p class="mt-1 text-gray-500 text-sm">Укажите в граммах</p>
								</div>
							</div>

							<!-- Expiration Date -->
							<div class="gap-3 grid">
								<Label for="expiration">Срок годности (дней)</Label>
								<Input
									id="expiration"
									type="number"
									class="w-full"
									v-model="ingredientExpiration"
									placeholder="Введите срок годности"
								/>
							</div>

							<!-- Measurement Unit -->
							<div class="gap-3 grid">
								<Label for="unit">Единица измерения</Label>
								<Select>
									<SelectTrigger
										id="unit"
										aria-label="Выберите единицу измерения"
									>
										<SelectValue placeholder="Выберите единицу измерения" />
									</SelectTrigger>
									<SelectContent>
										<SelectItem value="грамм">Грамм</SelectItem>
										<SelectItem value="литр">Литр</SelectItem>
										<SelectItem value="миллилитр">Миллилитр</SelectItem>
										<SelectItem value="килограмм">Килограмм</SelectItem>
									</SelectContent>
								</Select>
							</div>
						</div>
					</CardContent>
				</Card>
			</div>

			<!-- Media Block -->
			<div class="items-start gap-4 grid auto-rows-max">
				<Card>
					<CardHeader>
						<CardTitle>Медиа</CardTitle>
						<CardDescription>Загрузите изображение ингредиента.</CardDescription>
					</CardHeader>
					<CardContent>
						<div class="gap-4 grid">
							<div class="flex flex-col gap-3">
								<Label>Изображение</Label>
								<Input
									type="file"
									accept="image/*"
									class="file-input"
								/>
							</div>
						</div>
					</CardContent>
				</Card>
			</div>
		</div>

		<!-- Footer for Mobile -->
		<div class="flex justify-center items-center gap-2 md:hidden">
			<Button variant="outline"> Отменить </Button>
			<Button> Сохранить </Button>
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
import { ChevronLeft } from 'lucide-vue-next'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

// Reactive properties for ingredient details
const ingredientName = ref('')
const ingredientCalories = ref(0)
const ingredientProteins = ref(0)
const ingredientCarbs = ref(0)
const ingredientFats = ref(0)
const ingredientExpiration = ref(0)

// Router for navigation
const router = useRouter()

const onBackClick = () => {
  router.push({ name: 'ADMIN_INGREDIENTS' })
}
</script>

<style scoped>
/* Add any custom styles here */
</style>
