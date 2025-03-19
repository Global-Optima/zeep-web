<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import { useToast } from '@/core/components/ui/toast'
import { getKaspiConfig, KaspiService, saveKaspiConfig } from '@/core/integrations/kaspi.service'

import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { onMounted } from 'vue'
import * as z from 'zod'

const formSchema = toTypedSchema(
  z.object({
    host: z.string().min(1,
      'Неверный формат IP-адреса'
    ),
    name: z.string().min(2, 'Имя интеграции должно содержать минимум 2 символа'),
  })
);

const { toast } = useToast();

const { handleSubmit, isSubmitting, setValues } = useForm({
  validationSchema: formSchema,
});

onMounted(() => {
  const config = getKaspiConfig()
  if (config) {
    setValues(config);
  }
});

const onSubmit = handleSubmit(async (values) => {
  saveKaspiConfig(values)
  const kaspi = new KaspiService(values);

  try {
    await kaspi.registerTerminal();

    toast({
      title: 'Интеграция успешно завершена',
      description: `Терминал ${values.name} зарегистрирован.`,
      variant: "success"
    });
  } catch (error) {
    toast({
      title: 'Ошибка интеграции',
      description: error instanceof Error ? error.message : 'Неизвестная ошибка',
      variant: 'destructive',
    });
  }
});
</script>

<template>
	<form
		class="space-y-6"
		@submit="onSubmit"
	>
		<!-- IP-адрес POS-терминала -->
		<FormField
			v-slot="{ componentField }"
			name="host"
		>
			<FormItem>
				<FormLabel>Адрес POS-терминала</FormLabel>
				<FormControl>
					<Input
						type="text"
						placeholder="smartpos.kaspipos.kz"
						v-bind="componentField"
					/>
				</FormControl>
				<FormDescription>
					Укажите адрес POS-терминала внутри сети для подключения.
				</FormDescription>
				<FormMessage />
			</FormItem>
		</FormField>

		<!-- Имя интеграции -->
		<FormField
			v-slot="{ componentField }"
			name="name"
		>
			<FormItem>
				<FormLabel>Имя интеграции</FormLabel>
				<FormControl>
					<Input
						type="text"
						placeholder="Введите имя интеграции"
						v-bind="componentField"
					/>
				</FormControl>
				<FormDescription>
					Введите название, которое будет использоваться для этой интеграции.
				</FormDescription>
				<FormMessage />
			</FormItem>
		</FormField>

		<!-- Кнопка отправки -->
		<Button
			type="submit"
			:disabled="isSubmitting"
		>
			<span v-if="!isSubmitting">Сохранить настройки и интегрировать</span>
			<span v-else>Подключение...</span>
		</Button>
	</form>
</template>

<style scoped></style>
