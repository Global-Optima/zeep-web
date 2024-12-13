<template>
	<main class="flex-1 items-start gap-4 grid">
		<!-- Header -->
		<div class="flex justify-between items-center gap-2">
			<div class="flex items-center gap-2">
				<Input
					placeholder="Поиск добавок"
					class="bg-white w-full md:w-64"
				/>
				<DropdownMenu>
					<DropdownMenuTrigger as-child>
						<Button
							variant="outline"
							class="gap-2"
						>
							<ListFilter class="w-4 h-4" />
							<span class="sr-only sm:not-sr-only sm:whitespace-nowrap">Фильтр</span>
						</Button>
					</DropdownMenuTrigger>
					<DropdownMenuContent align="end">
						<DropdownMenuLabel>Фильтровать по</DropdownMenuLabel>
						<DropdownMenuSeparator />
						<DropdownMenuItem checked>Активные</DropdownMenuItem>
						<DropdownMenuItem>Черновики</DropdownMenuItem>
						<DropdownMenuItem>Архив</DropdownMenuItem>
					</DropdownMenuContent>
				</DropdownMenu>
			</div>

			<div class="flex items-center gap-2">
				<Button variant="outline">
					<span class="sr-only sm:not-sr-only sm:whitespace-nowrap">Экспорт</span>
				</Button>
				<Button>
					<span class="sr-only sm:not-sr-only sm:whitespace-nowrap">Добавить</span>
				</Button>
			</div>
		</div>

		<!-- Additives Table -->
		<Card>
			<CardContent class="mt-4">
				<Table>
					<TableHeader>
						<TableRow>
							<TableHead class="hidden w-[100px] sm:table-cell"></TableHead>
							<TableHead>Название</TableHead>
							<TableHead>Категория</TableHead>
							<TableHead>Базовая цена</TableHead>
							<TableHead>Размер</TableHead>
						</TableRow>
					</TableHeader>
					<TableBody>
						<TableRow
							v-for="additive in additives"
							:key="additive.id"
							@click="onAdditiveClick(additive.id)"
							class="hover:bg-slate-50 cursor-pointer"
						>
							<TableCell class="hidden sm:table-cell">
								<img
									:src="additive.imageUrl"
									alt="Изображение добавки"
									class="bg-gray-100 rounded-md aspect-square object-contain"
									height="64"
									width="64"
								/>
							</TableCell>
							<TableCell class="font-medium">{{ additive.name }}</TableCell>
							<TableCell>{{ additive.category }}</TableCell>
							<TableCell>{{ formatPrice(additive.basePrice) }}</TableCell>
							<TableCell>{{ additive.size }}</TableCell>
						</TableRow>
					</TableBody>
				</Table>
			</CardContent>
		</Card>
	</main>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { Card, CardContent } from '@/core/components/ui/card'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from '@/core/components/ui/dropdown-menu'
import { Input } from '@/core/components/ui/input'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/core/components/ui/table'
import { getRoutePath } from '@/core/config/routes.config'
import { formatPrice } from '@/core/utils/price.utils'
import { ListFilter } from 'lucide-vue-next'
import { ref } from 'vue'
import { useRouter } from 'vue-router'
const router = useRouter()

const onAdditiveClick = (additiveId: number) => {
  router.push(`${getRoutePath("ADMIN_ADDITIVES")}/${additiveId}`)
}

// Mock data for additives
const additives = ref([
  {
    id: 1,
    imageUrl: 'fake_url_1',
    name: 'Сахар',
    category: "Подсластители",
    basePrice: 50,
    size: '400 г',
  },
  {
    id: 2,
    imageUrl: 'fake_url_2',
    name: 'Миндальное молоко',
    category: "Молоко",
    basePrice: 70,
    size: '1 литр',
  },
  {
    id: 3,
    imageUrl: 'fake_url_3',
    name: 'Сироп карамельный',
    category: "Сиропы",
    basePrice: 120,
    size: '500 мл',
  },
  {
    id: 4,
    imageUrl: 'fake_url_4',
    name: 'Шоколадная стружка',
    category: "Посыпки",
    basePrice: 200,
    size: '300 г',
  },
  {
    id: 5,
    imageUrl: 'fake_url_5',
    name: 'Ванильный экстракт',
    category: "Сиропы",
    basePrice: 250,
    size: '100 мл',
  },
  {
    id: 6,
    imageUrl: 'fake_url_6',
    name: 'Корица',
    category: "Посыпки",
    basePrice: 180,
    size: '200 г',
  },
])
</script>

<style scoped>
/* Add any custom styles here */
</style>
