<script setup lang="ts"></script>

<template>
	<div class="w-full min-h-screen overflow-x-hidden no-scrollbar">
		<main class="relative w-full h-full pb-32">
			<RouterView />
		</main>
	</div>
</template>

<style lang="scss">
body {
  background: radial-gradient(circle, #F5F5F7 80%, rgb(228, 253, 226) 100%);
  padding-top: env(safe-area-inset-top);
}
</style>
