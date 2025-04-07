<template>
	<Card>
		<CardHeader>
			<CardTitle>Список ингредиентов</CardTitle>
			<CardDescription> Ингредиенты, необходимые для заготовки </CardDescription>
		</CardHeader>

		<CardContent>
			<Table>
				<TableHeader>
					<TableRow>
						<TableHead>Ингредиент</TableHead>
						<TableHead>Категория</TableHead>
						<TableHead>Объем</TableHead>
					</TableRow>
				</TableHeader>
				<TableBody>
					<TableRow
						v-for="(item, index) in storeProvision.ingredients"
						:key="index"
					>
						<TableCell>{{ item.ingredient.name }}</TableCell>
						<TableCell>{{ item.ingredient.category }}</TableCell>
						<TableCell>{{ item.volume }} {{ item.ingredient.unit.name.toLowerCase() }}</TableCell>
					</TableRow>
				</TableBody>
			</Table>
		</CardContent>
	</Card>
</template>

<script setup lang="ts">
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/core/components/ui/card'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/core/components/ui/table'
import {
  type StoreProvisionDetailsDTO
} from '@/modules/admin/store-provisions/models/store-provision.models'

const {storeProvision} = defineProps<{
  storeProvision: StoreProvisionDetailsDTO
}>()
</script>
