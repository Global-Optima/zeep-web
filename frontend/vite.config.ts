import vue from '@vitejs/plugin-vue'
import autoprefixer from 'autoprefixer'
import { fileURLToPath, URL } from 'node:url'
import tailwind from 'tailwindcss'
import { defineConfig } from 'vite'
import viteCompression from 'vite-plugin-compression'
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
		VitePWA({
			registerType: 'autoUpdate',
			manifest: {
				name: 'ZEEP',
				short_name: 'ZEEP',
				description:
					'Новая эра кофеен и чайных домов, где скорость и цифровой подход становятся синонимом удобства.',
				background_color: '#F5F5F7',
				theme_color: '#F5F5F7',
				display: 'standalone',
				lang: 'ru',
				start_url: '/',
				icons: [
					{ src: '/android-icon-36x36.png', sizes: '36x36', type: 'image/png' },
					{ src: '/android-icon-48x48.png', sizes: '48x48', type: 'image/png' },
					{ src: '/android-icon-72x72.png', sizes: '72x72', type: 'image/png' },
					{ src: '/android-icon-96x96.png', sizes: '96x96', type: 'image/png' },
					{ src: '/android-icon-144x144.png', sizes: '144x144', type: 'image/png' },
					{ src: '/android-icon-192x192.png', sizes: '192x192', type: 'image/png' },
				],
			},
			devOptions: {
				enabled: true,
				type: 'module',
			},
		}),
		viteCompression({
			algorithm: 'gzip',
			ext: '.gz',
			threshold: 10240,
		}),
		viteCompression({
			algorithm: 'brotliCompress',
			ext: '.br',
			threshold: 10240,
		}),
	],
	resolve: {
		alias: {
			'@': fileURLToPath(new URL('./src', import.meta.url)),
		},
	},
	build: {
		chunkSizeWarningLimit: 500,
		rollupOptions: {
			output: {
				manualChunks: {
					vendor: ['vue', 'vue-router', 'pinia'],
				},
			},
		},
	},
	server: {
		hmr: {
			protocol: 'ws',
			host: 'localhost',
			port: 5173,
			path: '/vite-hmr',
			timeout: 30000,
			overlay: false,
		},
		watch: {
			usePolling: true,
		},
	},
})
