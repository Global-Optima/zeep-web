<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger
} from '@/core/components/ui/dialog'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage
} from '@/core/components/ui/form'
import { Stepper, StepperItem, StepperTitle, StepperTrigger } from '@/core/components/ui/stepper'

import { toTypedSchema } from '@vee-validate/zod'
import { ref } from 'vue'
import * as z from 'zod'

const formSchema = [
  // Step 1
  z.object({
    userName: z.string().optional(),
  }),
  // Step 2
  z.object({
    paymentMethod: z.union([z.literal('bankCard'), z.literal('cash')]),
  }),
  // Removed Step 3
]

const stepIndex = ref(1)
const steps = [
  {
    step: 1,
    title: 'Введите имя пользователя',
    description: 'Или сгенерируйте случайную метку',
  },
  {
    step: 2,
    title: 'Выберите способ оплаты',
    description: 'Выберите карту или наличные',
  },
  // Removed Step 3
]

function onSubmit(values: unknown) {
  console.log('Order created:', values)
  // Optionally, close the dialog or show a success message here
}
</script>

<template>
	<Form
		v-slot="{ meta, values, validate, setFieldValue }"
		as=""
		keep-values
		:validation-schema="toTypedSchema(formSchema[stepIndex - 1])"
		@submit="onSubmit"
	>
		<Dialog>
			<DialogTrigger as-child>
				<button
					class="px-6 py-4 rounded-2xl bg-primary text-primary-foreground text-base sm:text-2xl flex items-center gap-2"
				>
					Подтвердить
					<Icon
						icon="mingcute:right-line"
						class="text-base sm:text-2xl"
					/>
				</button>
			</DialogTrigger>
			<DialogContent class="w-full max-w-[90vw] p-6">
				<DialogHeader>
					<DialogTitle class="text-2xl">Оформить заказ</DialogTitle>
					<DialogDescription class="text-lg">
						Пожалуйста, заполните необходимые данные для оформления заказа.
					</DialogDescription>
				</DialogHeader>

				<Stepper
					v-slot="{ isNextDisabled, isPrevDisabled, nextStep, prevStep }"
					v-model="stepIndex"
					class="block w-full"
				>
					<form
						@submit.prevent="() => {
              validate()
              if (stepIndex === steps.length && meta.valid) {
                onSubmit(values)
              }
            }"
					>
						<div class="hidden">
							<StepperItem
								v-for="step in steps"
								:key="step.step"
								:step="step.step"
							>
								<StepperTrigger as-child> </StepperTrigger>

								<div class="mt-5 flex flex-col items-center text-center">
									<StepperTitle></StepperTitle>
								</div>
							</StepperItem>
						</div>

						<div class="flex flex-col gap-4 mt-4">
							<!-- Step 1: User Name -->
							<template v-if="stepIndex === 1">
								<FormField
									v-slot="{ componentField }"
									name="userName"
								>
									<FormItem>
										<FormLabel class="text-xl font-normal">
											Введите ваше имя для отображения заказа
										</FormLabel>
										<FormControl>
											<input
												type="text"
												class="w-full border border-border rounded-2xl px-4 py-4 text-xl"
												placeholder="Алихан"
												v-bind="componentField"
											/>
										</FormControl>
										<FormMessage />
									</FormItem>
								</FormField>
								<Button
									variant="ghost"
									class="mt-2 text-xl"
									@click="setFieldValue('userName', 'User' + Math.floor(Math.random() * 1000))"
								>
									Сгенерировать случайную метку
								</Button>
							</template>

							<!-- Step 2: Payment Method -->
							<template v-if="stepIndex === 2">
								<FormField name="paymentMethod">
									<FormItem>
										<FormLabel class="text-xl font-normal">Способ оплаты</FormLabel>
										<FormControl>
											<div class="flex gap-4">
												<div
													:class="[
                            'p-4 border rounded cursor-pointer',
                            values.paymentMethod === 'bankCard' ? 'border-primary' : 'border-gray-300'
                          ]"
													@click="setFieldValue('paymentMethod', 'bankCard')"
												>
													<p class="text-xl">Банковская карта</p>
												</div>
												<div
													:class="[
                            'p-4 border rounded cursor-pointer',
                            values.paymentMethod === 'cash' ? 'border-primary' : 'border-gray-300'
                          ]"
													@click="setFieldValue('paymentMethod', 'cash')"
												>
													<p class="text-xl">Наличные</p>
												</div>
											</div>
										</FormControl>
										<FormMessage />
									</FormItem>
								</FormField>
							</template>
						</div>

						<!-- Navigation Buttons -->
						<div class="flex items-center justify-between mt-4">
							<button
								:disabled="isPrevDisabled"
								variant="outline"
								size="lg"
								@click="prevStep()"
								class="px-6 py-4 rounded-2xl border-border border-2 text-base sm:text-xl"
							>
								Назад
							</button>

							<div class="flex items-center gap-3">
								<button
									v-if="stepIndex < steps.length"
									type="button"
									:disabled="isNextDisabled"
									@click="meta.valid && nextStep()"
									class="px-6 py-4 rounded-2xl bg-primary text-primary-foreground text-base sm:text-xl"
								>
									Далее
								</button>

								<button
									v-else
									type="submit"
									:disabled="!meta.valid"
									class="px-6 py-4 rounded-2xl bg-primary text-primary-foreground text-base sm:text-xl"
								>
									Отправить
								</button>
							</div>
						</div>
					</form>
				</Stepper>
			</DialogContent>
		</Dialog>
	</Form>
</template>
