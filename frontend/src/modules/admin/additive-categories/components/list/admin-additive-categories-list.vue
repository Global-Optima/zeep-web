<template>
	<Table class="bg-white rounded-xl">
		<TableHeader>
			<TableRow>
				<TableHead class="p-4">Название</TableHead>
				<TableHead class="p-4">Описание</TableHead>
				<TableHead class="p-4">Топпингов</TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="additiveCategory in additiveCategories"
				:key="additiveCategory.id"
				class="h-12 cursor-pointer"
				@click="goToSupplierDetails(additiveCategory.id)"
			>
				<TableCell class="p-4">
					<span class="font-medium">{{ additiveCategory.name }}</span>
				</TableCell>

				<TableCell class="p-4">
					<span>{{ additiveCategory.description }}</span>
				</TableCell>

				<TableCell class="p-4">
					<span>{{ additiveCategory.additives.length }}</span>
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
import type { AdditiveCategories } from '@/modules/admin/additives/models/additives.model'
import { useRouter } from 'vue-router'

const {additiveCategories} = defineProps<{additiveCategories: AdditiveCategories[]}>()

const router = useRouter();

const goToSupplierDetails = (productCategoryId: number) => {
  router.push(`/admin/additive-categories/${productCategoryId}`);
};
</script>
