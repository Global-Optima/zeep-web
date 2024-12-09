<script setup lang="ts">
import {
	Tabs,
	TabsContent,
	TabsList,
	TabsTrigger,
} from '@/core/components/ui/tabs'
import LoginForm from "@/modules/auth/components/login/store-login-form.vue"
</script>

<template>
	<div class="w-full h-full flex items-center justify-center bg-gray-200">
		<Tabs
			default-value="cafe"
			class="w-[400px]"
		>
			<TabsList class="grid w-full grid-cols-2">
				<TabsTrigger
					value="cafe"
					class="text-base"
				>
					Кафе
				</TabsTrigger>
				<TabsTrigger
					value="warehouse"
					class="text-base"
				>
					Склад
				</TabsTrigger>
			</TabsList>
			<TabsContent value="cafe">
				<LoginForm />
			</TabsContent>
			<TabsContent value="warehouse">
				<LoginForm />
			</TabsContent>
		</Tabs>
	</div>
</template>
