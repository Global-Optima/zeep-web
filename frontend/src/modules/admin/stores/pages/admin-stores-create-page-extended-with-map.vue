<template>
	<div class="flex md:flex-row flex-col gap-6 mx-auto">
		<!-- Left Side: Store Details Form -->
		<div class="w-full md:w-2/3">
			<Card>
				<CardHeader>
					<CardTitle>{{ isEditing ? 'Обновить магазин' : 'Создать магазин' }}</CardTitle>
					<CardDescription>
						Заполните форму ниже, чтобы {{ isEditing ? 'обновить' : 'создать' }} магазин.
					</CardDescription>
				</CardHeader>
				<CardContent>
					<form
						@submit.prevent="onSubmit"
						class="gap-6 grid"
					>
						<!-- Store Name -->
						<FormField
							name="name"
							v-slot="{ field, errorMessage }"
						>
							<FormItem>
								<FormLabel>Название магазина</FormLabel>
								<FormControl>
									<Input
										v-model="field.value"
										placeholder="Введите название магазина"
									/>
								</FormControl>
								<FormMessage v-if="errorMessage">{{ errorMessage }}</FormMessage>
							</FormItem>
						</FormField>
						<!-- Address -->
						<FormField
							name="address"
							v-slot="{ field, errorMessage }"
						>
							<FormItem>
								<FormLabel>Адрес</FormLabel>
								<FormControl>
									<Input
										v-model="field.value"
										placeholder="Введите адрес"
									/>
								</FormControl>
								<FormMessage v-if="errorMessage">{{ errorMessage }}</FormMessage>
							</FormItem>
						</FormField>
						<!-- Contact Information -->
						<div class="gap-4 grid grid-cols-1 sm:grid-cols-2">
							<!-- Phone -->
							<FormField
								name="phone"
								v-slot="{ field, errorMessage }"
							>
								<FormItem>
									<FormLabel>Телефон</FormLabel>
									<FormControl>
										<Input
											v-model="field.value"
											placeholder="+7 (___) ___-__-__"
										/>
									</FormControl>
									<FormMessage v-if="errorMessage">{{ errorMessage }}</FormMessage>
								</FormItem>
							</FormField>
							<!-- Email -->
							<FormField
								name="email"
								v-slot="{ field, errorMessage }"
							>
								<FormItem>
									<FormLabel>Электронная почта</FormLabel>
									<FormControl>
										<Input
											type="email"
											v-model="field.value"
											placeholder="example@example.com"
										/>
									</FormControl>
									<FormMessage v-if="errorMessage">{{ errorMessage }}</FormMessage>
								</FormItem>
							</FormField>
						</div>
						<!-- Description -->
						<FormField
							name="description"
							v-slot="{ field, errorMessage }"
						>
							<FormItem>
								<FormLabel>Описание</FormLabel>
								<FormControl>
									<Textarea
										v-model="field.value"
										placeholder="Краткое описание магазина"
									/>
								</FormControl>
								<FormMessage v-if="errorMessage">{{ errorMessage }}</FormMessage>
							</FormItem>
						</FormField>
						<!-- Action Buttons -->
						<div class="flex gap-4 mt-6">
							<Button
								type="submit"
								class="flex-1"
							>
								{{ isEditing ? 'Обновить' : 'Создать' }}
							</Button>
							<Button
								variant="outline"
								class="flex-1"
								@click="handleCancel"
							>
								Отмена
							</Button>
						</div>
					</form>
				</CardContent>
			</Card>
		</div>

		<!-- Right Side: Map and Working Hours -->
		<div class="flex flex-col gap-6 w-full md:w-1/3">
			<!-- Map Integration -->
			<Card>
				<CardHeader>
					<CardTitle>Расположение магазина</CardTitle>
					<CardDescription>Выберите местоположение магазина на карте</CardDescription>
				</CardHeader>
				<CardContent>
					<AdminStoresCreateMap v-model="store.location" />
				</CardContent>
			</Card>
			<!-- Working Hours -->
			<Card>
				<CardHeader>
					<CardTitle>Часы работы</CardTitle>
					<CardDescription>Установите часы работы для каждого дня</CardDescription>
				</CardHeader>
				<CardContent>
					<AdminStoresCreateWorkHours v-model="store.workingHours" />
				</CardContent>
			</Card>
		</div>
	</div>
</template>

<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import * as z from 'zod'


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
  FormMessage,
} from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import { Textarea } from '@/core/components/ui/textarea'
import AdminStoresCreateMap from '@/modules/admin/stores/components/create/admin-stores-create-map.vue'
import AdminStoresCreateWorkHours from '@/modules/admin/stores/components/create/admin-stores-create-work-hours.vue'

const router = useRouter();
const isEditing = ref(false); // Set to true if editing an existing store

// Define the Zod schema for form validation
const schema = toTypedSchema(
  z.object({
    name: z.string().min(2, 'Название должно содержать минимум 2 символа').max(100, 'Название должно содержать не более 100 символов'),
    address: z.string().min(5, 'Адрес должен содержать минимум 5 символов'),
    phone: z.string().min(7, 'Телефон должен содержать минимум 7 символов').max(15, 'Телефон должен содержать не более 15 символов'),
    email: z.string().email('Введите действительный адрес электронной почты'),
    description: z.string().optional(),
  })
);

// Initialize the form with Vee-Validate
const { handleSubmit } = useForm({
  validationSchema: schema,
  initialValues: {
    name: '',
    address: '',
    phone: '',
    email: '',
    description: '',
  },
});

// Store data (including location and working hours)
const store = ref({
  location: {lat: 0, lng:0, address: ""} as {address: string, lat: number; lng: number },
  workingHours: {},
});

// Handle form submission
const onSubmit = handleSubmit((formValues) => {
  const storeData = {
    ...formValues,
    location: store.value.location,
    workingHours: store.value.workingHours,
  };
  if (isEditing.value) {
  } else {
  }
  router.push('/admin/stores'); // Redirect after saving
});

// Handle cancel action
const handleCancel = () => {
  router.push('/admin/stores');
};
</script>
