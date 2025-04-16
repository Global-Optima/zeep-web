<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { computed, ref } from 'vue'
import * as z from 'zod'

import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/core/components/ui/table'
import { ChevronDown, ChevronLeft } from 'lucide-vue-next'

import AdminSelectProvisionDialog from '@/modules/admin/provisions/components/admin-select-provision-dialog.vue'

import type { ProvisionDTO, ProvisionDetailsDTO } from '@/modules/admin/provisions/models/provision.models'
import type { CreateStoreProvisionDTO } from '@/modules/admin/store-provisions/models/store-provision.models'

import { provisionsService } from '@/modules/admin/provisions/services/provisions.service'
import { getRouteName } from "@/core/config/routes.config"

// Emits definition
const emits = defineEmits<{
  onSubmit: [dto: CreateStoreProvisionDTO]
  onCancel: []
}>()

// Validation Schema for the store provision details
const createStoreProvisionSchema = toTypedSchema(
  z.object({
    volume: z.number().min(0.0001, 'Введите объем заготовки'),
    expirationInMinutes: z.number().min(0, 'Введите срок годности в минутах'),
  })
)

// Form Setup
const { handleSubmit, resetForm, setFieldValue, values } = useForm({
  validationSchema: createStoreProvisionSchema,
})

// State for provision selection
const openProvisionDialog = ref(false)
const selectedProvision = ref<ProvisionDTO | null>(null)

// When a provision is selected from the dialog, update state and set default volume.
// Then the useQuery will fetch its full details.
function selectProvision(provision: ProvisionDTO) {
  selectedProvision.value = provision
  openProvisionDialog.value = false
  setFieldValue('volume', provision.absoluteVolume)
  setFieldValue('expirationInMinutes', provision.defaultExpirationInMinutes)
}

// Query to fetch the full provision details (including ingredients)
const provisionDetailsQuery = useQuery<ProvisionDetailsDTO>({
  queryKey: computed(() =>['admib-provision-details', selectedProvision.value?.id]),
  queryFn: async () => {
    if (!selectedProvision.value) {
      throw new Error('No provision selected')
    }
    return provisionsService.getProvisionById(selectedProvision.value.id)
  },
  enabled: computed(() => !!selectedProvision.value),
})

// Compute the technical map scaled to the store provision’s volume.
// Scaled quantity = templateIngredient.quantity * (storeVolume / templateAbsoluteVolume)
const scaledTechnicalMap = computed(() => {
  const details = provisionDetailsQuery.data
  if (!details || !values.volume || !details.value) {
    return []
  }

  const ratio = values.volume / details.value.absoluteVolume
  return details.value.ingredients.map(item => ({
    ingredientId: item.ingredient.id,
    name: item.ingredient.name,
    category: item.ingredient.category,
    quantity: item.quantity,
    scaledQuantity: item.quantity * ratio,
    unit: item.ingredient.unit.name,
  }))
})

// Handle form submission: build and emit the CreateStoreProvisionDTO
const onSubmit = handleSubmit((formValues) => {
  if (!selectedProvision.value) {
    // Optionally show an error message.
    return
  }
  const dto: CreateStoreProvisionDTO = {
    provisionId: selectedProvision.value.id,
    volume: formValues.volume,
    expirationInMinutes: formValues.expirationInMinutes,
  }
  emits('onSubmit', dto)
})

const onCancel = () => {
  resetForm()
  selectedProvision.value = null
  emits('onCancel')
}
</script>

<template>
	<div class="flex-1 gap-4 grid auto-rows-max mx-auto max-w-4xl">
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
				Создать заготовку
			</h1>
			<div class="hidden md:flex items-center gap-2 md:ml-auto">
				<Button
					variant="outline"
					type="button"
					@click="onCancel"
					>Отменить</Button
				>
				<Button
					type="submit"
					@click="onSubmit"
					>Сохранить</Button
				>
			</div>
		</div>

		<!-- Provision Template Selection -->
		<Card>
			<CardHeader>
				<CardTitle>Шаблон заготовки</CardTitle>
				<CardDescription>
					Выберите заготовку, которая будет использоваться как шаблон для расчетов.
				</CardDescription>
			</CardHeader>
			<CardContent>
				<div v-if="selectedProvision">
					<Button
						variant="outline"
						class="gap-3"
						@click="openProvisionDialog = true"
					>
						<span>
							{{ selectedProvision.name }} - {{ selectedProvision.absoluteVolume }}
							{{ selectedProvision.unit.name.toLowerCase()}}
						</span>

						<ChevronDown class="size-4 text-gray-600" />
					</Button>
				</div>
				<div v-else>
					<Button
						variant="outline"
						@click="openProvisionDialog = true"
					>
						Выбрать заготовку
					</Button>
				</div>
			</CardContent>
		</Card>

		<!-- Store Provision Details Form -->
		<Card v-if="selectedProvision">
			<CardHeader>
				<CardTitle>Детали заготовки</CardTitle>
				<CardDescription> Укажите требуемый объем и срок годности. </CardDescription>
			</CardHeader>
			<CardContent>
				<div class="gap-6 grid">
					<FormField
						name="volume"
						v-slot="{ componentField }"
					>
						<FormItem>
							<FormLabel>Объем ({{ selectedProvision.unit.name.toLowerCase() }})</FormLabel>
							<FormControl>
								<Input
									id="volume"
									type="number"
									v-bind="componentField"
									placeholder="Введите объем заготовки"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>
					<FormField
						name="expirationInMinutes"
						v-slot="{ componentField }"
					>
						<FormItem>
							<FormLabel>Срок годности (минут)</FormLabel>
							<FormControl>
								<Input
									id="expirationInMinutes"
									type="number"
									v-bind="componentField"
									placeholder="Введите срок годности в минутах"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>
				</div>
			</CardContent>
		</Card>

		<Card v-if="selectedProvision">
			<CardHeader>
				<CardTitle>Технологическая карта</CardTitle>
				<CardDescription> Технологическая карта на заданный обьем </CardDescription>
			</CardHeader>
			<CardContent>
				<Table v-if="selectedProvision">
					<TableHeader>
						<TableRow>
							<TableHead>Название</TableHead>
							<TableHead>Категория</TableHead>
							<TableHead>Абсолютный обьем</TableHead>
							<TableHead>Итоговый обьем</TableHead>
						</TableRow>
					</TableHeader>
					<TableBody>
						<TableRow
							v-for="ingredient in scaledTechnicalMap"
							:key="ingredient.ingredientId"
						>
							<TableCell>
								<RouterLink
									:to="{name: getRouteName('ADMIN_INGREDIENTS_DETAILS'), params: {id: ingredient.ingredientId}}"
									target="_blank"
									class="hover:text-primary underline transition-colors duration-300 underline-offset-4"
								>
									{{ ingredient.name }}
								</RouterLink>
							</TableCell>
							<TableCell>{{ ingredient.category.name }}</TableCell>
							<TableCell>{{ ingredient.quantity }} {{ ingredient.unit.toLowerCase() }}</TableCell>
							<TableCell class="font-semibold">
								{{ ingredient.scaledQuantity.toFixed(2) }}
								{{ ingredient.unit.toLowerCase() }}
							</TableCell>
						</TableRow>
					</TableBody>
				</Table>
				<div v-else>
					<p class="text-gray-500 text-sm">Выберите заготовку, чтобы увидеть техническую карту.</p>
				</div>
			</CardContent>
		</Card>

		<!-- Provision Selection Dialog -->
		<AdminSelectProvisionDialog
			:open="openProvisionDialog"
			@close="openProvisionDialog = false"
			@select="selectProvision"
		/>

		<!-- Mobile Footer -->
		<div class="md:hidden flex justify-center items-center gap-2">
			<Button
				variant="outline"
				@click="onCancel"
				>Отменить</Button
			>
			<Button
				type="submit"
				@click="onSubmit"
				>Сохранить</Button
			>
		</div>
	</div>
</template>
