<template>
	<Table>
		<TableHeader>
			<TableRow>
				<TableHead>Название</TableHead>
				<TableHead>Город</TableHead>
				<TableHead>Адресс</TableHead>
				<TableHead>Телефон</TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="supplier in suppliers"
				:key="supplier.id"
				@click="onSupplierClick(supplier.id)"
				class="hover:bg-slate-50 cursor-pointer"
			>
				<TableCell class="py-4 font-medium">{{ supplier.name }}</TableCell>
				<TableCell>{{ supplier.city }}</TableCell>
				<TableCell>{{ supplier.address }}</TableCell>
				<TableCell>{{ supplier.contactPhone }}</TableCell>
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
import type { SupplierDTO } from '@/modules/admin/suppliers/models/suppliers.model'
import { useRouter } from 'vue-router'

const {suppliers} = defineProps<{suppliers: SupplierDTO[]}>()

const router = useRouter();

const onSupplierClick = (id: number) => {
  router.push(`/admin/suppliers/${id}`);
};
</script>
