<template>
	<AdminFranchiseCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import AdminFranchiseCreateForm from '@/modules/admin/franchisees/components/create/admin-franchisee-create-form.vue'
import type { CreateFranchiseeDTO } from '@/modules/admin/franchisees/models/franchisee.model'
import { franchiseeService } from '@/modules/admin/franchisees/services/franchisee.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const createMutation = useMutation({
	mutationFn: (dto: CreateFranchiseeDTO) => franchiseeService.create(dto),
	onMutate: () => {
		toast({
			title: 'Создание...',
			description: 'Создание нового франчайзи. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-franchisees'] })
		toast({
			title: 'Успех!',
			description: 'Франчайзи успешно создан.',
		})
		router.back()
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при создании франчайзи.',
			variant: 'destructive',
		})
	},
})

function handleCreate(dto: CreateFranchiseeDTO) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
