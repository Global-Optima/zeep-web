<template>
	<Table>
		<TableHeader>
			<TableRow>
				<TableHead class="hidden w-[100px] sm:table-cell"></TableHead>
				<TableHead>Название</TableHead>
				<TableHead>Категория</TableHead>
				<TableHead>Базовая цена</TableHead>
				<TableHead>Размер</TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="additive in additives"
				:key="additive.id"
				@click="onAdditiveClick(additive.id)"
				class="hover:bg-slate-50 cursor-pointer"
			>
				<TableCell class="hidden sm:table-cell">
					<img
						:src="additive.imageUrl"
						alt="Изображение добавки"
						class="bg-gray-100 rounded-md aspect-square object-contain"
						height="64"
						width="64"
					/>
				</TableCell>
				<TableCell class="font-medium">{{ additive.name }}</TableCell>
				<TableCell>{{ additive.category.name }}</TableCell>
				<TableCell>{{ formatPrice(additive.price) }}</TableCell>
				<TableCell>{{ additive.size }}</TableCell>
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
import { formatPrice } from '@/core/utils/price.utils'
import type { Additives } from '@/modules/admin/additives/models/additives.model'
import { useRouter } from 'vue-router'

const {additives} = defineProps<{additives: Additives[]}>()

const router = useRouter();

const onAdditiveClick = (additiveID: number) => {
  router.push(`/admin/additives/${additiveID}`);
};
</script>
