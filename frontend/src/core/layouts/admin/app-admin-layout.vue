<script setup lang="ts">

import { Button } from '@/core/components/ui/button'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from '@/core/components/ui/dropdown-menu'
import { Input } from '@/core/components/ui/input'
import { Sheet, SheetContent, SheetTrigger } from '@/core/components/ui/sheet'
import { getRouteName } from '@/core/config/routes.config'
import AppAdminSidebar from '@/core/layouts/admin/app-admin-sidebar.vue'
import { CircleUser, Coffee, Menu, Package2, Search, Store } from 'lucide-vue-next'
import { useRouter } from 'vue-router'


const router = useRouter()

const onKioskClick = () => {
	router.push({name: getRouteName('KIOSK_HOME')})
}
const onBaristaClick = () => {
	router.push({name: getRouteName('KIOSK_ORDERS')})
}
</script>

<template>
	<div class="grid md:grid-cols-[220px_1fr] lg:grid-cols-[280px_1fr] bg-white w-full min-h-screen">
		<div class="md:block hidden bg-muted/40 border-r">
			<div class="flex flex-col gap-2 h-full max-h-screen">
				<div class="flex items-center px-4 lg:px-6 border-b h-14 lg:h-[60px]">
					<a
						href="/admin"
						class="flex items-center gap-2 font-semibold"
					>
						<Package2 class="w-6 h-6" />
						<span class="">Zeep</span>
					</a>
				</div>

				<div class="flex-1">
					<AppAdminSidebar />
				</div>

				<div class="mt-auto p-4">
					<Button
						variant="outline"
						@click="onBaristaClick"
						class="flex items-center gap-3 mb-2 py-6 w-full"
					>
						<Coffee class="w-5 h-5" />
						<p class="text-base">Бариста</p>
					</Button>

					<Button
						@click="onKioskClick"
						class="flex items-center gap-3 py-6 w-full"
					>
						<Store class="w-5 h-5" />
						<p class="text-base">Киоск</p>
					</Button>
				</div>
			</div>
		</div>

		<div class="flex flex-col">
			<header class="flex items-center gap-4 bg-muted/40 px-4 lg:px-6 border-b h-14 lg:h-[60px]">
				<Sheet>
					<SheetTrigger as-child>
						<Button
							variant="outline"
							size="icon"
							class="md:hidden shrink-0"
						>
							<Menu class="w-5 h-5" />
							<span class="sr-only">Toggle navigation menu</span>
						</Button>
					</SheetTrigger>

					<SheetContent
						side="left"
						class="flex flex-col"
					>
						<AppAdminSidebar />

						<div class="mt-auto">
							<Button
								variant="outline"
								@click="onBaristaClick"
								class="flex items-center gap-3 mb-1 py-6 w-full"
							>
								<Store class="w-5 h-5" />
								<p class="text-base">Бариста</p>
							</Button>

							<Button
								@click="onKioskClick"
								class="flex items-center gap-3 py-6 w-full"
							>
								<Store class="w-5 h-5" />
								<p class="text-base">Киоск</p>
							</Button>
						</div>
					</SheetContent>
				</Sheet>
				<div class="flex-1 w-full">
					<form>
						<div class="relative">
							<Search class="top-2.5 left-2.5 absolute w-4 h-4 text-muted-foreground" />
							<Input
								type="search"
								placeholder="Поиск"
								class="bg-background shadow-none pl-8 w-full md:w-2/3 lg:w-1/3 appearance-none"
							/>
						</div>
					</form>
				</div>
				<DropdownMenu>
					<DropdownMenuTrigger as-child>
						<Button
							variant="secondary"
							size="icon"
							class="rounded-full"
						>
							<CircleUser class="w-5 h-5" />
							<span class="sr-only">Toggle user menu</span>
						</Button>
					</DropdownMenuTrigger>
					<DropdownMenuContent align="end">
						<DropdownMenuLabel>My Account</DropdownMenuLabel>
						<DropdownMenuSeparator />
						<DropdownMenuItem>Settings</DropdownMenuItem>
						<DropdownMenuItem>Support</DropdownMenuItem>
						<DropdownMenuSeparator />
						<DropdownMenuItem>Logout</DropdownMenuItem>
					</DropdownMenuContent>
				</DropdownMenu>
			</header>
			<main class="flex flex-col flex-1 gap-4 lg:gap-6 p-4 lg:p-6">
				<div class="flex items-center">
					<h1 class="font-semibold text-lg md:text-2xl">Inventory</h1>
				</div>
				<div
					class="flex flex-1 justify-center items-center shadow-sm border border-dashed rounded-lg"
				>
					<div class="flex flex-col items-center gap-1 text-center">
						<h3 class="font-bold text-2xl tracking-tight">You have no products</h3>
						<p class="text-muted-foreground text-sm">
							You can start selling as soon as you add a product.
						</p>
						<Button class="mt-4"> Add Product </Button>
					</div>
				</div>
			</main>
		</div>
	</div>
</template>
