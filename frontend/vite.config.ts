import { fileURLToPath, URL } from 'node:url'

import VueI18nPlugin from '@intlify/unplugin-vue-i18n/vite'
import vue from '@vitejs/plugin-vue'
import autoprefixer from 'autoprefixer'
import { dirname, resolve } from 'node:path'
import tailwind from 'tailwindcss'
import { defineConfig } from 'vite'
import { VitePWA } from 'vite-plugin-pwa'

export default defineConfig({
	css: {
		postcss: {
			plugins: [tailwind(), autoprefixer()],
		},
		preprocessorOptions: {
			scss: {
				api: 'modern-compiler',
			},
		},
	},
	plugins: [
		vue(),
		VueI18nPlugin({
			include: resolve(dirname(fileURLToPath(import.meta.url)), './src/core/locales/**'),
		}),
		VitePWA({
			registerType: 'autoUpdate',
			injectRegister: 'auto',
			workbox: {
				runtimeCaching: [
					{
						urlPattern: ({ request }) =>
							request.destination === 'style' ||
							request.destination === 'script' ||
							request.destination === 'worker',
						handler: 'StaleWhileRevalidate',
						options: {
							cacheName: 'static-resources',
							expiration: {
								maxEntries: 50,
								maxAgeSeconds: 30 * 24 * 60 * 60, // 30 days
							},
						},
					},
					{
						urlPattern: ({ request }) => request.destination === 'image',
						handler: 'CacheFirst',
						options: {
							cacheName: 'images',
							expiration: {
								maxEntries: 100,
								maxAgeSeconds: 60 * 24 * 60 * 60, // 60 days
							},
						},
					},
				],
			},
			devOptions: {
				enabled: true,
			},
		}),
	],
	resolve: {
		alias: {
			'@': fileURLToPath(new URL('./src', import.meta.url)),
		},
	},
})
