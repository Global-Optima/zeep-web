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

				<!-- Variants -->
				<Card>
					<CardHeader>
						<CardTitle>Варианты</CardTitle>
						<CardDescription>
							Добавьте размеры, базовые цены и отметьте стандартный вариант.
						</CardDescription>
					</CardHeader>
					<CardContent>
						<Table>
							<TableHeader>
								<TableRow>
									<TableHead>Название</TableHead>
									<TableHead>Размер</TableHead>
									<TableHead>Базовая цена</TableHead>
									<TableHead>По умолчанию</TableHead>
								</TableRow>
							</TableHeader>
							<TableBody>
								<TableRow
									v-for="(variant, index) in variants"
									:key="index"
								>
									<TableCell>
										<Select>
											<SelectTrigger>
												<SelectValue :placeholder="variant.name || 'Выберите название'" />
											</SelectTrigger>
											<SelectContent>
												<SelectItem value="S">S</SelectItem>
												<SelectItem value="M">M</SelectItem>
												<SelectItem value="L">L</SelectItem>
												<SelectItem value="XL">XL</SelectItem>
												<SelectItem value="XXL">XXL</SelectItem>
											</SelectContent>
										</Select>
									</TableCell>
									<TableCell>
										<Input
											v-model="variant.size"
											type="text"
											placeholder="Например: 500 мл"
										/>
									</TableCell>
									<TableCell>
										<Input
											v-model="variant.basePrice"
											type="number"
											placeholder="0.00"
										/>
									</TableCell>
									<TableCell>
										<Input
											type="radio"
											:checked="variant.isDefault"
											name="default-variant"
											class="bg-primary form-radio w-5 h-5"
											@change="setDefaultVariant(index)"
										/>
									</TableCell>
								</TableRow>
							</TableBody>
						</Table>
					</CardContent>
					<CardFooter class="justify-center p-4 border-t">
						<Button
							variant="ghost"
							class="gap-1"
							@click="addVariant"
						>
							<Plus class="w-4 h-4" />
							Добавить вариант
						</Button>
					</CardFooter>
				</Card>

				<Card>
					<CardHeader>
						<div class="flex justify-between items-start gap-4">
							<div>
								<CardTitle>Топпинги по умолчанию</CardTitle>
								<CardDescription class="mt-2">
									Добавьте топпинги, которые идут по умолчанию с товаром.
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

				<Card>
					<CardHeader>
						<div class="flex justify-between items-start gap-4">
							<div>
								<CardTitle>Техническая карта</CardTitle>
								<CardDescription class="mt-2">
									Выберите ингредиенты и укажите их вес для данного продукта.
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
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/core/components/ui/card'
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
import { Textarea } from '@/core/components/ui/textarea'
import { getRouteName } from '@/core/config/routes.config'
import { ChevronLeft, Plus } from 'lucide-vue-next'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const productName = ref('')
const productDescription = ref('')
const variants = ref([
  { name: 'M', size: '500 мл', basePrice: 299.99, isDefault: true },
  { name: 'L', size: '1 литр', basePrice: 399.99, isDefault: false },
])

const defaultAdditives = ref([
  { id: 1, name: 'Карамельный сироп', size: '500 мл', imageUrl: 'caramel.jpg' },
  { id: 2, name: 'Ванильный сироп', size: '500 мл', imageUrl: 'vanilla.jpg' },
  { id: 3, name: 'Шоколадный сироп', size: '500 мл', imageUrl: 'chocolate.jpg' },
])

const technicalCardIngredients = ref([
    {
        id: 1,
        imageUrl: 'ingredient1.jpg',
        name: 'Кофе',
        unit: 'грамм',
        weight: 30,
    },
    {
        id: 2,
        imageUrl: 'ingredient2.jpg',
        name: 'Молоко',
        unit: 'мл',
        weight: 100,
    },
]);


const addVariant = () => {
  variants.value.push({ name: '', size: '', basePrice: 0, isDefault: false })
}

const setDefaultVariant = (index: number) => {
  variants.value.forEach((variant, i) => {
    variant.isDefault = i === index
  })
}

const onBackClick = () => {
  router.push({name: getRouteName("ADMIN_PRODUCTS")})
}
</script>
