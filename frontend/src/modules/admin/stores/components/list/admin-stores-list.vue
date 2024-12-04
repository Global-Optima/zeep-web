<template>
	<div>
		<!-- If no stores, display message -->
		<p
			v-if="stores.length === 0"
			class="text-muted-foreground"
		>
			Магазины не найдены
		</p>
		<!-- If there are stores, display the table -->
		<Table
			v-else
			class="bg-white rounded-xl"
		>
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
					v-for="store in filteredStores"
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
						{{ store.address }}
					</TableCell>
					<!-- Phone -->
					<TableCell class="p-4">
						{{ formatPhoneNumber(store.phone) }}
					</TableCell>
					<!-- Status -->
					<TableCell class="p-4">
						<span
							:class="[
                'inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium',
                STATUS_COLOR[store.status],
              ]"
						>
							{{ STATUS_FORMATTED[store.status] }}
						</span>
					</TableCell>
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
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

// Mock data for stores
interface Store {
  id: number;
  name: string;
  address: string;
  phone: string;
  status: 'active' | 'inactive';
}

const stores = ref<Store[]>([
  {
    id: 1,
    name: 'Магазин Центральный',
    address: 'ул. Ленина, 10, Москва',
    phone: '+74951234567',
    status: 'active',
  },
  {
    id: 2,
    name: 'Магазин Северный',
    address: 'ул. Садовая, 5, Санкт-Петербург',
    phone: '+78121234567',
    status: 'inactive',
  },
  // Add more stores as needed
]);

const router = useRouter();



// Computed property for filtered stores
const filteredStores = computed(() => {
  const filtered = stores.value;


  return filtered;
});

// Navigate to store details
const goToStore = (storeId: number) => {
  router.push(`/admin/stores/${storeId}`);
};

// Format phone number
const formatPhoneNumber = (phone: string) => {
  return phone.replace(/(\+7)(\d{3})(\d{3})(\d{2})(\d{2})/, '$1 ($2) $3-$4-$5');
};

// Status colors and labels
const STATUS_COLOR: Record<string, string> = {
  active: 'bg-green-100 text-green-800',
  inactive: 'bg-red-100 text-red-800',
};

const STATUS_FORMATTED: Record<string, string> = {
  active: 'Активен',
  inactive: 'Неактивен',
};
</script>
