<template>
	<Dialog
		:open="open"
		@update:open="$emit('close')"
	>
		<DialogContent :include-close-button="false">
			<DialogHeader>
				<DialogTitle>Укажите причину</DialogTitle>
			</DialogHeader>

			<Textarea
				v-model="comment"
				placeholder="Причина отклонения"
				rows="3"
			/>

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
		</DialogContent>
	</Dialog>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/core/components/ui/dialog'
import { Textarea } from '@/core/components/ui/textarea'
import type { RejectStockRequestStatusDTO } from '@/modules/admin/stock-requests/models/stock-requests.model'
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
