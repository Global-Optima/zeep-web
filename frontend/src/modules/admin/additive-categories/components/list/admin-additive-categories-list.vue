<template>
	<Table class="bg-white rounded-xl">
		<TableHeader>
			<TableRow>
				<TableHead class="p-4">Название</TableHead>
				<TableHead class="p-4">Описание</TableHead>
				<TableHead class="p-4">Кол-во топпингов</TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="category in categories"
				:key="category.id"
				class="h-12 cursor-pointer"
				@click="goToSupplierDetails(category.id)"
			>
				<TableCell class="p-4 font-medium"> {{ category.name }} </TableCell>
				<TableCell> {{ category.description }} </TableCell>
				<TableCell> {{ category.additives.length }} </TableCell>
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
import type { AdditiveCategoryDTO } from '@/modules/admin/additives/models/additives.model'
import { useRouter } from 'vue-router'

const {categories} = defineProps<{categories: AdditiveCategoryDTO[]}>()

const router = useRouter();

const goToSupplierDetails = (productCategoryId: number) => {
  router.push(`/admin/additive-categories/${productCategoryId}`);
};
</script>
