<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import * as z from 'zod'

// UI Components
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel } from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import type { StoreAdditiveDTO, UpdateStoreAdditiveDTO } from '@/modules/admin/store-additives/models/store-additves.model'
import { ChevronLeft, ExternalLink } from 'lucide-vue-next'
import { getRouteName } from "@/core/config/routes.config"


const emits = defineEmits<{
  onSubmit: [updatedAdditive: UpdateStoreAdditiveDTO]
  onCancel: []
}>()

const { initialAdditive, readonly = false } = defineProps<{
  initialAdditive: StoreAdditiveDTO
  readonly?: boolean
}>()

const updateStoreAdditiveSchema = toTypedSchema(
  z.object({
    storePrice: z
      .number()
      .min(0, 'Цена не может быть меньше 0')
      .describe('Введите новую цену в кафе'),
  })
)

const { handleSubmit, resetForm } = useForm({
  validationSchema: updateStoreAdditiveSchema,
  initialValues: {
    storePrice: initialAdditive.storePrice,
  },
})

const onSubmit = handleSubmit((formValues) => {
  if (readonly) return
  const dto: UpdateStoreAdditiveDTO = {
    storePrice: formValues.storePrice,
  }
  emits('onSubmit', dto)
})

const onCancel = () => {
  resetForm()
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
				type="button"
				@click="onCancel"
			>
				<ChevronLeft class="w-5 h-5" />
				<span class="sr-only">Назад</span>
			</Button>

			<div class="flex items-center gap-3">
				<h1 class="font-semibold text-xl tracking-tight shrink-0">
					{{ initialAdditive.name }}
				</h1>
				<RouterLink
					:to="{name: getRouteName('ADMIN_ADDITIVE_DETAILS'), params: {id: initialAdditive.additiveId}}"
					target="_blank"
					class="hover:text-primary transition-colors duration-300"
				>
					<ExternalLink class="size-5 -mt-1" />
				</RouterLink>
			</div>

			<div
				v-if="!readonly"
				class="md:flex items-center gap-2 hidden md:ml-auto"
			>
				<Button
					variant="outline"
					type="button"
					@click="onCancel"
				>
					Отменить
				</Button>
				<Button
					type="submit"
					@click="onSubmit"
				>
					Сохранить
				</Button>
			</div>
		</div>

		<!-- Main Content -->
		<Card>
			<CardHeader>
				<CardTitle>Детали модификатора</CardTitle>
				<CardDescription v-if="!readonly">Измените цену модификатора для кафе.</CardDescription>
			</CardHeader>
			<CardContent class="space-y-6">
				<!-- Base Price -->
				<FormField name="basePrice">
					<FormItem>
						<FormLabel>Базовая цена</FormLabel>
						<FormControl>
							<Input
								id="storePrice"
								type="number"
								:model-value="initialAdditive.basePrice"
								readonly
							/>
						</FormControl>
					</FormItem>
				</FormField>

				<!-- Store Price -->
				<FormField
					name="storePrice"
					v-slot="{ componentField }"
				>
					<FormItem>
						<FormLabel>Цена в кафе</FormLabel>
						<FormControl>
							<Input
								id="storePrice"
								type="number"
								v-bind="componentField"
								placeholder="Введите цену в кафе"
								:readonly="readonly"
							/>
						</FormControl>
					</FormItem>
				</FormField>
			</CardContent>
		</Card>

		<!-- Mobile Footer -->
		<div
			v-if="!readonly"
			class="flex justify-center items-center gap-2 md:hidden"
		>
			<Button
				variant="outline"
				type="button"
				@click="onCancel"
			>
				Отменить
			</Button>
			<Button
				type="submit"
				@click="onSubmit"
			>
				Сохранить
			</Button>
		</div>
	</div>
</template>
