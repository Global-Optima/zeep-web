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
import { KaspiService, type KaspiConfig } from '@/core/integrations/kaspi.service'

import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { onMounted } from 'vue'
import * as z from 'zod'

const formSchema = toTypedSchema(
  z.object({
    posIpAddress: z.string().regex(
      /^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$/,
      'Неверный формат IP-адреса'
    ),
    integrationName: z.string().min(2, 'Имя интеграции должно содержать минимум 2 символа'),
  })
);

const { toast } = useToast();

const { handleSubmit, isSubmitting, setValues } = useForm({
  validationSchema: formSchema,
  initialValues: {
    posIpAddress: '',
    integrationName: '',
  },
});

onMounted(() => {
  const savedConfig = localStorage.getItem('ZEEP_KASPI_CONFIG');
  if (savedConfig) {
    try {
      const parsedConfig: KaspiConfig = JSON.parse(savedConfig);
      setValues(parsedConfig);
    } catch (error) {
      console.warn('Ошибка загрузки конфигурации:', error);
    }
  }
});

const onSubmit = handleSubmit(async (values) => {
  localStorage.setItem('ZEEP_KASPI_CONFIG', JSON.stringify(values));
  const kaspi = new KaspiService(values);

  try {
    await kaspi.registerTerminal(values.integrationName);

    toast({
      title: '✅ Интеграция успешно завершена',
      description: `Терминал ${values.integrationName} зарегистрирован.`,
    });
  } catch (error) {
    toast({
      title: '❌ Ошибка интеграции',
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
			name="posIpAddress"
		>
			<FormItem>
				<FormLabel>IP-адрес POS-терминала</FormLabel>
				<FormControl>
					<Input
						type="text"
						placeholder="192.168.x.x"
						v-bind="componentField"
					/>
				</FormControl>
				<FormDescription> Укажите IP-адрес POS-терминала для подключения. </FormDescription>
				<FormMessage />
			</FormItem>
		</FormField>

		<!-- Имя интеграции -->
		<FormField
			v-slot="{ componentField }"
			name="integrationName"
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
