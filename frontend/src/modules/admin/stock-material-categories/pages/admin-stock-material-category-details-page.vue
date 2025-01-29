<template>
	<p v-if="!categoryDetails">Категория не найдена</p>

	<AdminStockMaterialCategoryDetailsForm
		v-else
		:category="categoryDetails"
		@onSubmit="handleUpdate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminStockMaterialCategoryDetailsForm from '@/modules/admin/stock-material-categories/components/details/admin-stock-material-category-details-form.vue'
import type { UpdateStockMaterialCategoryDTO } from '@/modules/admin/stock-material-categories/models/stock-material-categories.model'
import { stockMaterialCategoryService } from '@/modules/admin/stock-material-categories/services/stock-materials.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const categoryId = route.params.id as string

const { data: categoryDetails } = useQuery({
	queryKey: computed(() => ['admin-stock-material-categories-details', categoryId]),
	queryFn: () => stockMaterialCategoryService.getById(Number(categoryId)),
	enabled: !isNaN(Number(categoryId)),
})

const updateMutation = useMutation({
	mutationFn: ({ id, dto }: { id: number; dto: UpdateStockMaterialCategoryDTO }) =>
		stockMaterialCategoryService.update(id, dto),
	onMutate: () => {
		toast({
			title: 'Обновление...',
			description: 'Обновление данных категории материалов. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-stock-material-categories'] })
		queryClient.invalidateQueries({ queryKey: ['admin-stock-material-categories-details', categoryId] })
		toast({
			title: 'Успех!',
			description: 'Данные категории материалов успешно обновлены.',
		})
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при обновлении категории материалов.',
			variant: 'destructive',
		})
	},
})

function handleUpdate(updatedData: UpdateStockMaterialCategoryDTO) {
	if (isNaN(Number(categoryId))) {
		toast({
			title: 'Ошибка',
			description: 'Неверный идентификатор категории.',
			variant: 'destructive',
		})
		return router.back()
	}

	updateMutation.mutate({ id: Number(categoryId), dto: updatedData })
}

function handleCancel() {
	router.back()
}
</script>
