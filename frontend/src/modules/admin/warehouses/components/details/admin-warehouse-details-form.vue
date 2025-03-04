<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { ref } from 'vue'
import * as z from 'zod'

// UI Components
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import AdminSelectRegionDialog from '@/modules/admin/regions/components/admin-select-region-dialog.vue'
import type { RegionDTO } from '@/modules/admin/regions/models/regions.model'
import type { UpdateWarehouseDTO, WarehouseDTO } from '@/modules/admin/warehouses/models/warehouse.model'
import { ChevronLeft } from 'lucide-vue-next'

// Props & Events
const {warehouse, readonly} = defineProps<{
  warehouse: WarehouseDTO
  readonly?: boolean
}>()

const emits = defineEmits<{
  (e: 'onSubmit', dto: UpdateWarehouseDTO): void
  (e: 'onCancel'): void
}>()

// Validation Schema
const updateWarehouseSchema = toTypedSchema(
  z.object({
    name: z.string().min(1, 'Введите название склада'),
    address: z.string().min(1, 'Введите адрес склада'),
    regionId: z.number().min(1, 'Введите регион склада')
  })
)

// Form Setup
const { handleSubmit, resetForm, setFieldValue } = useForm({
  validationSchema: updateWarehouseSchema,
  initialValues: {
    name: warehouse.name,
    address: warehouse.facilityAddress.address,
    regionId: warehouse.region.id
  }
})

// Handlers
const onSubmit = handleSubmit(async (formValues) => {
  if (readonly) return

  // Ensure the latest values are used
  emits('onSubmit', {
    name: formValues.name,
    facilityAddress: { address: formValues.address },
    regionId: formValues.regionId
  })
})

const onCancel = () => {
  resetForm({
    values: {
      name: warehouse.name,
      address: warehouse.facilityAddress.address,
      regionId: warehouse.region.id
    }
  })

  emits('onCancel')
}

// Region selection dialog
const openRegionDialog = ref(false)
const selectedRegion = ref<RegionDTO | null>(warehouse.region || null)

function selectRegion(region: RegionDTO) {
  selectedRegion.value = region
  openRegionDialog.value = false
  setFieldValue('regionId', region.id)
}
</script>

<template>
	<div class="flex-1 gap-4 grid auto-rows-max mx-auto max-w-6xl">
		<!-- Header -->
		<div class="flex items-center gap-4">
			<Button
				variant="outline"
				size="icon"
				@click="onCancel"
			>
				<ChevronLeft class="w-5 h-5" />
				<span class="sr-only">Назад</span>
			</Button>
			<h1 class="flex-1 sm:grow-0 font-semibold text-xl tracking-tight whitespace-nowrap shrink-0">
				Обновить склад - {{ warehouse.name }}
			</h1>

			<div
				class="hidden md:flex items-center gap-2 md:ml-auto"
				v-if="!readonly"
			>
				<Button
					variant="outline"
					type="button"
					@click="onCancel"
					:disabled="readonly"
					>Отменить</Button
				>
				<Button
					type="submit"
					@click="onSubmit"
					:disabled="readonly"
					>Сохранить</Button
				>
			</div>
		</div>

		<!-- Main Content -->
		<div class="gap-4 grid md:grid-cols-[1fr_250px] lg:grid-cols-3">
			<div class="items-start gap-4 grid lg:col-span-2 auto-rows-max">
				<Card>
					<CardHeader>
						<CardTitle>Редактирование склада</CardTitle>
						<CardDescription>Обновите детали склада.</CardDescription>
					</CardHeader>
					<CardContent>
						<form
							@submit.prevent="onSubmit"
							class="gap-6 grid"
						>
							<FormField
								name="name"
								v-slot="{ componentField }"
							>
								<FormItem>
									<FormLabel>Название склада</FormLabel>
									<FormControl>
										<Input
											id="name"
											v-bind="componentField"
											:readonly="readonly"
											placeholder="Введите название склада"
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>

							<FormField
								name="address"
								v-slot="{ componentField }"
							>
								<FormItem>
									<FormLabel>Название склада</FormLabel>
									<FormControl>
										<Input
											id="address"
											v-bind="componentField"
											:readonly="readonly"
											placeholder="Введите название склада"
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>
						</form>
					</CardContent>
				</Card>
			</div>

			<!-- Region Selection -->
			<div class="items-start gap-4 grid auto-rows-max">
				<Card>
					<CardHeader>
						<CardTitle>Регион</CardTitle>
						<CardDescription>Выберите регион склада</CardDescription>
					</CardHeader>
					<CardContent>
						<div>
							<Button
								variant="link"
								type="button"
								class="mt-0 p-0 h-fit text-primary underline"
								:disabled="readonly"
								@click="openRegionDialog = true"
							>
								{{ selectedRegion?.name || 'Регион не выбран' }}
							</Button>
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
				type="button"
				@click="onCancel"
				:disabled="readonly"
				>Отменить</Button
			>
			<Button
				type="submit"
				@click="onSubmit"
				:disabled="readonly"
				>Сохранить</Button
			>
		</div>

		<AdminSelectRegionDialog
			:open="openRegionDialog"
			@close="openRegionDialog = false"
			@select="selectRegion"
		/>
	</div>
</template>
