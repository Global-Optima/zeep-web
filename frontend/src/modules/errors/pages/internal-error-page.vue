<template>
	<div class="flex justify-center items-center p-4 min-h-screen text-center">
		<div class="max-w-[26rem]">
			<h1 class="mb-4 font-semibold text-2xl">Внутренняя ошибка сервера</h1>
			<p class="mb-6 text-gray-600">
				Произошла непредвиденная ошибка. Попробуйте обновить страницу или вернитесь на главную,
				чтобы продолжить работу.
			</p>
			<Button @click="$router.push('/')">На главную</Button>
		</div>
	</div>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
</script>

<style scoped>
/* Добавьте свой дополнительный стиль при необходимости */
</style>
