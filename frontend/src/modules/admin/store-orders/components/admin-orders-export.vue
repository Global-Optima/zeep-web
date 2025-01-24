<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from '@/core/components/ui/dialog'
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage
} from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/core/components/ui/select'
import { useToast } from '@/core/components/ui/toast'
import type { OrdersExportFilterQuery } from '@/modules/orders/models/orders.models'
import { ordersService } from '@/modules/orders/services/orders.service'
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import * as z from 'zod'

// Validation schema for the form
const formSchema = toTypedSchema(
  z.object({
    startDate: z.string().optional().refine(
      (value) => !value || !isNaN(Date.parse(value)),
      { message: 'Некорректная дата начала' }
    ),
    endDate: z.string().optional().refine(
      (value) => !value || !isNaN(Date.parse(value)),
      { message: 'Некорректная дата окончания' }
    ),
    language: z.enum(['kk', 'ru', 'en'], { required_error: 'Выберите язык' }),
  }).refine(
      (data) =>
        !data.startDate ||
        !data.endDate ||
        new Date(data.startDate).getTime() <= new Date(data.endDate).getTime(),
      {
        message: 'Дата начала должна быть раньше даты окончания',
        path: ['startDate'],
      }
    )
)

const {toast} = useToast()

// Form setup
const { handleSubmit, isFieldDirty } = useForm({
  validationSchema: formSchema,
  initialValues: {
    language: "ru"
  }
})

// Submit handler
const onSubmit = handleSubmit(async (values) => {
  const filter: OrdersExportFilterQuery = {
    ...values,
    startDate: values.startDate ? new Date(values.startDate).toISOString() : undefined,
    endDate: values.endDate ? new Date(values.endDate).toISOString() : undefined,
  }

  await ordersService.exportOrders(filter).catch(() => {
    toast({
      title: 'Ошибка',
      description: "'Не удалось выполнить действие. Попробуйте снова.'",
      variant: 'destructive',
    })
  })
})
</script>

<template>
	<Dialog>
		<DialogTrigger as-child>
			<Button variant="outline">Экспорт</Button>
		</DialogTrigger>
		<DialogContent
			class="sm:max-w-[425px]"
			:include-close-button="false"
		>
			<DialogHeader>
				<DialogTitle>Экспорт заказов</DialogTitle>
			</DialogHeader>
			<form
				class="space-y-6"
				@submit="onSubmit"
			>
				<!-- Start Date -->
				<FormField
					v-slot="{ componentField }"
					name="startDate"
					:validate-on-blur="!isFieldDirty"
				>
					<FormItem>
						<FormLabel>Дата начала</FormLabel>
						<FormControl>
							<Input
								type="date"
								v-bind="componentField"
							/>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<!-- End Date -->
				<FormField
					v-slot="{ componentField }"
					name="endDate"
					:validate-on-blur="!isFieldDirty"
				>
					<FormItem>
						<FormLabel>Дата окончания</FormLabel>
						<FormControl>
							<Input
								type="date"
								v-bind="componentField"
							/>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<!-- Language -->
				<FormField
					v-slot="{ componentField }"
					name="language"
					:validate-on-blur="!isFieldDirty"
				>
					<FormItem>
						<FormLabel>Язык</FormLabel>
						<FormControl>
							<Select v-bind="componentField">
								<SelectTrigger>
									<SelectValue />
								</SelectTrigger>
								<SelectContent>
									<SelectItem value="ru">Русский</SelectItem>
									<SelectItem value="kk">Казахский</SelectItem>
									<SelectItem value="en">Английский</SelectItem>
								</SelectContent>
							</Select>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<DialogFooter>
					<Button type="submit">Экспортировать</Button>
				</DialogFooter>
			</form>
		</DialogContent>
	</Dialog>
</template>
