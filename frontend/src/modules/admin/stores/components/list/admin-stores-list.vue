<template>
	<Table class="bg-white rounded-xl">
		<TableHeader>
			<TableRow>
				<TableHead class="p-4">Название</TableHead>
				<TableHead class="p-4">Адрес</TableHead>
				<TableHead class="p-4">Телефон</TableHead>
				<TableHead class="p-4">Статус</TableHead>
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
						v-if="store.isFranchise"
						:class="[
                'inline-flex items-center rounded-full px-2.5 py-0.5 text-xs bg-green-100 text-green-80',
              ]"
					>
						Франшиза
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
import type { StoreDTO } from '@/modules/admin/stores/models/stores.models'
import { useRouter } from 'vue-router'

const {stores} = defineProps<{stores: StoreDTO[]}>()

const router = useRouter();


// Navigate to store details
const goToStore = (storeId: number) => {
  router.push(`/admin/stores/${storeId}`);
};

// Format phone number
const formatPhoneNumber = (phone: string) => {
  return phone.replace(/(\+7)(\d{3})(\d{3})(\d{2})(\d{2})/, '$1 ($2) $3-$4-$5');
};
</script>
