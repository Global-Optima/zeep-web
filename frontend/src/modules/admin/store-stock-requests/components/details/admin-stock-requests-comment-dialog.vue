<template>
	<Dialog
		:open="open"
		@close="$emit('close')"
	>
		<DialogHeader>
			<DialogTitle>Укажите причину</DialogTitle>
		</DialogHeader>
		<DialogContent>
			<Textarea
				v-model="comment"
				placeholder="Причина отклонения"
				rows="3"
			/>
		</DialogContent>
		<DialogFooter>
			<Button
				variant="outline"
				@click="$emit('close')"
				>Отмена</Button
			>
			<Button
				variant="default"
				@click="submitComment"
			>
				ОК
			</Button>
		</DialogFooter>
	</Dialog>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/core/components/ui/dialog'
import { Textarea } from '@/core/components/ui/textarea'
import type { RejectStockRequestStatusDTO } from '@/modules/admin/store-stock-requests/models/stock-requests.model'
import { ref } from 'vue'

 defineProps<{
  open: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'submit', payload: RejectStockRequestStatusDTO): void
}>()

const comment = ref('')

function submitComment() {
  emit('submit', { comment: comment.value })
}
</script>
