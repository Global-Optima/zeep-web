<template>
	<AdminProvisionCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
		:initialProvision="provisionDetails"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import { useAxiosLocaleToast, type AxiosLocalizedError } from '@/core/hooks/use-axios-locale-toast.hooks'
import AdminProvisionCreateForm from "@/modules/admin/provisions/components/create/admin-provision-create-form.vue"
import type { CreateProvisionDTO } from "@/modules/admin/provisions/models/provision.models"
import { provisionsService } from "@/modules/admin/provisions/services/provisions.service"
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const route = useRoute()
const queryClient = useQueryClient()
const { toast } = useToast()

const templateProvisionId = route.query.templateProvisionId as string

const { data: provisionDetails } = useQuery({
	queryKey: computed(() => ['admin-provision-details', templateProvisionId]),
	queryFn: () => provisionsService.getProvisionById(Number(templateProvisionId)),
	enabled: computed(() =>!isNaN(Number(templateProvisionId))),
})

const {toastLocalizedError} = useAxiosLocaleToast()

const createMutation = useMutation({
	mutationFn: (dto: CreateProvisionDTO) => provisionsService.createProvision(dto),
	onMutate: () => {
		toast({
			title: 'Создание...',
			description: 'Создание новой заготовки. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-provisions'] })
		toast({
			title: 'Успех!',
      variant: 'success',
			description: 'Заготовка успешна создана.',
		})
		router.back()
	},
	onError: (err: AxiosLocalizedError) => {
    toastLocalizedError(err, 'Произошла ошибка при создании заготовки.')
	},
})

function handleCreate(dto: CreateProvisionDTO) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
