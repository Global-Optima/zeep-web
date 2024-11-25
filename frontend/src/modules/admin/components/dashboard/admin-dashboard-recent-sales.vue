<script setup lang="ts">
import {
  Avatar,
  AvatarFallback,
  AvatarImage,
} from '@/core/components/ui/avatar'
</script>

<template>
	<div class="space-y-8">
		<div class="flex items-center">
			<Avatar class="w-9 h-9">
				<AvatarImage
					src="/avatars/01.png"
					alt="Avatar"
				/>
				<AvatarFallback>OM</AvatarFallback>
			</Avatar>
			<div class="space-y-1 ml-4">
				<p class="font-medium text-sm leading-none">Olivia Martin</p>
				<p class="text-muted-foreground text-sm">olivia.martin@email.com</p>
			</div>
			<div class="ml-auto font-medium">+$1,999.00</div>
		</div>
		<div class="flex items-center">
			<Avatar class="flex justify-center items-center space-y-0 border w-9 h-9">
				<AvatarImage
					src="/avatars/02.png"
					alt="Avatar"
				/>
				<AvatarFallback>JL</AvatarFallback>
			</Avatar>
			<div class="space-y-1 ml-4">
				<p class="font-medium text-sm leading-none">Jackson Lee</p>
				<p class="text-muted-foreground text-sm">jackson.lee@email.com</p>
			</div>
			<div class="ml-auto font-medium">+$39.00</div>
		</div>
		<div class="flex items-center">
			<Avatar class="w-9 h-9">
				<AvatarImage
					src="/avatars/03.png"
					alt="Avatar"
				/>
				<AvatarFallback>IN</AvatarFallback>
			</Avatar>
			<div class="space-y-1 ml-4">
				<p class="font-medium text-sm leading-none">Isabella Nguyen</p>
				<p class="text-muted-foreground text-sm">isabella.nguyen@email.com</p>
			</div>
			<div class="ml-auto font-medium">+$299.00</div>
		</div>
		<div class="flex items-center">
			<Avatar class="w-9 h-9">
				<AvatarImage
					src="/avatars/04.png"
					alt="Avatar"
				/>
				<AvatarFallback>WK</AvatarFallback>
			</Avatar>
			<div class="space-y-1 ml-4">
				<p class="font-medium text-sm leading-none">William Kim</p>
				<p class="text-muted-foreground text-sm">will@email.com</p>
			</div>
			<div class="ml-auto font-medium">+$99.00</div>
		</div>
		<div class="flex items-center">
			<Avatar class="w-9 h-9">
				<AvatarImage
					src="/avatars/05.png"
					alt="Avatar"
				/>
				<AvatarFallback>SD</AvatarFallback>
			</Avatar>
			<div class="space-y-1 ml-4">
				<p class="font-medium text-sm leading-none">Sofia Davis</p>
				<p class="text-muted-foreground text-sm">sofia.davis@email.com</p>
			</div>
			<div class="ml-auto font-medium">+$39.00</div>
		</div>
	</div>
</template>
