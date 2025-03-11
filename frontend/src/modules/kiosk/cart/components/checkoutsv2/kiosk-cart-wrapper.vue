<script setup lang="ts">
import { useToast } from '@/core/components/ui/toast'
import { getRouteName } from '@/core/config/routes.config'
import type { LocalizedError } from '@/core/models/errors.model'
import type { CreateOrderDTO, OrderDTO } from '@/modules/admin/store-orders/models/orders.models'
import { ordersService } from '@/modules/admin/store-orders/services/orders.service'
import { CheckoutStep, type PaymentMethod } from '@/modules/kiosk/cart/models/kiosk-cart.models'
import { useCartStore } from '@/modules/kiosk/cart/stores/cart.store'
import { useMutation } from '@tanstack/vue-query'
import type { AxiosError } from 'axios'
import { defineAsyncComponent, onMounted, shallowReactive } from 'vue'
import { useRouter } from 'vue-router'


// Step State (Persisted)
const stepState = shallowReactive({
  currentStep: (sessionStorage.getItem('currentStep') as CheckoutStep) || CheckoutStep.CUSTOMER,
  customerName: sessionStorage.getItem('customerName') || '',
  selectedPayment: sessionStorage.getItem('selectedPayment') ?? null ,
  orderId: Number(sessionStorage.getItem('orderId')) || null,
  qrCodeUrl: '',
});

const cartStore = useCartStore();
const { toast } = useToast();
const router = useRouter();

// Lazy-loaded Step Components
const CheckoutCustomer = defineAsyncComponent(() =>
  import('@/modules/kiosk/cart/components/checkoutsv2/kiosk-cart-customer.vue')
);
const CheckoutPayment = defineAsyncComponent(() =>
  import('@/modules/kiosk/cart/components/checkoutsv2/kiosk-cart-payment.vue')
);
const CheckoutConfirmation = defineAsyncComponent(() =>
  import('@/modules/kiosk/cart/components/checkoutsv2/kiosk-cart-success.vue')
);

// Create Order Mutation
const createOrderMutation = useMutation<OrderDTO, AxiosError<LocalizedError>, CreateOrderDTO>({
  mutationFn: (orderDTO: CreateOrderDTO) => ordersService.createOrder(orderDTO),
  onSuccess: (order) => {
    stepState.orderId = order.id;
    sessionStorage.setItem('orderId', order.id.toString());
    stepState.currentStep = CheckoutStep.PAYMENT;
    sessionStorage.setItem('currentStep', CheckoutStep.PAYMENT);
  },
  onError: (error) => {
    toast({
      title: 'Ошибка',
      description: error.response?.data.message.ru ?? 'Ошибка при создании заказа',
    });
  },
});

// If there's an existing orderId in sessionStorage, jump to Payment on reload
onMounted(() => {
  if (stepState.orderId) {
    stepState.currentStep = CheckoutStep.PAYMENT;
  }
});

// Step Navigation
const proceedToNextStep = async () => {
  if (stepState.currentStep === CheckoutStep.CUSTOMER) {
    // Build CreateOrderDTO from cart
    const orderDTO: CreateOrderDTO = {
      customerName: stepState.customerName,
      subOrders: Object.values(cartStore.cartItems).map(item => ({
        storeProductSizeId: item.size.id,
        quantity: item.quantity,
        storeAdditivesIds: item.additives.map(add => add.id),
      })),
    };
    await createOrderMutation.mutateAsync(orderDTO);

    // If success, onSuccess sets stepState to 'payment'
  }
};

// Called when Payment is done or user closes
const resetCheckoutFlow = () => {
  stepState.currentStep = CheckoutStep.CUSTOMER;
  sessionStorage.clear();
  cartStore.clearCart();
  cartStore.toggleModal();
  router.push({ name: getRouteName("KIOSK_HOME") });
};

const closeCheckoutModal = () => {
  stepState.currentStep = CheckoutStep.CUSTOMER;
  sessionStorage.clear();
  cartStore.toggleModal();
}

// If you want to proceed from Payment → Confirmation instead of closing
const onPaymentSuccess = () => {
  // Example: show a Confirmation step with a QR code
  stepState.currentStep = CheckoutStep.CONFIRMATION;
  sessionStorage.setItem('currentStep', CheckoutStep.CONFIRMATION);

  // For demonstration, we set a sample QR code
  stepState.qrCodeUrl = 'https://cabinet.kofd.kz/consumer';
};

const updateSelectedPayment = (val: PaymentMethod | null) => {
  stepState.selectedPayment = val;
  sessionStorage.setItem('selectedPayment', val || '');
}
</script>

<template>
	<!-- Step: Customer -->
	<CheckoutCustomer
		v-if="stepState.currentStep === 'customer'"
		:open="true"
		v-model:customer-name="stepState.customerName"
		@toggle="closeCheckoutModal"
		@next="proceedToNextStep"
	/>

	<!-- Step: Payment -->
	<CheckoutPayment
		v-else-if="stepState.currentStep === 'payment'"
		:isOpen="true"
		:selectedPayment="stepState.selectedPayment as PaymentMethod ?? null"
		:orderId="stepState.orderId"
		@back="stepState.currentStep = CheckoutStep.CUSTOMER"
		@close="closeCheckoutModal"
		@update:selectedPayment="updateSelectedPayment"
		@proceed="onPaymentSuccess"
	/>

	<!-- Step: Confirmation -->
	<CheckoutConfirmation
		v-else-if="stepState.currentStep === 'confirmation'"
		:isOpen="true"
		:qrCodeUrl="stepState.qrCodeUrl"
		@close="resetCheckoutFlow"
		@proceed="resetCheckoutFlow"
	/>
</template>
