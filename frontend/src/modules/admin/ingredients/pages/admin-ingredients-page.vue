<template>
	<main class="flex-1 items-start gap-4 grid">
		<!-- Header -->
		<div class="flex justify-between items-center gap-2">
			<div class="flex items-center gap-2">
				<Input
					placeholder="Поиск ингредиентов"
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
						<DropdownMenuItem checked>Доступные</DropdownMenuItem>
						<DropdownMenuItem>Скрытые</DropdownMenuItem>
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

		<!-- Ingredients Table -->
		<Card>
			<CardContent class="mt-4">
				<Table>
					<TableHeader>
						<TableRow>
							<TableHead class="hidden w-[100px] sm:table-cell"></TableHead>
							<TableHead>Название</TableHead>
						</TableRow>
					</TableHeader>
					<TableBody>
						<TableRow
							v-for="ingredient in ingredients"
							:key="ingredient.id"
							@click="onIngredientClick(ingredient.id)"
							class="hover:bg-slate-50 cursor-pointer"
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
import { ListFilter } from 'lucide-vue-next'
import { ref } from 'vue'
import { useRouter } from 'vue-router'
const router = useRouter()

const onIngredientClick = (ingredientId: number) => {
  router.push(`${getRoutePath("ADMIN_INGREDIENTS")}/${ingredientId}`)
}

// Mock data for ingredients
const ingredients = ref([
  {
    id: 1,
    imageUrl: 'fake_url_1',
    name: 'Мука',
  },
  {
    id: 2,
    imageUrl: 'fake_url_2',
    name: 'Яйца',
  },
  {
    id: 3,
    imageUrl: 'fake_url_3',
    name: 'Соль',
  },
  {
    id: 4,
    imageUrl: 'fake_url_4',
    name: 'Масло сливочное',
  },
])
</script>

<style scoped>
/* Add any custom styles here */
</style>
