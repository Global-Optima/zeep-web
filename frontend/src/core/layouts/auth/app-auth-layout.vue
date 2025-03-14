<script setup lang="ts">
import LocaleSwitch from "@/core/components/locale-switch/LocaleSwitch.vue"
import PageLoader from '@/core/components/page-loader/PageLoader.vue'
</script>

<template>
	<div class="relative bg-background w-full h-screen overflow-x-hidden">
		<div class="top-4 right-4 absolute"><LocaleSwitch /></div>
		<main class="w-full h-full">
			<RouterView v-slot="{ Component, route }">
				<template v-if="Component">
					<Transition
						name="fade"
						mode="out-in"
					>
						<Suspense>
							<div :key="route.matched[0]?.path">
								<component :is="Component" />
							</div>
							<template #fallback>
								<PageLoader />
							</template>
						</Suspense>
					</Transition>
				</template>
			</RouterView>
		</main>
	</div>
</template>

<style scoped></style>
