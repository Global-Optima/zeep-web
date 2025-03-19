<template>
	<AdminUnitsCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminUnitsCreateForm from '@/modules/admin/units/components/create/admin-units-create-form.vue'
import type { CreateUnitDTO } from '@/modules/admin/units/models/units.model'
import { unitsService } from '@/modules/admin/units/services/units.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const createMutation = useMutation({
	mutationFn: (dto: CreateUnitDTO) => unitsService.createUnit(dto),
	onMutate: () => {
		toast({
			title: 'Создание...',
			description: 'Пожалуйста, подождите, создается новая единица измерения.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-units'] })
		toast({
			title: 'Успех!',
variant: 'success',
			description: 'Единица измерения успешно создана.',
		})
		router.back()
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при создании единицы измерения.',
			variant: 'destructive',
		})
	},
})

function handleCreate(dto: CreateUnitDTO) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
