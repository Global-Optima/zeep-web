// vite.config.ts
import { fileURLToPath, URL } from "node:url";
import vue from "file:///C:/Users/Diar/Desktop/ggnetworks/zeep-web/frontend/node_modules/@vitejs/plugin-vue/dist/index.mjs";
import autoprefixer from "file:///C:/Users/Diar/Desktop/ggnetworks/zeep-web/frontend/node_modules/autoprefixer/lib/autoprefixer.js";
import tailwind from "file:///C:/Users/Diar/Desktop/ggnetworks/zeep-web/frontend/node_modules/tailwindcss/lib/index.js";
import { defineConfig } from "file:///C:/Users/Diar/Desktop/ggnetworks/zeep-web/frontend/node_modules/vite/dist/node/index.js";
import viteCompression from "file:///C:/Users/Diar/Desktop/ggnetworks/zeep-web/frontend/node_modules/vite-plugin-compression/dist/index.mjs";
import { VitePWA } from "file:///C:/Users/Diar/Desktop/ggnetworks/zeep-web/frontend/node_modules/vite-plugin-pwa/dist/index.js";
var __vite_injected_original_import_meta_url = "file:///C:/Users/Diar/Desktop/ggnetworks/zeep-web/frontend/vite.config.ts";
var vite_config_default = defineConfig({
  css: {
    postcss: {
      plugins: [tailwind(), autoprefixer()]
    },
    preprocessorOptions: {
      scss: {
        api: "modern-compiler"
      }
    }
  },
  plugins: [
    vue(),
    VitePWA({
      registerType: "autoUpdate",
      injectRegister: "auto",
      workbox: {
        runtimeCaching: [
          {
            urlPattern: ({ request }) => request.destination === "style" || request.destination === "script" || request.destination === "worker",
            handler: "StaleWhileRevalidate",
            options: {
              cacheName: "static-resources",
              expiration: {
                maxEntries: 50,
                maxAgeSeconds: 30 * 24 * 60 * 60
                // 30 days
              }
            }
          },
          {
            urlPattern: ({ request }) => request.destination === "image",
            handler: "CacheFirst",
            options: {
              cacheName: "images",
              expiration: {
                maxEntries: 100,
                maxAgeSeconds: 60 * 24 * 60 * 60
                // 60 days
              }
            }
          }
        ]
      },
      devOptions: {
        enabled: true
      }
    }),
    viteCompression()
  ],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", __vite_injected_original_import_meta_url))
    }
  }
});
export {
  vite_config_default as default
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcudHMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbImNvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9kaXJuYW1lID0gXCJDOlxcXFxVc2Vyc1xcXFxEaWFyXFxcXERlc2t0b3BcXFxcZ2duZXR3b3Jrc1xcXFx6ZWVwLXdlYlxcXFxmcm9udGVuZFwiO2NvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9maWxlbmFtZSA9IFwiQzpcXFxcVXNlcnNcXFxcRGlhclxcXFxEZXNrdG9wXFxcXGdnbmV0d29ya3NcXFxcemVlcC13ZWJcXFxcZnJvbnRlbmRcXFxcdml0ZS5jb25maWcudHNcIjtjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfaW1wb3J0X21ldGFfdXJsID0gXCJmaWxlOi8vL0M6L1VzZXJzL0RpYXIvRGVza3RvcC9nZ25ldHdvcmtzL3plZXAtd2ViL2Zyb250ZW5kL3ZpdGUuY29uZmlnLnRzXCI7aW1wb3J0IHsgZmlsZVVSTFRvUGF0aCwgVVJMIH0gZnJvbSAnbm9kZTp1cmwnXHJcblxyXG5pbXBvcnQgdnVlIGZyb20gJ0B2aXRlanMvcGx1Z2luLXZ1ZSdcclxuaW1wb3J0IGF1dG9wcmVmaXhlciBmcm9tICdhdXRvcHJlZml4ZXInXHJcbmltcG9ydCB0YWlsd2luZCBmcm9tICd0YWlsd2luZGNzcydcclxuaW1wb3J0IHsgZGVmaW5lQ29uZmlnIH0gZnJvbSAndml0ZSdcclxuaW1wb3J0IHZpdGVDb21wcmVzc2lvbiBmcm9tICd2aXRlLXBsdWdpbi1jb21wcmVzc2lvbidcclxuaW1wb3J0IHsgVml0ZVBXQSB9IGZyb20gJ3ZpdGUtcGx1Z2luLXB3YSdcclxuXHJcbmV4cG9ydCBkZWZhdWx0IGRlZmluZUNvbmZpZyh7XHJcblx0Y3NzOiB7XHJcblx0XHRwb3N0Y3NzOiB7XHJcblx0XHRcdHBsdWdpbnM6IFt0YWlsd2luZCgpLCBhdXRvcHJlZml4ZXIoKV0sXHJcblx0XHR9LFxyXG5cdFx0cHJlcHJvY2Vzc29yT3B0aW9uczoge1xyXG5cdFx0XHRzY3NzOiB7XHJcblx0XHRcdFx0YXBpOiAnbW9kZXJuLWNvbXBpbGVyJyxcclxuXHRcdFx0fSxcclxuXHRcdH0sXHJcblx0fSxcclxuXHRwbHVnaW5zOiBbXHJcblx0XHR2dWUoKSxcclxuXHRcdFZpdGVQV0Eoe1xyXG5cdFx0XHRyZWdpc3RlclR5cGU6ICdhdXRvVXBkYXRlJyxcclxuXHRcdFx0aW5qZWN0UmVnaXN0ZXI6ICdhdXRvJyxcclxuXHRcdFx0d29ya2JveDoge1xyXG5cdFx0XHRcdHJ1bnRpbWVDYWNoaW5nOiBbXHJcblx0XHRcdFx0XHR7XHJcblx0XHRcdFx0XHRcdHVybFBhdHRlcm46ICh7IHJlcXVlc3QgfSkgPT5cclxuXHRcdFx0XHRcdFx0XHRyZXF1ZXN0LmRlc3RpbmF0aW9uID09PSAnc3R5bGUnIHx8XHJcblx0XHRcdFx0XHRcdFx0cmVxdWVzdC5kZXN0aW5hdGlvbiA9PT0gJ3NjcmlwdCcgfHxcclxuXHRcdFx0XHRcdFx0XHRyZXF1ZXN0LmRlc3RpbmF0aW9uID09PSAnd29ya2VyJyxcclxuXHRcdFx0XHRcdFx0aGFuZGxlcjogJ1N0YWxlV2hpbGVSZXZhbGlkYXRlJyxcclxuXHRcdFx0XHRcdFx0b3B0aW9uczoge1xyXG5cdFx0XHRcdFx0XHRcdGNhY2hlTmFtZTogJ3N0YXRpYy1yZXNvdXJjZXMnLFxyXG5cdFx0XHRcdFx0XHRcdGV4cGlyYXRpb246IHtcclxuXHRcdFx0XHRcdFx0XHRcdG1heEVudHJpZXM6IDUwLFxyXG5cdFx0XHRcdFx0XHRcdFx0bWF4QWdlU2Vjb25kczogMzAgKiAyNCAqIDYwICogNjAsIC8vIDMwIGRheXNcclxuXHRcdFx0XHRcdFx0XHR9LFxyXG5cdFx0XHRcdFx0XHR9LFxyXG5cdFx0XHRcdFx0fSxcclxuXHRcdFx0XHRcdHtcclxuXHRcdFx0XHRcdFx0dXJsUGF0dGVybjogKHsgcmVxdWVzdCB9KSA9PiByZXF1ZXN0LmRlc3RpbmF0aW9uID09PSAnaW1hZ2UnLFxyXG5cdFx0XHRcdFx0XHRoYW5kbGVyOiAnQ2FjaGVGaXJzdCcsXHJcblx0XHRcdFx0XHRcdG9wdGlvbnM6IHtcclxuXHRcdFx0XHRcdFx0XHRjYWNoZU5hbWU6ICdpbWFnZXMnLFxyXG5cdFx0XHRcdFx0XHRcdGV4cGlyYXRpb246IHtcclxuXHRcdFx0XHRcdFx0XHRcdG1heEVudHJpZXM6IDEwMCxcclxuXHRcdFx0XHRcdFx0XHRcdG1heEFnZVNlY29uZHM6IDYwICogMjQgKiA2MCAqIDYwLCAvLyA2MCBkYXlzXHJcblx0XHRcdFx0XHRcdFx0fSxcclxuXHRcdFx0XHRcdFx0fSxcclxuXHRcdFx0XHRcdH0sXHJcblx0XHRcdFx0XSxcclxuXHRcdFx0fSxcclxuXHRcdFx0ZGV2T3B0aW9uczoge1xyXG5cdFx0XHRcdGVuYWJsZWQ6IHRydWUsXHJcblx0XHRcdH0sXHJcblx0XHR9KSxcclxuXHRcdHZpdGVDb21wcmVzc2lvbigpLFxyXG5cdF0sXHJcblx0cmVzb2x2ZToge1xyXG5cdFx0YWxpYXM6IHtcclxuXHRcdFx0J0AnOiBmaWxlVVJMVG9QYXRoKG5ldyBVUkwoJy4vc3JjJywgaW1wb3J0Lm1ldGEudXJsKSksXHJcblx0XHR9LFxyXG5cdH0sXHJcbn0pXHJcbiJdLAogICJtYXBwaW5ncyI6ICI7QUFBc1YsU0FBUyxlQUFlLFdBQVc7QUFFelgsT0FBTyxTQUFTO0FBQ2hCLE9BQU8sa0JBQWtCO0FBQ3pCLE9BQU8sY0FBYztBQUNyQixTQUFTLG9CQUFvQjtBQUM3QixPQUFPLHFCQUFxQjtBQUM1QixTQUFTLGVBQWU7QUFQaU0sSUFBTSwyQ0FBMkM7QUFTMVEsSUFBTyxzQkFBUSxhQUFhO0FBQUEsRUFDM0IsS0FBSztBQUFBLElBQ0osU0FBUztBQUFBLE1BQ1IsU0FBUyxDQUFDLFNBQVMsR0FBRyxhQUFhLENBQUM7QUFBQSxJQUNyQztBQUFBLElBQ0EscUJBQXFCO0FBQUEsTUFDcEIsTUFBTTtBQUFBLFFBQ0wsS0FBSztBQUFBLE1BQ047QUFBQSxJQUNEO0FBQUEsRUFDRDtBQUFBLEVBQ0EsU0FBUztBQUFBLElBQ1IsSUFBSTtBQUFBLElBQ0osUUFBUTtBQUFBLE1BQ1AsY0FBYztBQUFBLE1BQ2QsZ0JBQWdCO0FBQUEsTUFDaEIsU0FBUztBQUFBLFFBQ1IsZ0JBQWdCO0FBQUEsVUFDZjtBQUFBLFlBQ0MsWUFBWSxDQUFDLEVBQUUsUUFBUSxNQUN0QixRQUFRLGdCQUFnQixXQUN4QixRQUFRLGdCQUFnQixZQUN4QixRQUFRLGdCQUFnQjtBQUFBLFlBQ3pCLFNBQVM7QUFBQSxZQUNULFNBQVM7QUFBQSxjQUNSLFdBQVc7QUFBQSxjQUNYLFlBQVk7QUFBQSxnQkFDWCxZQUFZO0FBQUEsZ0JBQ1osZUFBZSxLQUFLLEtBQUssS0FBSztBQUFBO0FBQUEsY0FDL0I7QUFBQSxZQUNEO0FBQUEsVUFDRDtBQUFBLFVBQ0E7QUFBQSxZQUNDLFlBQVksQ0FBQyxFQUFFLFFBQVEsTUFBTSxRQUFRLGdCQUFnQjtBQUFBLFlBQ3JELFNBQVM7QUFBQSxZQUNULFNBQVM7QUFBQSxjQUNSLFdBQVc7QUFBQSxjQUNYLFlBQVk7QUFBQSxnQkFDWCxZQUFZO0FBQUEsZ0JBQ1osZUFBZSxLQUFLLEtBQUssS0FBSztBQUFBO0FBQUEsY0FDL0I7QUFBQSxZQUNEO0FBQUEsVUFDRDtBQUFBLFFBQ0Q7QUFBQSxNQUNEO0FBQUEsTUFDQSxZQUFZO0FBQUEsUUFDWCxTQUFTO0FBQUEsTUFDVjtBQUFBLElBQ0QsQ0FBQztBQUFBLElBQ0QsZ0JBQWdCO0FBQUEsRUFDakI7QUFBQSxFQUNBLFNBQVM7QUFBQSxJQUNSLE9BQU87QUFBQSxNQUNOLEtBQUssY0FBYyxJQUFJLElBQUksU0FBUyx3Q0FBZSxDQUFDO0FBQUEsSUFDckQ7QUFBQSxFQUNEO0FBQ0QsQ0FBQzsiLAogICJuYW1lcyI6IFtdCn0K
