<template>
	<Table>
		<TableHeader>
			<TableRow>
				<TableHead>Название</TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="region in regions"
				:key="region.id"
				@click="onAdditiveClick(region.id)"
				class="hover:bg-slate-50 cursor-pointer"
			>
				<TableCell class="py-4 font-medium">{{ region.name }}</TableCell>
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
import type { RegionDTO } from '@/modules/admin/regions/models/regions.model'
import { useRouter } from 'vue-router'

const {regions} = defineProps<{regions: RegionDTO[]}>()

const router = useRouter();

const onAdditiveClick = (additiveID: number) => {
  router.push(`/admin/regions/${additiveID}`);
};
</script>
