<script setup lang="ts">
import PageLoader from '@/core/components/page-loader/PageLoader.vue'
import { Toaster } from '@/core/components/ui/toast'
</script>

<template>
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

	<Toaster />
</template>
