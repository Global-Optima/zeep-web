<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import { FormField, FormItem, FormLabel, FormMessage } from '@/core/components/ui/form'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/core/components/ui/select'
import { useToast } from '@/core/components/ui/toast'
import { getEmployeeShortName } from '@/core/utils/user-formatting.utils'
import { EMPLOYEE_ROLES_FORMATTED, EMPLOYEE_TYPES_FORMATTED, EmployeeRole, EmployeeType } from '@/modules/admin/employees/models/employees.models'
import { employeesService } from '@/modules/admin/employees/services/employees.service'
import { useMutation, useQuery } from '@tanstack/vue-query'
import { toTypedSchema } from '@vee-validate/zod'
import { ChevronDown, ChevronLeft } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { computed, defineAsyncComponent, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import * as z from 'zod'

// Route and Router
const route = useRoute()
const router = useRouter()

const { toast } = useToast()

// Reactive State
const openWorkplaceDialog = ref(false)
const selectedWorkplace = ref<{ id: number; name: string } | null>(null)
const selectedType = ref<EmployeeType | null>(null)
const selectedRole = ref<EmployeeRole | null>(null)

// Fetch Current Employee
const { data: employee } = useQuery({
  queryKey: ['employee', route.params.id],
  queryFn: () => employeesService.getEmployeeById(Number(route.params.id)),
})

// Zod Validation Schema
const schema = toTypedSchema(
  z.object({
    employeeType: z.nativeEnum(EmployeeType),
    role: z.nativeEnum(EmployeeRole),
    workplaceId: z.number().min(1, 'Выберите место работы'),
  })
)

// Form Setup with Vee-Validate
const { handleSubmit, setFieldValue } = useForm({
  validationSchema: schema,
})


const employeeTypes: EmployeeType[] = [
	EmployeeType.STORE,
	EmployeeType.WAREHOUSE,
	EmployeeType.FRANCHISEE,
	EmployeeType.REGION,
]

// Available Roles Based on Selected Type
const availableRoles = computed(() => {
  if (!selectedType.value) return []
  switch (selectedType.value) {
    case EmployeeType.STORE:
      return [EmployeeRole.STORE_MANAGER, EmployeeRole.BARISTA]
    case EmployeeType.WAREHOUSE:
      return [EmployeeRole.WAREHOUSE_MANAGER, EmployeeRole.WAREHOUSE_EMPLOYEE]
    case EmployeeType.FRANCHISEE:
      return [EmployeeRole.FRANCHISEE_MANAGER, EmployeeRole.FRANCHISEE_OWNER]
    case EmployeeType.REGION:
      return [EmployeeRole.REGION_WAREHOUSE_MANAGER]
    default:
      return []
  }
})

// Workplace Dialog Component Based on Selected Type
const workplaceDialogComponent = computed(() => {
  switch (selectedType.value) {
    case EmployeeType.STORE:
      return defineAsyncComponent(() => import('@/modules/admin/stores/components/admin-select-store-dialog.vue'))
    case EmployeeType.WAREHOUSE:
      return defineAsyncComponent(() => import('@/modules/admin/warehouses/components/admin-select-warehouse-dialog.vue'))
    case EmployeeType.FRANCHISEE:
      return defineAsyncComponent(() => import('@/modules/admin/franchisees/components/admin-select-franchisee-dialog.vue'))
    case EmployeeType.REGION:
      return defineAsyncComponent(() => import('@/modules/admin/regions/components/admin-select-region-dialog.vue'))
    default:
      return null
  }
})

// Handle Type Change
const onTypeChange = (type: EmployeeType) => {
  selectedType.value = type
  selectedRole.value = null
  setFieldValue('role', undefined)
  selectedWorkplace.value = null
  setFieldValue('employeeType', type)
}

// Handle Role Change
const onRoleChange = (role: EmployeeRole) => {
  selectedRole.value = role
  setFieldValue('role', role)
}

// Handle Workplace Selection
const selectWorkplace = (workplace: { id: number; name: string }) => {
  selectedWorkplace.value = workplace
  openWorkplaceDialog.value = false
  setFieldValue('workplaceId', workplace.id)
}

// Handle Cancel
const handleCancel = () => {
  router.back()
}

// Mutation for Reassigning Employee Type
const { mutate: reassignEmployeeType, isPending } = useMutation({
  mutationFn: (data: { employeeType: EmployeeType; role: EmployeeRole; workplaceId: number }) =>
    employeesService.reassignEmployeeType(Number(route.params.id), data),
  onSuccess: () => {
    toast({
      title: 'Успех',
      description: 'Сотрудник успешно переназначен.',
      variant: 'default',
    })
    router.back()
  },
  onError: (error) => {
    toast({
      title: 'Ошибка',
      description: 'Не удалось переназначить сотрудника.',
      variant: 'destructive',
    })
    console.error('Ошибка при переназначении сотрудника:', error)
  },
})

// Handle Submit
const onSubmit = handleSubmit((values) => {
  if (!selectedType.value || !selectedRole.value || !selectedWorkplace.value) {
    toast({
      title: 'Ошибка валидации',
      description: 'Пожалуйста, заполните все обязательные поля.',
      variant: 'destructive',
    })

    return
  }
  reassignEmployeeType({
    employeeType: values.employeeType,
    role: values.role,
    workplaceId: values.workplaceId,
  })
})
</script>

<template>
	<div
		v-if="employee"
		class="flex-1 gap-4 grid auto-rows-max mx-auto max-w-4xl"
	>
		<!-- Header -->
		<div class="flex items-center gap-4">
			<Button
				variant="outline"
				size="icon"
				@click="handleCancel"
			>
				<ChevronLeft class="w-5 h-5" />
				<span class="sr-only">Назад</span>
			</Button>
			<h1 class="flex-1 sm:grow-0 font-semibold text-xl tracking-tight whitespace-nowrap shrink-0">
				Переназначение {{ getEmployeeShortName(employee) }}
			</h1>
			<div class="hidden md:flex items-center gap-2 md:ml-auto">
				<Button
					variant="outline"
					type="button"
					@click="handleCancel"
				>
					Отменить
				</Button>
				<Button
					type="submit"
					form="reassignForm"
					:disabled="isPending"
				>
					Сохранить
				</Button>
			</div>
		</div>
		<!-- Main Content -->
		<Card>
			<CardHeader>
				<CardTitle>Переназначение сотрудника</CardTitle>
				<CardDescription>Укажите новый тип, роль и место работы для сотрудника.</CardDescription>
			</CardHeader>
			<CardContent>
				<!-- Form -->
				<form
					id="reassignForm"
					@submit.prevent="onSubmit"
					class="space-y-6"
				>
					<!-- Employee Type -->
					<FormField
						name="employeeType"
						v-slot="{ componentField }"
					>
						<FormItem>
							<FormLabel>Тип сотрудника</FormLabel>
							<Select
								v-bind="componentField"
								@update:model-value="(type) => onTypeChange(type as EmployeeType)"
							>
								<SelectTrigger>
									<SelectValue placeholder="Выберите тип" />
								</SelectTrigger>
								<SelectContent>
									<SelectItem
										v-for="type in employeeTypes"
										:key="type"
										:value="type"
									>
										{{ EMPLOYEE_TYPES_FORMATTED[type] }}
									</SelectItem>
								</SelectContent>
							</Select>
							<FormMessage />
						</FormItem>
					</FormField>
					<!-- Role -->
					<FormField
					 	:key="selectedType"
						name="role"
						v-if="selectedType"
						v-slot="{ componentField }"
					>
						<FormItem>
							<FormLabel>Роль</FormLabel>
							<Select
								v-bind="componentField"
								@update:model-value="(role) => onRoleChange(role as EmployeeRole)"
							>
								<SelectTrigger>
									<SelectValue placeholder="Выберите роль" />
								</SelectTrigger>
								<SelectContent>
									<SelectItem
										v-for="role in availableRoles"
										:key="role"
										:value="role"
									>
										{{ EMPLOYEE_ROLES_FORMATTED[role] }}
									</SelectItem>
								</SelectContent>
							</Select>
							<FormMessage />
						</FormItem>
					</FormField>
					<!-- Workplace -->
					<FormField
						name="workplaceId"
						v-if="selectedRole"
					>
						<FormItem>
							<FormLabel>Место работы</FormLabel>
							<div
								@click="openWorkplaceDialog = true"
								class="flex justify-between items-center gap-4 px-3 py-2 border rounded-md text-sm cursor-pointer"
							>
								<span>{{ selectedWorkplace?.name || 'Выберите место работы' }}</span>
								<ChevronDown class="size-5 text-gray-400" />
							</div>
							<FormMessage />
						</FormItem>
					</FormField>
					<!-- Dialog for Selecting Workplace -->
					<component
						:is="workplaceDialogComponent"
						:open="openWorkplaceDialog"
						@close="openWorkplaceDialog = false"
						@select="selectWorkplace"
					/>
				</form>
			</CardContent>
		</Card>
		<!-- Footer -->
		<div class="md:hidden flex justify-center items-center gap-2">
			<Button
				variant="outline"
				@click="handleCancel"
			>
				Отменить
			</Button>
			<Button
				type="submit"
				form="reassignForm"
				:disabled="isPending"
			>
				Сохранить
			</Button>
		</div>
	</div>
</template>
