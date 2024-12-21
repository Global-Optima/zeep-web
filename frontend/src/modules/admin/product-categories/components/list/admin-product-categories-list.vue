<template>
	<Table class="bg-white rounded-xl">
		<TableHeader>
			<TableRow>
				<TableHead class="p-4">Название</TableHead>
				<TableHead class="p-4">Описание</TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="productCategory in productCategories"
				:key="productCategory.id"
				class="h-12 cursor-pointer"
				@click="goToSupplierDetails(productCategory.id)"
			>
				<TableCell class="p-4">
					<span class="font-medium">{{ productCategory.name }}</span>
				</TableCell>

				<TableCell class="p-4">
					<span class="max-w-12 font-medium truncate">{{ productCategory.description }}</span>
				</TableCell>
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
import type { ProductCategories } from '@/modules/admin/product-categories/models/product-categories.model'
import { useRouter } from 'vue-router'

const {productCategories} = defineProps<{productCategories: ProductCategories[]}>()

const router = useRouter();

const goToSupplierDetails = (productCategoryId: number) => {
  router.push(`/admin/product-categories/${productCategoryId}`);
};
</script>
