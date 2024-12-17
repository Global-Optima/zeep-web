<template>
	<div>
		<p
			v-if="suppliers.length === 0"
			class="text-muted-foreground"
		>
			Поставщики не найдены
		</p>

		<Table
			v-else
			class="bg-white rounded-xl"
		>
			<TableHeader>
				<TableRow>
					<TableHead class="p-4">Название</TableHead>
					<TableHead class="p-4">Адрес</TableHead>
					<TableHead class="p-4">Телефон</TableHead>
					<TableHead class="p-4">Email</TableHead>
				</TableRow>
			</TableHeader>
			<TableBody>
				<TableRow
					v-for="supplier in suppliers"
					:key="supplier.id"
					class="h-12 cursor-pointer"
					@click="goToSupplierDetails(supplier.id)"
				>
					<TableCell class="p-4">
						<span class="font-medium">{{ supplier.name }}</span>
					</TableCell>

					<TableCell class="p-4">
						{{ supplier.address }}
					</TableCell>

					<TableCell class="p-4"> {{ formatPhoneNumber(supplier.contactPhone) }} </TableCell>
					<TableCell class="p-4"> {{ formatPhoneNumber(supplier.contactEmail) }} </TableCell>
				</TableRow>
			</TableBody>
		</Table>
	</div>
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
import type { Suppliers } from '@/modules/admin/suppliers/models/suppliers.model'
import { useRouter } from 'vue-router'

const {suppliers} = defineProps<{suppliers: Suppliers[]}>()

const router = useRouter();

const goToSupplierDetails = (storeId: number) => {
  router.push(`/admin/suppliers/${storeId}`);
};

const formatPhoneNumber = (phone: string) => {
  return phone.replace(/(\+7)(\d{3})(\d{3})(\d{2})(\d{2})/, '$1 ($2) $3-$4-$5');
};
</script>
