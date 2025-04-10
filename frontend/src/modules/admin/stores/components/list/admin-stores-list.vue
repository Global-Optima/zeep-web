<template>
	<Table class="bg-white rounded-xl">
		<TableHeader>
			<TableRow>
				<TableHead class="p-4">Название</TableHead>
				<TableHead class="p-4">Адрес</TableHead>
				<TableHead class="p-4">Телефон</TableHead>
				<TableHead class="p-4">Тип</TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="store in stores"
				:key="store.id"
				class="h-12 cursor-pointer"
				@click="goToStore(store.id)"
			>
				<!-- Store Name -->
				<TableCell class="p-4">
					<span class="font-medium">{{ store.name }}</span>
				</TableCell>
				<!-- Address -->
				<TableCell class="p-4">
					{{ store.facilityAddress.address }}
				</TableCell>
				<!-- Phone -->
				<TableCell class="p-4"> {{ formatPhoneNumber(store.contactPhone) }} </TableCell>
				<!-- Status -->
				<TableCell class="p-4">
					<span
						:class="[
              'inline-flex items-center rounded-full px-2.5 py-0.5 text-xs',
              store.franchisee ? 'bg-green-100 text-green-800' : 'bg-blue-100 text-blue-800'
              ]"
					>
						{{ store.franchisee ? 'Франшиза' : 'Официальное' }}
					</span>
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
import { formatPhoneNumber } from '@/core/utils/fomat-phone-number.utils'
import type { StoreDTO } from '@/modules/admin/stores/models/stores.models'
import { useRouter } from 'vue-router'

const {stores} = defineProps<{stores: StoreDTO[]}>()

const router = useRouter();


// Navigate to store details
const goToStore = (storeId: number) => {
  router.push(`/admin/stores/${storeId}`);
};
</script>
