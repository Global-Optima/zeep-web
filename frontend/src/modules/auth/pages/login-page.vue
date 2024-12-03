<template>
	<div class="flex justify-center items-center bg-gray-100 w-full h-full">
		<Card class="">
			<CardHeader>
				<CardTitle class="text-xl sm:text-2xl">Вход для сотрудника</CardTitle>
				<CardDescription>Введите данные для входа в портал сотрудника</CardDescription>
			</CardHeader>
			<CardContent>
				<LoginForm />
			</CardContent>
		</Card>
	</div>
</template>

<script setup lang="ts">
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader
} from '@/core/components/ui/card'
import LoginForm from '@/modules/auth/components/login/login-form.vue'
</script>
