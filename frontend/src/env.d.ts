/**
 * env.d.ts — 环境类型声明
 *
 * 为 Vite 和 Vue 单文件组件提供 TypeScript 类型支持：
 *   - vite/client：提供 import.meta.env 等 Vite 特有 API 的类型
 *   - *.vue 模块声明：让 TypeScript 能识别 .vue 文件导入
 */

/// <reference types="vite/client" />

// 声明 .vue 单文件组件的模块类型
// 使 TypeScript 编译器能正确处理 import Xxx from './Xxx.vue' 语句
declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<Record<string, unknown>, Record<string, unknown>, unknown>
  export default component
}