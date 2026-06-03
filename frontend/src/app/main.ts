/**
 * main.ts — 藏叶前端入口文件
 *
 * 初始化 Vue 3 应用并挂载到 DOM：
 *   1. 创建 Vue 应用实例（App.vue 为根组件）
 *   2. 注册 Vue Router（路由管理）
 *   3. 加载全局样式（global.css）
 *   4. 挂载到 index.html 中的 #app 元素
 */

import { createApp } from 'vue'
import App from './App.vue'       // 根组件
import router from '@/router'     // 路由配置
import '@/styles/global.css'      // 全局样式（CSS Reset + 基础样式）

// 创建应用 → 注册路由 → 挂载到 DOM
createApp(App).use(router).mount('#app')