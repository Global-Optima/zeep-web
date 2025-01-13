<template>
	<Table class="bg-white rounded-xl">
		<TableHeader>
			<TableRow>
				<TableHead class="p-4">Название</TableHead>
				<TableHead class="p-4">Множитель</TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="unit in units"
				:key="unit.id"
				class="h-12 cursor-pointer"
				@click="goToSupplierDetails(unit.id)"
			>
				<TableCell class="p-4 font-medium"> {{ unit.name }} </TableCell>
				<TableCell> {{ unit.conversionFactor }} </TableCell>
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
import type { UnitDTO } from '@/modules/admin/units/models/units.model'
import { useRouter } from 'vue-router'

const {units} = defineProps<{units: UnitDTO[]}>()

const router = useRouter();

const goToSupplierDetails = (productCategoryId: number) => {
  router.push(`/admin/units/${productCategoryId}`);
};
</script>
