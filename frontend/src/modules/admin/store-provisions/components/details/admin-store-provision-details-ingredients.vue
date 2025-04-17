<template>
	<Card>
		<CardHeader>
			<CardTitle>Технологическая карта</CardTitle>
			<CardDescription> Сырье, необходимое для заготовки </CardDescription>
		</CardHeader>

		<CardContent>
			<Table>
				<TableHeader>
					<TableRow>
						<TableHead>Название</TableHead>
						<TableHead>Категория</TableHead>
						<TableHead>Объем</TableHead>
					</TableRow>
				</TableHeader>
				<TableBody>
					<TableRow
						v-for="(item, index) in storeProvision.ingredients"
						:key="index"
					>
						<TableCell>
							<RouterLink
								:to="{name: getRouteName('ADMIN_INGREDIENTS_DETAILS'), params: {id: item.ingredient.id}}"
								target="_blank"
								class="hover:text-primary underline transition-colors duration-300 underline-offset-4"
							>
								{{ item.ingredient.name }}
							</RouterLink>
						</TableCell>
						<TableCell>{{ item.ingredient.category.name }}</TableCell>
						<TableCell>{{ item.quantity }} {{ item.ingredient.unit.name.toLowerCase() }}</TableCell>
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
import { getRouteName } from "@/core/config/routes.config";
import {
  type StoreProvisionDetailsDTO
} from '@/modules/admin/store-provisions/models/store-provision.models'

const {storeProvision} = defineProps<{
  storeProvision: StoreProvisionDetailsDTO
}>()
</script>
