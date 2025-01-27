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
  server: {
    hmr: {
      host: "localhost",
      protocol: "ws",
      port: 5173
    }
  },
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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcudHMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbImNvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9kaXJuYW1lID0gXCJDOlxcXFxVc2Vyc1xcXFxEaWFyXFxcXERlc2t0b3BcXFxcZ2duZXR3b3Jrc1xcXFx6ZWVwLXdlYlxcXFxmcm9udGVuZFwiO2NvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9maWxlbmFtZSA9IFwiQzpcXFxcVXNlcnNcXFxcRGlhclxcXFxEZXNrdG9wXFxcXGdnbmV0d29ya3NcXFxcemVlcC13ZWJcXFxcZnJvbnRlbmRcXFxcdml0ZS5jb25maWcudHNcIjtjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfaW1wb3J0X21ldGFfdXJsID0gXCJmaWxlOi8vL0M6L1VzZXJzL0RpYXIvRGVza3RvcC9nZ25ldHdvcmtzL3plZXAtd2ViL2Zyb250ZW5kL3ZpdGUuY29uZmlnLnRzXCI7aW1wb3J0IHsgZmlsZVVSTFRvUGF0aCwgVVJMIH0gZnJvbSAnbm9kZTp1cmwnXHJcblxyXG5pbXBvcnQgdnVlIGZyb20gJ0B2aXRlanMvcGx1Z2luLXZ1ZSdcclxuaW1wb3J0IGF1dG9wcmVmaXhlciBmcm9tICdhdXRvcHJlZml4ZXInXHJcbmltcG9ydCB0YWlsd2luZCBmcm9tICd0YWlsd2luZGNzcydcclxuaW1wb3J0IHsgZGVmaW5lQ29uZmlnIH0gZnJvbSAndml0ZSdcclxuaW1wb3J0IHZpdGVDb21wcmVzc2lvbiBmcm9tICd2aXRlLXBsdWdpbi1jb21wcmVzc2lvbidcclxuaW1wb3J0IHsgVml0ZVBXQSB9IGZyb20gJ3ZpdGUtcGx1Z2luLXB3YSdcclxuXHJcbmV4cG9ydCBkZWZhdWx0IGRlZmluZUNvbmZpZyh7XHJcblx0c2VydmVyOiB7XHJcblx0XHRobXI6IHtcclxuXHRcdFx0aG9zdDogJ2xvY2FsaG9zdCcsXHJcblx0XHRcdHByb3RvY29sOiAnd3MnLFxyXG5cdFx0XHRwb3J0OiA1MTczLFxyXG5cdFx0fSxcclxuXHR9LFxyXG5cdGNzczoge1xyXG5cdFx0cG9zdGNzczoge1xyXG5cdFx0XHRwbHVnaW5zOiBbdGFpbHdpbmQoKSwgYXV0b3ByZWZpeGVyKCldLFxyXG5cdFx0fSxcclxuXHRcdHByZXByb2Nlc3Nvck9wdGlvbnM6IHtcclxuXHRcdFx0c2Nzczoge1xyXG5cdFx0XHRcdGFwaTogJ21vZGVybi1jb21waWxlcicsXHJcblx0XHRcdH0sXHJcblx0XHR9LFxyXG5cdH0sXHJcblx0cGx1Z2luczogW1xyXG5cdFx0dnVlKCksXHJcblx0XHRWaXRlUFdBKHtcclxuXHRcdFx0cmVnaXN0ZXJUeXBlOiAnYXV0b1VwZGF0ZScsXHJcblx0XHRcdGluamVjdFJlZ2lzdGVyOiAnYXV0bycsXHJcblx0XHRcdHdvcmtib3g6IHtcclxuXHRcdFx0XHRydW50aW1lQ2FjaGluZzogW1xyXG5cdFx0XHRcdFx0e1xyXG5cdFx0XHRcdFx0XHR1cmxQYXR0ZXJuOiAoeyByZXF1ZXN0IH0pID0+XHJcblx0XHRcdFx0XHRcdFx0cmVxdWVzdC5kZXN0aW5hdGlvbiA9PT0gJ3N0eWxlJyB8fFxyXG5cdFx0XHRcdFx0XHRcdHJlcXVlc3QuZGVzdGluYXRpb24gPT09ICdzY3JpcHQnIHx8XHJcblx0XHRcdFx0XHRcdFx0cmVxdWVzdC5kZXN0aW5hdGlvbiA9PT0gJ3dvcmtlcicsXHJcblx0XHRcdFx0XHRcdGhhbmRsZXI6ICdTdGFsZVdoaWxlUmV2YWxpZGF0ZScsXHJcblx0XHRcdFx0XHRcdG9wdGlvbnM6IHtcclxuXHRcdFx0XHRcdFx0XHRjYWNoZU5hbWU6ICdzdGF0aWMtcmVzb3VyY2VzJyxcclxuXHRcdFx0XHRcdFx0XHRleHBpcmF0aW9uOiB7XHJcblx0XHRcdFx0XHRcdFx0XHRtYXhFbnRyaWVzOiA1MCxcclxuXHRcdFx0XHRcdFx0XHRcdG1heEFnZVNlY29uZHM6IDMwICogMjQgKiA2MCAqIDYwLCAvLyAzMCBkYXlzXHJcblx0XHRcdFx0XHRcdFx0fSxcclxuXHRcdFx0XHRcdFx0fSxcclxuXHRcdFx0XHRcdH0sXHJcblx0XHRcdFx0XHR7XHJcblx0XHRcdFx0XHRcdHVybFBhdHRlcm46ICh7IHJlcXVlc3QgfSkgPT4gcmVxdWVzdC5kZXN0aW5hdGlvbiA9PT0gJ2ltYWdlJyxcclxuXHRcdFx0XHRcdFx0aGFuZGxlcjogJ0NhY2hlRmlyc3QnLFxyXG5cdFx0XHRcdFx0XHRvcHRpb25zOiB7XHJcblx0XHRcdFx0XHRcdFx0Y2FjaGVOYW1lOiAnaW1hZ2VzJyxcclxuXHRcdFx0XHRcdFx0XHRleHBpcmF0aW9uOiB7XHJcblx0XHRcdFx0XHRcdFx0XHRtYXhFbnRyaWVzOiAxMDAsXHJcblx0XHRcdFx0XHRcdFx0XHRtYXhBZ2VTZWNvbmRzOiA2MCAqIDI0ICogNjAgKiA2MCwgLy8gNjAgZGF5c1xyXG5cdFx0XHRcdFx0XHRcdH0sXHJcblx0XHRcdFx0XHRcdH0sXHJcblx0XHRcdFx0XHR9LFxyXG5cdFx0XHRcdF0sXHJcblx0XHRcdH0sXHJcblx0XHRcdGRldk9wdGlvbnM6IHtcclxuXHRcdFx0XHRlbmFibGVkOiB0cnVlLFxyXG5cdFx0XHR9LFxyXG5cdFx0fSksXHJcblx0XHR2aXRlQ29tcHJlc3Npb24oKSxcclxuXHRdLFxyXG5cdHJlc29sdmU6IHtcclxuXHRcdGFsaWFzOiB7XHJcblx0XHRcdCdAJzogZmlsZVVSTFRvUGF0aChuZXcgVVJMKCcuL3NyYycsIGltcG9ydC5tZXRhLnVybCkpLFxyXG5cdFx0fSxcclxuXHR9LFxyXG59KVxyXG4iXSwKICAibWFwcGluZ3MiOiAiO0FBQXNWLFNBQVMsZUFBZSxXQUFXO0FBRXpYLE9BQU8sU0FBUztBQUNoQixPQUFPLGtCQUFrQjtBQUN6QixPQUFPLGNBQWM7QUFDckIsU0FBUyxvQkFBb0I7QUFDN0IsT0FBTyxxQkFBcUI7QUFDNUIsU0FBUyxlQUFlO0FBUGlNLElBQU0sMkNBQTJDO0FBUzFRLElBQU8sc0JBQVEsYUFBYTtBQUFBLEVBQzNCLFFBQVE7QUFBQSxJQUNQLEtBQUs7QUFBQSxNQUNKLE1BQU07QUFBQSxNQUNOLFVBQVU7QUFBQSxNQUNWLE1BQU07QUFBQSxJQUNQO0FBQUEsRUFDRDtBQUFBLEVBQ0EsS0FBSztBQUFBLElBQ0osU0FBUztBQUFBLE1BQ1IsU0FBUyxDQUFDLFNBQVMsR0FBRyxhQUFhLENBQUM7QUFBQSxJQUNyQztBQUFBLElBQ0EscUJBQXFCO0FBQUEsTUFDcEIsTUFBTTtBQUFBLFFBQ0wsS0FBSztBQUFBLE1BQ047QUFBQSxJQUNEO0FBQUEsRUFDRDtBQUFBLEVBQ0EsU0FBUztBQUFBLElBQ1IsSUFBSTtBQUFBLElBQ0osUUFBUTtBQUFBLE1BQ1AsY0FBYztBQUFBLE1BQ2QsZ0JBQWdCO0FBQUEsTUFDaEIsU0FBUztBQUFBLFFBQ1IsZ0JBQWdCO0FBQUEsVUFDZjtBQUFBLFlBQ0MsWUFBWSxDQUFDLEVBQUUsUUFBUSxNQUN0QixRQUFRLGdCQUFnQixXQUN4QixRQUFRLGdCQUFnQixZQUN4QixRQUFRLGdCQUFnQjtBQUFBLFlBQ3pCLFNBQVM7QUFBQSxZQUNULFNBQVM7QUFBQSxjQUNSLFdBQVc7QUFBQSxjQUNYLFlBQVk7QUFBQSxnQkFDWCxZQUFZO0FBQUEsZ0JBQ1osZUFBZSxLQUFLLEtBQUssS0FBSztBQUFBO0FBQUEsY0FDL0I7QUFBQSxZQUNEO0FBQUEsVUFDRDtBQUFBLFVBQ0E7QUFBQSxZQUNDLFlBQVksQ0FBQyxFQUFFLFFBQVEsTUFBTSxRQUFRLGdCQUFnQjtBQUFBLFlBQ3JELFNBQVM7QUFBQSxZQUNULFNBQVM7QUFBQSxjQUNSLFdBQVc7QUFBQSxjQUNYLFlBQVk7QUFBQSxnQkFDWCxZQUFZO0FBQUEsZ0JBQ1osZUFBZSxLQUFLLEtBQUssS0FBSztBQUFBO0FBQUEsY0FDL0I7QUFBQSxZQUNEO0FBQUEsVUFDRDtBQUFBLFFBQ0Q7QUFBQSxNQUNEO0FBQUEsTUFDQSxZQUFZO0FBQUEsUUFDWCxTQUFTO0FBQUEsTUFDVjtBQUFBLElBQ0QsQ0FBQztBQUFBLElBQ0QsZ0JBQWdCO0FBQUEsRUFDakI7QUFBQSxFQUNBLFNBQVM7QUFBQSxJQUNSLE9BQU87QUFBQSxNQUNOLEtBQUssY0FBYyxJQUFJLElBQUksU0FBUyx3Q0FBZSxDQUFDO0FBQUEsSUFDckQ7QUFBQSxFQUNEO0FBQ0QsQ0FBQzsiLAogICJuYW1lcyI6IFtdCn0K
