import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
  // 构建输出到 Go 后端嵌入的目录
  build: {
    outDir: '../frontend/dist',
    emptyOutDir: true,
  },
  server: {
    port: 5173,
    // 代理 API 请求到 Go 后端
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:27138',
        changeOrigin: true,
      },
    },
  },
})
