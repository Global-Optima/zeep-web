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
				{{ 'Название добавки' }}
			</h1>

			<div class="md:flex items-center gap-2 hidden md:ml-auto">
				<Button variant="outline"> Отменить </Button>
				<Button> Сохранить </Button>
			</div>
		</div>

		<!-- Main Content -->
		<div class="gap-4 grid md:grid-cols-[1fr_250px] lg:grid-cols-3">
			<div class="items-start gap-4 grid lg:col-span-2 auto-rows-max">
				<!-- Additive Details -->
				<Card>
					<CardHeader>
						<CardTitle>Детали добавки</CardTitle>
						<CardDescription>Введите название, описание и цену добавки.</CardDescription>
					</CardHeader>
					<CardContent>
						<div class="gap-6 grid">
							<div class="gap-3 grid">
								<Label for="name">Название</Label>
								<Input
									id="name"
									type="text"
									class="w-full"
									v-model="additiveName"
									placeholder="Введите название добавки"
								/>
							</div>
							<div class="gap-3 grid">
								<Label for="description">Описание</Label>
								<Textarea
									id="description"
									v-model="additiveDescription"
									placeholder="Краткое описание добавки"
									class="min-h-32"
								/>
							</div>
							<div class="gap-3 grid">
								<Label for="price">Цена</Label>
								<Input
									id="price"
									type="number"
									class="w-full"
									v-model="additivePrice"
									placeholder="Введите цену добавки"
								/>
							</div>
							<div class="gap-3 grid">
								<Label for="size">Размер</Label>
								<Input
									id="size"
									type="text"
									class="w-full"
									v-model="additiveSize"
									placeholder="Например: 500 мл"
								/>
							</div>
						</div>
					</CardContent>
				</Card>
			</div>

			<!-- Media and Category Blocks -->
			<div class="items-start gap-4 grid auto-rows-max">
				<!-- Media Block -->
				<Card>
					<CardHeader>
						<CardTitle>Медиа</CardTitle>
						<CardDescription>Загрузите изображение добавки.</CardDescription>
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

				<!-- Category Block -->
				<Card>
					<CardHeader>
						<CardTitle>Категория</CardTitle>
					</CardHeader>
					<CardContent>
						<div class="gap-6 grid sm:grid-cols-1">
							<div class="gap-3 grid">
								<Label for="category">Категория</Label>
								<Select>
									<SelectTrigger
										id="category"
										aria-label="Выберите категорию"
									>
										<SelectValue placeholder="Выберите категорию" />
									</SelectTrigger>
									<SelectContent>
										<SelectItem value="соусы">Соусы</SelectItem>
										<SelectItem value="сиропы">Сиропы</SelectItem>
										<SelectItem value="топпинги">Топпинги</SelectItem>
									</SelectContent>
								</Select>
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
import { Textarea } from '@/core/components/ui/textarea'
import { ChevronLeft } from 'lucide-vue-next'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

// Reactive properties for additive details
const additiveName = ref('')
const additiveDescription = ref('')
const additivePrice = ref(0)
const additiveSize = ref('')

// Router for navigation
const router = useRouter()

const onBackClick = () => {
  router.push({ name: 'ADMIN_ADDITIVES' })
}
</script>

<style scoped>
/* Add any custom styles here */
</style>
