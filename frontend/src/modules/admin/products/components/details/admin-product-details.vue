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
				{{ 'Название продукта' }}
			</h1>

			<div class="md:flex items-center gap-2 hidden md:ml-auto">
				<Button variant="outline"> Отменить </Button>
				<Button> Сохранить </Button>
			</div>
		</div>

		<!-- Main Content -->
		<div class="gap-4 grid md:grid-cols-[1fr_250px] lg:grid-cols-3">
			<div class="items-start gap-4 grid lg:col-span-2 auto-rows-max">
				<!-- Product Details -->
				<Card>
					<CardHeader>
						<CardTitle>Детали товара</CardTitle>
						<CardDescription>Введите название и описание товара.</CardDescription>
					</CardHeader>
					<CardContent>
						<div class="gap-6 grid">
							<div class="gap-3 grid">
								<Label for="name">Название</Label>
								<Input
									id="name"
									type="text"
									class="w-full"
									v-model="productName"
									placeholder="Введите название продукта"
								/>
							</div>
							<div class="gap-3 grid">
								<Label for="description">Описание</Label>
								<Textarea
									id="description"
									v-model="productDescription"
									placeholder="Краткое описание продукта"
									class="min-h-32"
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
						<CardDescription>Загрузите изображение и видео для товара.</CardDescription>
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
							<div class="flex flex-col gap-3">
								<Label>Видео</Label>
								<Input
									type="file"
									accept="video/*"
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
										<SelectItem value="одежда">Одежда</SelectItem>
										<SelectItem value="электроника">Электроника</SelectItem>
										<SelectItem value="аксессуары">Аксессуары</SelectItem>
									</SelectContent>
								</Select>
							</div>
						</div>
					</CardContent>
				</Card>
			</div>
		</div>

		<div class="flex justify-center items-center gap-2 md:hidden">
			<Button variant="outline"> Отменить</Button>
			<Button>Сохранить</Button>
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
import { getRouteName } from '@/core/config/routes.config'
import { ChevronLeft } from 'lucide-vue-next'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const productName = ref('')
const productDescription = ref('')

const onBackClick = () => {
  router.push({name: getRouteName("ADMIN_PRODUCTS")})
}
</script>
