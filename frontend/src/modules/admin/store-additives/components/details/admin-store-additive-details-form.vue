<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import * as z from 'zod'

import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel } from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import type { StoreAdditiveDTO, UpdateStoreAdditiveDTO } from '@/modules/admin/store-additives/models/store-additves.model'
import { ChevronLeft } from 'lucide-vue-next'

const emits = defineEmits<{
  onSubmit: [updatedAdditive: UpdateStoreAdditiveDTO]
  onCancel: []
}>()

const {initialAdditive} = defineProps<{
  initialAdditive: StoreAdditiveDTO
}>()

const updateStoreAdditiveSchema = toTypedSchema(
  z.object({
    storePrice: z
      .number()
      .min(0, 'Цена не может быть меньше 0')
      .describe('Введите новую цену в магазине'),
  })
)

const { handleSubmit, resetForm } = useForm({
  validationSchema: updateStoreAdditiveSchema,
  initialValues: {
    storePrice: initialAdditive.storePrice,
  },
})

const onSubmit = handleSubmit((formValues) => {
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
			<h1 class="font-semibold text-xl tracking-tight shrink-0">{{ initialAdditive.name }}</h1>
			<div class="md:flex items-center gap-2 hidden md:ml-auto">
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

		<Card>
			<CardHeader>
				<CardTitle>Детали добавки</CardTitle>
				<CardDescription>Измените цену добавки для магазина.</CardDescription>
			</CardHeader>
			<CardContent class="space-y-6">
				<FormField name="basePrice">
					<FormItem>
						<FormLabel>Базовая цена</FormLabel>
						<FormControl>
							<Input
								id="storePrice"
								type="number"
								:model-value="initialAdditive.basePrice"
								disabled
							/>
						</FormControl>
					</FormItem>
				</FormField>

				<FormField
					name="storePrice"
					v-slot="{ componentField }"
				>
					<FormItem>
						<FormLabel>Цена в магазине</FormLabel>

						<FormControl>
							<Input
								id="storePrice"
								type="number"
								v-bind="componentField"
								placeholder="Введите цену в магазине"
							/>
						</FormControl>
					</FormItem>
				</FormField>
			</CardContent>
		</Card>

		<div class="flex justify-center items-center gap-2 md:hidden">
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
