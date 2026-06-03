/**
 * vite.config.ts — Vite 构建配置
 *
 * 配置前端开发服务器和构建流程：
 *   - Vue 3 插件（SFC 编译）
 *   - 路径别名 @ → src/
 *   - 开发代理：/api 请求转发到后端 Go 服务
 *
 * 开发模式：npm run dev 启动 Vite 开发服务器（默认 :5173），
 * 通过 proxy 将 API 请求代理到 127.0.0.1:27138。
 *
 * 生产模式：npm run build 构建到 dist/，
 * Go 后端通过 embed 将 dist/ 打包进二进制文件。
 */

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'   // Vue 3 单文件组件插件
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],                    // 启用 Vue SFC 编译

  resolve: {
    alias: {
      // @ 别名指向 src 目录，简化导入路径
      // 示例：import X from '@/api/request' → import X from 'src/api/request'
      '@': resolve(__dirname, 'src'),
    },
  },

  server: {
    proxy: {
      // 开发时将 /api 请求代理到 Go 后端（避免跨域问题）
      '/api': 'http://127.0.0.1:27138',
    },
  },
})