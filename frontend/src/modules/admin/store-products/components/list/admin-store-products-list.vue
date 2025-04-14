<template>
	<Table>
		<TableHeader>
			<TableRow>
				<TableHead class="hidden sm:table-cell w-[100px]"> </TableHead>
				<TableHead>Название</TableHead>
				<TableHead class="hidden md:table-cell">Категория</TableHead>
				<TableHead>Цена</TableHead>
				<TableHead class="hidden md:table-cell">Статус</TableHead>
				<TableHead v-if="canDelete"></TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="product in storeProducts"
				:key="product.id"
				class="hover:bg-gray-50 h-12 cursor-pointer"
				@click="onProductClick(product.id)"
			>
				<TableCell class="hidden sm:table-cell">
					<LazyImage
						:src="product.imageUrl"
						alt="Изображение продукта"
						class="rounded-md size-16 object-contain aspect-square"
					/>
				</TableCell>
				<TableCell class="font-medium">
					{{ product.name }}
				</TableCell>
				<TableCell class="hidden md:table-cell">
					{{ product.category.name }}
				</TableCell>
				<TableCell>
					{{ formatPrice(product.storePrice) }}
				</TableCell>
				<TableCell class="hidden md:table-cell">
					<p
						class="inline-flex items-center px-2.5 py-1 rounded-md w-fit text-xs"
						:class="getStatusClass(product)"
					>
						{{ getStatusLabel(product) }}
					</p>
				</TableCell>
				<TableCell
					class="flex justify-end"
					v-if="canDelete"
				>
					<Button
						variant="ghost"
						size="icon"
						@click="e => onDeleteProductClick(e, product.id)"
					>
						<Trash class="w-6 h-6 text-red-400" />
					</Button>
				</TableCell>
			</TableRow>
		</TableBody>
	</Table>
</template>

<script setup lang="ts">
import LazyImage from '@/core/components/lazy-image/LazyImage.vue'
import Button from '@/core/components/ui/button/Button.vue'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/core/components/ui/table'
import { toast } from '@/core/components/ui/toast'
import { useAxiosLocaleToast, type AxiosLocalizedError } from '@/core/hooks/use-axios-locale-toast.hooks'
import { useHasRole } from '@/core/hooks/use-has-roles.hook'
import { formatPrice } from '@/core/utils/price.utils'
import { EmployeeRole } from '@/modules/admin/employees/models/employees.models'
import type { StoreProductDTO } from '@/modules/admin/store-products/models/store-products.model'
import { storeProductsService } from '@/modules/admin/store-products/services/store-products.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { Trash } from 'lucide-vue-next'
import { useRouter } from 'vue-router'

type StoreProductStatus = 'available' | 'out_of_stock' | 'unavailable'

const STORE_PRODUCT_STATUS_COLOR: Record<StoreProductStatus, string> = {
  available: 'bg-green-100 text-green-800',
  unavailable: 'bg-yellow-100 text-yellow-800',
  out_of_stock: 'bg-red-100 text-red-800',
}

const STORE_PRODUCT_STATUS_FORMATTED: Record<StoreProductStatus, string> = {
  available: 'Доступен',
  unavailable: 'Недоступен',
  out_of_stock: 'Нет в наличии',
}

function computeStatus(storeProduct: StoreProductDTO): StoreProductStatus {
  if (!storeProduct.isAvailable) {
    return 'unavailable'
  }
  if (storeProduct.isOutOfStock) {
    return 'out_of_stock'
  }
  return 'available'
}

function getStatusClass(storeProduct: StoreProductDTO): string {
  return STORE_PRODUCT_STATUS_COLOR[computeStatus(storeProduct)]
}

function getStatusLabel(storeProduct: StoreProductDTO): string {
  return STORE_PRODUCT_STATUS_FORMATTED[computeStatus(storeProduct)]
}

const router = useRouter();
const queryClient = useQueryClient();
const { toastLocalizedError } = useAxiosLocaleToast()

const { storeProducts } = defineProps<{ storeProducts: StoreProductDTO[] }>();

const canDelete = useHasRole([EmployeeRole.STORE_MANAGER])

const { mutate: deleteStoreProduct } = useMutation({
	mutationFn: (id: number) => storeProductsService.deleteStoreProduct(id),
	onSuccess: () => {
		toast({ title: "Продукт удален из кафе" });
		queryClient.invalidateQueries({ queryKey: ['admin-store-products'] });
	},
	onError: (error: AxiosLocalizedError) => {
		toastLocalizedError(error, "Произошла ошибка при удалении.");
	},
});

const onProductClick = (storeProductId: number) => {
	router.push(`/admin/store-products/${storeProductId}`);
};

const onDeleteProductClick = (e: Event, id: number) => {
  if (!canDelete) return
	e.stopPropagation();
	deleteStoreProduct(id);
};
</script>
