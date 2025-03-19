<template>
	<div class="flex justify-center items-center w-full h-[80vh] overflow-hidden">
		<div
			class="flex flex-col justify-center items-center space-y-6 bg-white shadow-2xl shadow-gray-200 mx-auto px-10 py-6 rounded-[36px] w-fit"
		>
			<!-- Profile Icon with Initials -->
			<div
				class="relative flex justify-center items-center bg-gray-200 bg-gradient-to-r rounded-full size-24 font-semibold text-gray-700 text-3xl"
			>
				{{ initials }}
			</div>

			<!-- Welcome Message -->
			<div>
				<h1 class="font-bold text-gray-900 dark:text-gray-100 text-2xl text-center">
					Добро пожаловать, {{ currentEmployee?.firstName }}!
				</h1>
				<p class="mt-2 text-gray-500 dark:text-gray-300 text-center">
					Вы вошли как
					<span class="font-semibold text-gray-800 dark:text-gray-200">{{ formattedRole }}</span>
				</p>
			</div>

			<!-- Employee Details -->
			<div class="bg-gray-50 dark:bg-gray-800 shadow-sm p-4 rounded-lg w-full max-w-md">
				<div class="flex justify-between items-center text-gray-600 dark:text-gray-300">
					<span>Телефон:</span>
					<span
						class="font-medium text-gray-900 dark:text-gray-100"
						>{{ currentEmployee?.phone }}</span
					>
				</div>
				<div class="flex justify-between items-center mt-2 text-gray-600 dark:text-gray-300">
					<span>Email:</span>
					<span
						class="font-medium text-gray-900 dark:text-gray-100"
						>{{ currentEmployee?.email }}</span
					>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { EMPLOYEE_ROLES_FORMATTED } from '@/modules/admin/employees/models/employees.models'
import { useEmployeeAuthStore } from '@/modules/auth/store/employee-auth.store'
import { computed } from 'vue'

// Получаем текущего пользователя из стора
const { currentEmployee } = useEmployeeAuthStore()

// Получение инициалов сотрудника
const initials = computed(() => {
	return `${currentEmployee?.firstName?.charAt(0) ?? ''}${currentEmployee?.lastName?.charAt(0) ?? ''}`.toUpperCase()
})

// Форматирование роли и типа сотрудника
const formattedRole = computed(() => currentEmployee ? EMPLOYEE_ROLES_FORMATTED[currentEmployee.role] : 'Пользователь')
</script>

<style scoped>
/* Добавить плавные переходы для более приятного UX */
div {
	transition: all 0.3s ease-in-out;
}
</style>
