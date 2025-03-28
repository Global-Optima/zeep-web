<template>
	<AdminRegionEmployeesCreateForm
		@on-submit="handleCreate"
		@on-cancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { useToast } from '@/core/components/ui/toast/use-toast'
import type { CreateEmployeeDTO } from '@/modules/admin/employees/models/employees.models'
import AdminRegionEmployeesCreateForm from '@/modules/admin/employees/regions/components/create/admin-region-employees-create-form.vue'
import { regionEmployeeService } from '@/modules/admin/employees/regions/services/region-employees.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRoute, useRouter } from 'vue-router'
import {type AxiosLocalizedError, useAxiosLocaleToast} from "@/core/hooks/use-axios-locale-toast.hooks";

const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()
const { toastLocalizedError } = useAxiosLocaleToast()
const route = useRoute()
const regionId = route.params.regionId as string

const createMutation = useMutation({
	mutationFn: ({ dto, regionId }: { dto: CreateEmployeeDTO, regionId: number }) =>
		regionEmployeeService.createRegionEmployee(dto, regionId),
	onMutate: () => {
		toast({
			title: 'Добавление сотрудника...',
			description: 'Мы создаем нового сотрудника. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['region-employees'] })
		router.back()
		toast({
			title: 'Сотрудник добавлен!',
			description: 'Новый сотрудник успешно зарегистрирован в системе.',
		})
	},
  onError: (error: AxiosLocalizedError) => {
    toastLocalizedError(error, "Произошла ошибка при обновлении материалов запроса на склад.")
  }
})

function handleCreate(dto: CreateEmployeeDTO) {
	if (!regionId) {
		toast({
			title: 'Ошибка',
			description: 'Не удалось определить регион для сотрудника.',
			variant: 'destructive',
		})
		return
	}

	createMutation.mutate({ dto, regionId: Number(regionId) })
}

function handleCancel() {
	router.back()
}
</script>
