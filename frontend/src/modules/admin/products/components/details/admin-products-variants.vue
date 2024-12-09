<template>
	<main class="flex-1 items-start gap-4 grid">
		<div class="flex justify-between items-center gap-2">
			<div class="flex items-center gap-2">
				<Input
					placeholder="Поиск вариантов"
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

			<Button>
				<span class="sr-only sm:not-sr-only sm:whitespace-nowrap">Добавить</span>
			</Button>
		</div>

		<Card>
			<CardContent class="mt-4">
				<Table>
					<TableHeader>
						<TableRow>
							<TableHead>Название</TableHead>
							<TableHead>Размер</TableHead>
							<TableHead>Начальная цена</TableHead>
							<TableHead>По умолчанию</TableHead>
						</TableRow>
					</TableHeader>
					<TableBody>
						<TableRow
							v-for="variant in variants"
							:key="variant.id"
							@click="onVariantClick(variant.id)"
							class="hover:bg-slate-50"
						>
							<TableCell class="font-medium py-4">{{ variant.name }}</TableCell>
							<TableCell>{{ variant.size }}</TableCell>
							<TableCell>
								{{ variant.basePrice }}
							</TableCell>
							<TableCell>
								<input
									type="checkbox"
									:checked="variant.isDefault"
									@click="setDefaultVariant(variant.id)"
								/>
							</TableCell>
						</TableRow>
					</TableBody>
				</Table>
			</CardContent>
		</Card>
	</main>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button';
import { Card, CardContent } from '@/core/components/ui/card';
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from '@/core/components/ui/table';
import { Input } from '@/core/components/ui/input';
import { ref } from 'vue';
import { ListFilter } from "lucide-vue-next";
import { DropdownMenu, DropdownMenuTrigger, DropdownMenuContent,
DropdownMenuLabel,
DropdownMenuSeparator,
DropdownMenuItem } from "@/core/components/ui/dropdown-menu";
import { useRouter } from "vue-router";

const router = useRouter()

// Variants array
const variants = ref([
	{ id: 1, name: 'S', size: '250 мл', basePrice: 120, isDefault: false },
	{ id: 2, name: 'M', size: '350 мл', basePrice: 150, isDefault: true },
	{ id: 3, name: 'L', size: '500 мл', basePrice: 180, isDefault: false },
	{ id: 4, name: 'XL', size: '700 мл', basePrice: 220, isDefault: false },
]);

// Set default variant
const setDefaultVariant = (id: number) => {
	variants.value = variants.value.map((variant) => ({
		...variant,
		isDefault: variant.id === id, // Mark only the selected variant as default
	}));
};

const onVariantClick = (variantId: number) => {
	router.push(`../product-sizes/${variantId}`)
}
</script>

<style scoped>
/* Add custom styles if needed */
</style>
