<template>
	<Table>
		<TableHeader>
			<TableRow>
				<TableHead>Название</TableHead>
				<TableHead>Регион</TableHead>
				<TableHead>Адрес</TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="warehouse in warehouses"
				:key="warehouse.id"
				@click="onAdditiveClick(warehouse.id)"
				class="hover:bg-slate-50 cursor-pointer"
			>
				<TableCell class="py-4 font-medium">{{ warehouse.name }}</TableCell>
				<TableCell>{{ warehouse.region.name }}</TableCell>
				<TableCell>{{ warehouse.facilityAddress.address }}</TableCell>
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
import type { WarehouseDTO } from '@/modules/admin/warehouses/models/warehouse.model'
import { useRouter } from 'vue-router'
const {warehouses} = defineProps<{warehouses: WarehouseDTO[]}>()

const router = useRouter();

const onAdditiveClick = (additiveID: number) => {
  router.push(`/admin/warehouses/${additiveID}`);
};
</script>
