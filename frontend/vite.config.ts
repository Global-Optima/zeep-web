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
