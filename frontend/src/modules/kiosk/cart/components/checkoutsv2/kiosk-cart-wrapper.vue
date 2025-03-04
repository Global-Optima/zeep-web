<script setup lang="ts">
import { useToast } from '@/core/components/ui/toast'
import { getRouteName } from '@/core/config/routes.config'
import type { LocalizedError } from '@/core/models/errors.model'
import type { CreateOrderDTO } from '@/modules/admin/store-orders/models/orders.models'
import { ordersService } from '@/modules/admin/store-orders/services/orders.service'
import { useCartStore } from '@/modules/kiosk/cart/stores/cart.store'
import { useMutation } from '@tanstack/vue-query'
import type { AxiosError } from 'axios'
import { defineAsyncComponent, shallowReactive } from 'vue'
import { useRouter } from 'vue-router'

// Step State
const stepState = shallowReactive({
  currentStep: 'customer' as 'customer' | 'payment' | 'confirmation',
  customerName: '',
  selectedPayment: '',
  qrCodeUrl: '',
});

// Cart Store
const cartStore = useCartStore();
const {toast} = useToast()
const router = useRouter()

// Step Components
const CheckoutCustomer = defineAsyncComponent(() =>
  import('@/modules/kiosk/cart/components/checkoutsv2/kiosk-cart-customer.vue')
);
const CheckoutPayment = defineAsyncComponent(() =>
  import('@/modules/kiosk/cart/components/checkoutsv2/kiosk-cart-payment.vue')
);
const CheckoutConfirmation = defineAsyncComponent(() =>
  import('@/modules/kiosk/cart/components/checkoutsv2/kiosk-cart-success.vue')
);

// Mutations
const createOrderMutation = useMutation({
  mutationFn: (orderDTO: CreateOrderDTO) => ordersService.createOrder(orderDTO),
  onSuccess: () => {
    stepState.qrCodeUrl = "https://cabinet.kofd.kz/consumer";
    stepState.currentStep = 'confirmation';
  },
  onError: (error: AxiosError<LocalizedError>) => {
    toast({
      title: "Ошибка",
      description: error.response?.data.message.ru ?? "Ошибка при создании заказа",
    });
  },
});

const checkCustomerNameMutation = useMutation({
  mutationFn: (name: string) => ordersService.checkCustomerName({ customerName: name }),
  onError: (error: AxiosError<LocalizedError>) => {
    toast({
      title: "Ошибка",
      description: error.response?.data.message.ru ?? "Ошибка при проверке имени клиента",
    });
  },
});

// Navigation Logic
const proceedToNextStep = async () => {
  if (stepState.currentStep === 'customer') {
    try {
      await checkCustomerNameMutation.mutateAsync(stepState.customerName);
      stepState.currentStep = 'payment';
    } catch {
      return;
    }
  } else if (stepState.currentStep === 'payment') {
    const orderDTO: CreateOrderDTO = {
      customerName: stepState.customerName,
      subOrders: Object.values(cartStore.cartItems).map(item => ({
        storeProductSizeId: item.size.id,
        quantity: item.quantity,
        storeAdditivesIds: item.additives.map(add => add.id),
      })),
    };
    await createOrderMutation.mutateAsync(orderDTO);
  }
};

const goBack = () => {
  if (stepState.currentStep === 'payment') {
    stepState.currentStep = 'customer';
  } else if (stepState.currentStep === 'confirmation') {
    stepState.currentStep = 'payment';
  }
};

const closeSuccess = () => {
  stepState.currentStep = 'customer';
  stepState.customerName = '';
  stepState.selectedPayment = '';
  stepState.qrCodeUrl = '';
  cartStore.clearCart();
  cartStore.toggleModal()

  router.push({name: getRouteName("KIOSK_HOME")})
};


const closePayment = () => {
  stepState.currentStep = 'customer';
  stepState.customerName = '';
  stepState.selectedPayment = '';
  stepState.qrCodeUrl = '';
  cartStore.toggleModal()
};

const closeCustomer = () => {
  stepState.currentStep = 'customer';
  stepState.customerName = '';
  stepState.selectedPayment = '';
  stepState.qrCodeUrl = '';
  cartStore.toggleModal()
};
</script>

<template>
	<CheckoutCustomer
		v-if="stepState.currentStep === 'customer'"
		:open="true"
		@toggle="closeCustomer"
		@next="proceedToNextStep"
		v-model:customer-name="stepState.customerName"
	/>

	<CheckoutPayment
		v-if="stepState.currentStep === 'payment'"
		:isOpen="true"
		:selectedPayment="stepState.selectedPayment"
		@close="closePayment"
		@back="goBack"
		@proceed="(data) => { stepState.selectedPayment = data.paymentMethod; proceedToNextStep(); }"
	/>

	<CheckoutConfirmation
		v-if="stepState.currentStep === 'confirmation'"
		:isOpen="true"
		:qrCodeUrl="stepState.qrCodeUrl"
		@close="closeSuccess"
		@proceed="closeSuccess"
	/>
</template>
