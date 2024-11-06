<script setup lang="ts">
import LocaleSwitch from "@/core/components/locale-switch/LocaleSwitch.vue"
</script>

<template>
	<div class="bg-background w-full h-screen overflow-x-hidden relative">
		<div class="absolute right-4 top-4"><LocaleSwitch /></div>
		<main class="p-6 w-full h-full">
			<RouterView />
		</main>
	</div>
</template>

<style scoped></style>
