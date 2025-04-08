<template>
	<AdminAdditiveCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
		:isSubmitting="isPending"
		:initialAdditive="additiveDetails"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminAdditiveCreateForm from '@/modules/admin/additives/components/create/admin-additive-create-form.vue'
import type { CreateAdditiveDTO } from '@/modules/admin/additives/models/additives.model'
import { additivesService } from '@/modules/admin/additives/services/additives.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const templateAdditiveId = route.query.templateAdditiveId as string

const { data: additiveDetails } = useQuery({
	queryKey: computed(() =>['admin-additive-details', templateAdditiveId]),
	queryFn: () => additivesService.getAdditiveById(Number(templateAdditiveId)),
	enabled: computed(() => !isNaN(Number(templateAdditiveId))),
})

const {mutate, isPending} = useMutation({
	mutationFn: (dto: CreateAdditiveDTO) => additivesService.createAdditive(dto),
	onMutate: () => {
		toast({
			title: 'Создание...',
			description: 'Создание новой модификатора. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-additives'] })
		toast({
			title: 'Успех!',
      variant: 'success',
			description: 'Модификатор успешно создан.',
		})
		router.back()
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при создании модификатора.',
			variant: 'destructive',
		})
	},
})

function handleCreate(dto: CreateAdditiveDTO) {
	mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
