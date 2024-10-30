<script setup lang="ts">
import AppAdminLayout from "@/core/layouts/admin/app-admin-layout.vue";
import AppDefaultLayout from "@/core/layouts/default/app-default-layout.vue";
import AppKioskLayout from "@/core/layouts/kiosk/app-kiosk-layout.vue";
import { AppLayouts } from "@/core/types/routes.types";
import { useColorMode } from "@vueuse/core";

useColorMode()
</script>

<template>
	<AppKioskLayout v-if="$route.meta.layout === AppLayouts.KIOSK">
		<RouterView />
	</AppKioskLayout>

	<AppAdminLayout v-else-if="$route.meta.layout === AppLayouts.ADMIN">
		<RouterView />
	</AppAdminLayout>

	<AppDefaultLayout v-else>
		<RouterView />
	</AppDefaultLayout>
</template>
