<template>
	<Tabs default-value="details">
		<TabsList class="grid w-full grid-cols-2 mb-6">
			<TabsTrigger
				class="py-2"
				value="details"
			>
				Детали
			</TabsTrigger>
			<TabsTrigger
				class="py-2"
				value="variants"
			>
				Варианты
			</TabsTrigger>
		</TabsList>
		<TabsContent value="details">
			<AdminProductDetails />
		</TabsContent>

		<TabsContent value="variants"> <AdminProductsVariants /></TabsContent>
	</Tabs>
</template>

<script setup lang="ts">
import {
	Tabs,
	TabsContent,
	TabsList,
	TabsTrigger,
} from '@/core/components/ui/tabs'
import AdminProductDetails from "@/modules/admin/products/components/details/admin-product-details.vue";
import AdminProductsVariants from "@/modules/admin/products/components/details/admin-products-variants.vue";
</script>
