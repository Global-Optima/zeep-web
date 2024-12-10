<template>
	<div class="relative bg-white p-6 border rounded-xl">
		<!-- Employee Image -->
		<img
			class="mb-6 rounded-full w-48 h-48 object-cover"
			:src="employee.image"
			alt="Employee Image"
		/>
		<!-- Employee Attributes -->
		<div
			v-for="attribute in attributes"
			:key="attribute.label"
			class="mt-4"
		>
			<label class="block text-gray-400 text-sm">
				{{ attribute.label }}
			</label>
			<p class="mt-1 text-base text-gray-900">{{ attribute.value }}</p>
		</div>
		<!-- Action Buttons -->
		<div class="top-6 right-6 absolute">
			<Button
				size="icon"
				variant="ghost"
			>
				<Pencil class="w-6 h-6 text-gray-500" />
			</Button>
		</div>
	</div>
</template>

<script setup lang="ts">
import Button from '@/core/components/ui/button/Button.vue'
import { Pencil } from 'lucide-vue-next'
import { computed, ref } from 'vue'

interface Employee {
  image: string;
  name: string;
  position: string;
  email: string;
  phone: string;
  status: string;
  storeAddress: string;
  isActive: boolean;
}

const employee = ref<Employee>({
  image: 'https://via.placeholder.com/400x300',
  name: 'Иван Иванов',
  position: 'Менеджер по продажам',
  email: 'ivan.ivanov@example.com',
  phone: '+7 (495) 123-45-67',
  status: 'Активен',
  storeAddress: 'ул. Ленина, 123, Москва, Россия',
  isActive: true,
});

const attributes = computed(() => [
  { label: 'Имя', value: employee.value.name },
  { label: 'Должность', value: employee.value.position },
  { label: 'Email', value: employee.value.email },
  { label: 'Телефон', value: employee.value.phone },
  { label: 'Статус', value: employee.value.status },
  { label: 'Адрес магазина', value: employee.value.storeAddress },
]);
</script>

<style scoped></style>
