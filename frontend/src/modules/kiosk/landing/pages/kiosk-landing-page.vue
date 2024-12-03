<template>
	<div
		class="relative w-full h-screen"
		@click="handleClick"
	>
		<!-- Language Selector -->
		<div class="top-8 right-8 z-10 absolute">
			<LocaleSwitch />
		</div>

		<swiper
			class="w-full h-screen"
			:modules="modules"
			:autoplay="{
        delay: 3000,
        disableOnInteraction: false,
      }"
		>
			<swiper-slide
				v-for="(slide, index) in slides"
				:key="index"
			>
				<!-- Slide Content -->
				<div class="relative w-full h-full">
					<!-- Gradient Overlay Behind Text -->
					<div
						class="bottom-0 left-0 absolute bg-gradient-to-t from-transparent via-transparent to-black w-full h-1/2"
					></div>

					<!-- Text Content -->
					<div class="bottom-64 left-12 z-10 absolute">
						<p class="text-2xl text-white sm:text-5xl">{{ slide.title }}</p>
						<p class="mt-2 sm:mt-4 font-medium text-3xl text-white sm:text-5xl">
							{{ formatPrice(slide.price) }}
						</p>
					</div>

					<!-- Image -->
					<img
						:src="slide.image"
						:alt="slide.title"
						class="w-full h-full object-cover"
						loading="eager"
					/>
				</div>
			</swiper-slide>
			<!-- Pagination -->
			<div class="swiper-pagination"></div>
		</swiper>

		<!-- Touch Indicator -->
		<div
			class="bottom-20 z-10 absolute flex justify-center items-center gap-3 w-full text-gray-400 animate-bounce"
		>
			<MousePointerClickIcon class="w-6 sm:w-10 h-6 sm:h-10" />
			<p class="text-base sm:text-3xl">Перейти в каталог</p>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

import { Swiper, SwiperSlide } from 'swiper/vue'

import 'swiper/css'
import 'swiper/css/pagination'

import { Autoplay, Pagination } from 'swiper/modules'

import LandingImage0 from '@/core/assets/landing/image0.avif'
import LandingImage1 from '@/core/assets/landing/image1.avif'
import LandingImage2 from '@/core/assets/landing/image2.avif'
import LandingImage3 from '@/core/assets/landing/image3.avif'
import LocaleSwitch from '@/core/components/locale-switch/LocaleSwitch.vue'
import { getRouteName } from '@/core/config/routes.config'
import { formatPrice } from '@/core/utils/price.utils'
import { MousePointerClickIcon } from 'lucide-vue-next'

const slides = ref([
  {
    image: LandingImage0,
    title: 'Капучино с карамельным сиропом',
    price: 1400,
  },
  {
    image: LandingImage3,
    title: 'Эспрессо Фраппучино',
    price: 1700,
  },
  {
    image: LandingImage2,
    title: 'Черный чай с лаймом',
    price: 650,
  },
  {
    image: LandingImage1,
    title: 'Грейпфрутовый лимонад',
    price: 1250,
  },
])

const modules = [Autoplay, Pagination]

const router = useRouter()

const handleClick = () => {
  router.push({ name: getRouteName("KIOSK_HOME") })
}
</script>

<style scoped>
/* Enlarge pagination indicators */
.swiper-pagination-bullet {
  width: 16px;
  height: 16px;
  background-color: white;
  opacity: 0.7;
}

.swiper-pagination-bullet-active {
  background-color: white;
  opacity: 1;
}

/* Touch indicator animation */
@keyframes bounce {
  0%, 100% {
    transform: translateY(-20%);
  }
  50% {
    transform: translateY(0);
  }
}

.animate-bounce {
  animation: bounce 1s infinite;
}

.bg-gradient-to-t {
  background: linear-gradient(to top, rgba(0, 0, 0, 0.7), transparent);
}

.safe-top-locale {
  padding-top: env(safe-area-inset-top);
}
</style>
