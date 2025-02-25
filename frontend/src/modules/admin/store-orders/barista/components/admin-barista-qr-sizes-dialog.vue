<template>
	<Dialog
		:open="isModalOpen"
		@update:open="handleModalToggle"
	>
		<!-- КНОПКА-ТРЕЙГЕР -->
		<DialogTrigger as-child>
			<Button
				size="icon"
				variant="outline"
			>
				<Settings class="size-5" />
			</Button>
		</DialogTrigger>
		<!-- СОДЕРЖИМОЕ ДИАЛОГА -->
		<DialogContent
			class="sm:max-w-md"
			:include-close-button="false"
		>
			<DialogHeader>
				<DialogTitle>Размер для принтера QR-кодов</DialogTitle>
				<DialogDescription>
					Выберите стандартный размер или задайте собственные размеры.
				</DialogDescription>
			</DialogHeader>
			<!-- ФОРМА -->
			<form
				id="dialogForm"
				@submit="onSubmit"
				class="space-y-6"
			>
				<!-- SELECT ПРЕДОПРЕДЕЛЕННЫХ РАЗМЕРОВ ИЛИ "СВОЙ РАЗМЕР" -->
				<FormField
					name="standard"
					v-slot="{ componentField }"
				>
					<FormItem>
						<FormLabel>Стандарт принтера</FormLabel>
						<!-- shadcn-vue Select -->
						<Select
							v-bind="componentField"
							@update:modelValue="onStandardChange"
						>
							<FormControl>
								<!-- Текущее выбранное значение -->
								<SelectTrigger>
									<SelectValue placeholder="Выберите размер" />
								</SelectTrigger>
							</FormControl>
							<SelectContent>
								<SelectGroup>
									<SelectItem
										v-for="option in PREDEFINED_QR_SIZES"
										:key="option.label"
										:value="option.label"
									>
										{{ option.label }}
									</SelectItem>
									<!-- Свой размер -->
									<SelectItem :value="QRPrinterSizeName.CUSTOM">
										{{ QRPrinterSizeName.CUSTOM }}
									</SelectItem>
								</SelectGroup>
							</SelectContent>
						</Select>
						<FormMessage />
					</FormItem>
				</FormField>
				<!-- ПОЛЯ ДЛЯ "СВОЙ РАЗМЕР" -->
				<div
					v-if="isCustomSelected"
					class="gap-4 grid grid-cols-2"
				>
					<FormField
						name="width"
						v-slot="{ componentField }"
					>
						<FormItem>
							<FormLabel>Ширина (мм)</FormLabel>
							<FormControl>
								<Input
									type="number"
									step="0.1"
									min="10"
									v-bind="componentField"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>
					<FormField
						name="height"
						v-slot="{ componentField }"
					>
						<FormItem>
							<FormLabel>Высота (мм)</FormLabel>
							<FormControl>
								<Input
									type="number"
									step="0.1"
									min="10"
									v-bind="componentField"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>
				</div>
				<!-- КНОПКИ ДИАЛОГА -->
			</form>
			<DialogFooter class="sm:justify-start">
				<DialogClose as-child>
					<Button
						type="button"
						variant="secondary"
					>
						Отмена
					</Button>
				</DialogClose>
				<Button
					type="submit"
					variant="default"
					form="dialogForm"
				>
					Сохранить
				</Button>
			</DialogFooter>
		</DialogContent>
	</Dialog>
</template>

<script setup lang="ts">
import { useForm } from "vee-validate"
import { computed, ref } from "vue"
import * as z from "zod"
// shadcn-vue UI
import { Button } from "@/core/components/ui/button"
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/core/components/ui/dialog"
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/core/components/ui/form"
import { Input } from "@/core/components/ui/input"
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/core/components/ui/select"
import { useToast } from '@/core/components/ui/toast'
import { PREDEFINED_QR_SIZES, QRPrinterSizeName, SAVED_BARISTA_QR_SETTINGS_KEY } from '@/core/constants/qr.constants'
import { toTypedSchema } from '@vee-validate/zod'
import { Settings } from 'lucide-vue-next'

const { toast } = useToast();

/**
 * Схема валидации (Zod)
 */
const formSchema = z.object({
  standard: z.nativeEnum(QRPrinterSizeName),
  width: z.preprocess((val) => parseFloat(val as string), z.number().min(10)),
  height: z.preprocess((val) => parseFloat(val as string), z.number().min(10)),
});

/**
 * Преобразуем Zod → vee-validate
 */
const typedSchema = toTypedSchema(formSchema);

/**
 * ИНИЦИАЛЬНЫЕ ЗНАЧЕНИЯ
 */
const INITIAL_VALUES = {
  standard: QRPrinterSizeName.SIZE_80x80,
  width: 80,
  height: 80,
};

/**
 * Создаём форму
 */
const { handleSubmit, values, setValues, resetForm } = useForm({
  validationSchema: typedSchema,
  initialValues: INITIAL_VALUES,
});

/**
 * ОПРЕДЕЛЯЕМ, КОГДА ПОКАЗЫВАТЬ ПОЛЯ ДЛЯ "СВОЙ РАЗМЕР"
 */
const isCustomSelected = computed(() => values.standard === QRPrinterSizeName.CUSTOM);

/**
 * МОДАЛЬНОЕ ОКНО: Управление состоянием открытия
 */
const isModalOpen = ref(false);

/**
 * Загрузка данных из localStorage при открытии модального окна
 */
function loadFromLocalStorage() {
  const stored = localStorage.getItem(SAVED_BARISTA_QR_SETTINGS_KEY);
  if (stored) {
    try {
      const parsed = JSON.parse(stored);
      if (parsed && typeof parsed === "object") {
        setValues(parsed);
      }
    } catch (err) {
      console.error("Ошибка загрузки из LocalStorage:", err);
    }
  } else {
    resetForm({ values: INITIAL_VALUES });
  }
}

/**
 * Обработчик открытия/закрытия модального окна
 */
function handleModalToggle(isOpen: boolean) {
  isModalOpen.value = isOpen;

  if (isOpen) {
    loadFromLocalStorage();
  }
}

/**
 * ОБРАБОТЧИК ВЫБОРА СТАНДАРТА:
 */
function onStandardChange(selected: string) {
  if (selected !== QRPrinterSizeName.CUSTOM) {
    const found = PREDEFINED_QR_SIZES.find((size) => size.label === selected);
    if (found) {
      setValues({
        standard: found.label,
        width: found.width,
        height: found.height,
      });
    }
  } else {
    setValues({
      standard: QRPrinterSizeName.CUSTOM,
      width: 60,
      height: 60,
    });
  }
}

/**
 * СОХРАНЕНИЕ (SUBMIT) В localStorage
 */
const onSubmit = handleSubmit((values) => {
  localStorage.setItem(SAVED_BARISTA_QR_SETTINGS_KEY, JSON.stringify(values));
  toast({
    title: "Успех",
    description: "Размер для принтера QR-кодов сохранен",
  });
  isModalOpen.value = false;
});
</script>
