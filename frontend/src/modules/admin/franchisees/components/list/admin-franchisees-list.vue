<template>
	<Table>
		<TableHeader>
			<TableRow>
				<TableHead>Название</TableHead>
				<TableHead>Описание</TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="franchisee in franchisees"
				:key="franchisee.id"
				@click="onAdditiveClick(franchisee.id)"
				class="hover:bg-slate-50 cursor-pointer"
			>
				<TableCell class="py-4 font-medium">{{ franchisee.name }}</TableCell>
				<TableCell
					class="line-clamp-1 max-w-lg"
					>{{ franchisee.description ?? "Отсутсвует" }}</TableCell
				>
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
import type { FranchiseeDTO } from '@/modules/admin/franchisees/models/franchisee.model'
import { useRouter } from 'vue-router'

const {franchisees} = defineProps<{franchisees: FranchiseeDTO[]}>()

const router = useRouter();

const onAdditiveClick = (additiveID: number) => {
  router.push(`/admin/franchisees/${additiveID}`);
};
</script>
