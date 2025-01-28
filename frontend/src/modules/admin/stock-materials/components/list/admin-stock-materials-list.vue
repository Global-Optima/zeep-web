<template>
	<Table>
		<TableHeader>
			<TableRow>
				<TableHead>Название</TableHead>
				<TableHead>Категория</TableHead>
				<TableHead>Ингредиент</TableHead>
				<TableHead>Упаковка</TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="stockMaterial in stockMaterials"
				:key="stockMaterial.id"
				@click="onStockMaterialClick(stockMaterial.id)"
				class="hover:bg-slate-50 cursor-pointer"
			>
				<TableCell class="py-4 font-medium">{{ stockMaterial.name }}</TableCell>
				<TableCell>{{ stockMaterial.category.name }}</TableCell>
				<TableCell>{{ stockMaterial.ingredient.name }}</TableCell>
				<TableCell>{{ stockMaterial.size }} {{ stockMaterial.unit.name }}</TableCell>
			</TableRow>
		</TableBody>
	</Table>
</template>

<script setup lang="ts">
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/core/components/ui/table'
import type { StockMaterialsDTO } from '@/modules/admin/stock-materials/models/stock-materials.model'
import { useRouter } from 'vue-router'

const {stockMaterials} = defineProps<{stockMaterials: StockMaterialsDTO[]}>()

const router = useRouter();

const onStockMaterialClick = (id: number) => {
  router.push(`/admin/stock-materials/${id}`);
};
</script>
