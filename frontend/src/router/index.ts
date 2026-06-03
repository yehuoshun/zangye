/**
 * router/index.ts — 藏叶前端路由配置
 *
 * 使用 Vue Router 管理页面导航：
 *   - 根路径 '/' 加载 MainLayout（侧边栏布局）
 *   - 默认子路由指向 Dashboard 页面（懒加载）
 *   - 所有未匹配路径重定向到首页
 *
 * 路由模式：createWebHistory（HTML5 History 模式）
 * 生产环境需后端配合 SPA 回退（main.go 已实现）。
 */

import { createRouter, createWebHistory } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'

// 路由表定义
const routes = [
  {
    path: '/',                   // 根路径
    component: MainLayout,       // 使用侧边栏布局包裹
    children: [
      {
        path: '',                // 默认子路由（访问 '/' 时渲染）
        name: 'Dashboard',       // 路由名称，用于编程式导航
        // 懒加载：仅在访问时加载 Dashboard 页面组件，减小首屏体积
        component: () => import('@/pages/Dashboard/DashboardPage.vue'),
      },
    ],
  },
  {
    // 通配符路由：匹配所有未定义的路径
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    redirect: '/',               // 重定向到首页
  },
]

// 创建路由实例
export default createRouter({
  history: createWebHistory(),  // HTML5 History 模式（无 # 号）
  routes,
})