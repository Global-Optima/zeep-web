<template>
	<div class="flex-1 gap-4 grid auto-rows-max mx-auto max-w-6xl">
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
				{{  store.name }}
			</h1>

			<div
				v-if="!readonly"
				class="hidden md:flex items-center gap-2 md:ml-auto"
			>
				<Button
					variant="outline"
					type="button"
					@click="handleCancel"
					>Отменить</Button
				>
				<Button
					type="submit"
					@click="submitForm"
					>Сохранить</Button
				>
			</div>
		</div>

		<!-- Main Content -->
		<div class="gap-4 grid md:grid-cols-[1fr_250px] lg:grid-cols-3">
			<div class="items-start gap-4 grid lg:col-span-2 auto-rows-max">
				<Card>
					<CardHeader>
						<CardTitle>Детали кафе</CardTitle>
						<CardDescription>Заполните форму ниже, чтобы обновить кафе.</CardDescription>
					</CardHeader>
					<CardContent>
						<form
							@submit="submitForm"
							class="gap-6 grid"
						>
							<FormField
								name="name"
								v-slot="{ componentField }"
							>
								<FormItem>
									<FormLabel>Название кафе</FormLabel>
									<FormControl>
										<Input
											v-bind="componentField"
											placeholder="Введите название кафе"
											:readonly="readonly"
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>

							<!-- Facility Address -->
							<FormField
								name="facilityAddress.address"
								v-slot="{ componentField }"
							>
								<FormItem>
									<FormLabel>Адрес кафе</FormLabel>
									<FormControl>
										<Input
											v-bind="componentField"
											placeholder="Введите адрес"
											:readonly="readonly"
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>

							<!-- Contact Phone and Email -->
							<div class="flex gap-4">
								<div class="w-1/2">
									<FormField
										name="contactPhone"
										v-slot="{ componentField }"
									>
										<FormItem>
											<FormLabel>Контактный телефон</FormLabel>
											<FormControl>
												<Input
													v-bind="componentField"
													placeholder="+7XXXXXXXXXX"
													:readonly="readonly"
												/>
											</FormControl>
											<FormMessage />
										</FormItem>
									</FormField>
								</div>
								<div class="w-1/2">
									<FormField
										name="contactEmail"
										v-slot="{ componentField }"
									>
										<FormItem>
											<FormLabel>Контактный Email</FormLabel>
											<FormControl>
												<Input
													type="email"
													v-bind="componentField"
													placeholder="example@example.com"
													:readonly="readonly"
												/>
											</FormControl>
											<FormMessage />
										</FormItem>
									</FormField>
								</div>
							</div>

							<!-- Store Hours -->
							<FormField
								name="storeHours"
								v-slot="{ componentField }"
							>
								<FormItem>
									<FormLabel>Часы работы</FormLabel>
									<FormControl>
										<Input
											v-bind="componentField"
											placeholder="Введите часы работы (например, 9:00-18:00)"
											:readonly="readonly"
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>
						</form>
					</CardContent>
				</Card>
			</div>

			<div class="items-start gap-4 grid auto-rows-max">
				<Card>
					<CardHeader>
						<CardTitle>Франчайзи (опционально)</CardTitle>
						<CardDescription>Выберите франчайзи</CardDescription>
					</CardHeader>
					<CardContent>
						<div class="flex items-center gap-4">
							<Button
								variant="link"
								class="mt-0 p-0 h-fit text-primary underline"
								@click="openFranchiseeDialog = true"
							>
								{{ selectedFranchisee?.name || 'Франчайзи не выбран' }}
							</Button>

							<button
								v-if="selectedFranchisee"
								@click="selectFranchisee(null)"
							>
								<X class="size-4 text-gray-600" />
							</button>
						</div>
					</CardContent>
				</Card>

				<Card>
					<CardHeader>
						<CardTitle>Склад</CardTitle>
						<CardDescription>Выберите склад</CardDescription>
					</CardHeader>
					<CardContent>
						<div v-if="!readonly">
							<Button
								variant="link"
								class="mt-0 p-0 h-fit text-primary underline"
								@click="openWarehouseDialog = true"
							>
								{{ selectedWarehouse?.name || 'Склад не выбран' }}
							</Button>
						</div>
						<div v-else>
							<span>{{ selectedWarehouse?.name || 'Склад не выбран' }}</span>
						</div>
					</CardContent>
				</Card>
			</div>
		</div>

		<!-- Footer -->
		<div
			v-if="!readonly"
			class="md:hidden flex justify-center items-center gap-2"
		>
			<Button
				variant="outline"
				@click="handleCancel"
				>Отменить</Button
			>
			<Button
				type="submit"
				@click="submitForm"
				>Сохранить</Button
			>
		</div>

		<AdminSelectFranchiseeDialog
			v-if="!readonly"
			:open="openFranchiseeDialog"
			@close="openFranchiseeDialog = false"
			@select="selectFranchisee"
		/>
		<AdminSelectWarehouseDialog
			v-if="!readonly"
			:open="openWarehouseDialog"
			@close="openWarehouseDialog = false"
			@select="selectWarehouse"
		/>
	</div>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/core/components/ui/card'
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage
} from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import { phoneValidationSchema } from '@/core/validators/phone.validator'
import type { FranchiseeDTO } from '@/modules/admin/franchisees/models/franchisee.model'
import type { UpdateStoreDTO } from '@/modules/admin/stores/models/stores-dto.model'
import type { StoreDTO } from '@/modules/admin/stores/models/stores.models'
import type { WarehouseDTO } from '@/modules/admin/warehouses/models/warehouse.model'
import { toTypedSchema } from '@vee-validate/zod'
import { ChevronLeft, X } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { defineAsyncComponent, ref } from 'vue'
import * as z from 'zod'

const AdminSelectFranchiseeDialog = defineAsyncComponent(() => import('@/modules/admin/franchisees/components/admin-select-franchisee-dialog.vue'))
const AdminSelectWarehouseDialog = defineAsyncComponent(() => import('@/modules/admin/warehouses/components/admin-select-warehouse-dialog.vue'))

const {readonly, store} = defineProps<{ readonly?: boolean, store: StoreDTO }>()

const emit = defineEmits<{
	(e: 'onSubmit', formValues: UpdateStoreDTO): void
	(e: 'onCancel'): void
}>()

const schema = toTypedSchema(
	z.object({
		name: z.string().min(2, 'Название должно содержать минимум 2 символа'),
		facilityAddress: z.object({
			address: z.string().min(5, 'Адрес должен содержать минимум 5 символов'),
		}),
		contactPhone: phoneValidationSchema,
		contactEmail: z.string().email('Введите действительный адрес электронной почты'),
		storeHours: z.string().min(5, 'Часы работы должны быть указаны'),
    warehouseId: z.number().min(1, 'Введите склад'),
    franchiseId: z.number().optional()
	}),
)

const { handleSubmit, resetForm, setFieldValue } = useForm({
	validationSchema: schema,
  initialValues: {
    name: store.name,
    warehouseId: store.warehouse.id,
    facilityAddress: {
      address: store.facilityAddress.address
    },
    contactPhone: store.contactPhone,
    contactEmail: store.contactEmail,
    storeHours: store.storeHours,
    franchiseId: store.franchisee?.id,
  }
})

const submitForm = handleSubmit((formValues) => {
  const dto: UpdateStoreDTO = {
    name: formValues.name,
    warehouseId: formValues.warehouseId,
    facilityAddress: {
      address: formValues.facilityAddress.address
    },
    isActive: store.isActive,
    contactPhone: formValues.contactPhone,
    contactEmail: formValues.contactEmail,
    storeHours: formValues.storeHours,
    franchiseId:  formValues.franchiseId ?? null
  }
  emit('onSubmit', dto)
})

const handleCancel = () => { resetForm(); emit('onCancel') }

const openWarehouseDialog = ref(false)
const openFranchiseeDialog = ref(false)
const selectedWarehouse = ref<WarehouseDTO | null>(store.warehouse)
const selectedFranchisee = ref<FranchiseeDTO  | null>(store.franchisee || null)

function selectWarehouse(warehouse: WarehouseDTO) {
  selectedWarehouse.value = warehouse;
  openWarehouseDialog.value = false
  setFieldValue('warehouseId', warehouse.id)
}

function selectFranchisee(franchisee: FranchiseeDTO | null) {
  if (!franchisee) {
    selectedFranchisee.value = null
    setFieldValue('franchiseId', undefined)
    return
  }

  selectedFranchisee.value = franchisee
  openFranchiseeDialog.value = false
  setFieldValue('franchiseId', franchisee.id)
}
</script>
