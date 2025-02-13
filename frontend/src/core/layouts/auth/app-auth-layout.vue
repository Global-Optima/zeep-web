<script setup lang="ts">
import LocaleSwitch from "@/core/components/locale-switch/LocaleSwitch.vue"
</script>

<template>
	<div class="relative bg-background w-full h-screen overflow-x-hidden">
		<div class="top-4 right-4 absolute"><LocaleSwitch /></div>
		<main class="w-full h-full">
			<RouterView />
		</main>
	</div>
</template>

<style scoped></style>
