<template>
	<Table class="bg-white rounded-xl">
		<TableHeader>
			<TableRow>
				<TableHead class="p-4">Название</TableHead>
				<TableHead class="p-4">Описание</TableHead>
				<TableHead class="p-4"></TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="productCategory in productCategories"
				:key="productCategory.id"
				class="h-12 cursor-pointer"
				@click="goToSupplierDetails(productCategory.id)"
			>
				<TableCell class="p-4">
					<span class="font-medium">{{ productCategory.name }}</span>
				</TableCell>

				<TableCell class="p-4">
					<span class="max-w-12 font-medium truncate">{{ productCategory.description }}</span>
				</TableCell>

				<TableCell class="flex justify-end">
					<Button
						variant="ghost"
						size="icon"
						@click="e => onDeleteClick(e, productCategory.id)"
					>
						<Trash class="w-6 h-6 text-red-400" />
					</Button>
				</TableCell>
			</TableRow>
		</TableBody>
	</Table>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
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
import type { ProductCategoryDTO } from '@/modules/kiosk/products/models/product.model'
import { productsService } from '@/modules/kiosk/products/services/products.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { Trash } from 'lucide-vue-next'
import { useRouter } from 'vue-router'

const {productCategories} = defineProps<{productCategories: ProductCategoryDTO[]}>()

const router = useRouter();
const queryClient = useQueryClient()
const { toastLocalizedError } = useAxiosLocaleToast()


const { mutate: deleteMutation } = useMutation({
	mutationFn: (id: number) => productsService.deleteProductCategory(id),
	onSuccess: () => {
		toast({ title: 'Успешное удаление' })
		queryClient.invalidateQueries({ queryKey: ['admin-product-categories'] })
	},
	onError: (error: AxiosLocalizedError) => {
		toastLocalizedError(error, 'Произошла ошибка при удалении' )
	},
})

const onDeleteClick = (e: Event, id: number) => {
	e.stopPropagation()

	const confirmed = window.confirm('Вы уверены, что хотите удалить?')
	if (confirmed) {
		deleteMutation(id)
	}
}


const goToSupplierDetails = (productCategoryId: number) => {
  router.push(`/admin/product-categories/${productCategoryId}`);
};
</script>
