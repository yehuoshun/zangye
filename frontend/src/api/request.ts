/**
 * request.ts — HTTP 请求封装
 *
 * 基于 fetch API 的轻量请求工具，提供类型安全的 GET 请求。
 * 所有 API 调用集中通过此模块，便于统一处理错误、添加认证头等。
 *
 * 使用方式：
 *   import { get } from '@/api/request'
 *   const data = await get<MyType>('/api/endpoint')
 */

// 生产环境中前端和后端在同一域名下，无需设置 base URL
const BASE = ''

/**
 * 发送 GET 请求并返回 JSON 解析后的数据。
 *
 * @param url  - API 路径（相对于 BASE）
 * @returns Promise<T> - 解析后的响应数据
 * @throws  Error - 当 HTTP 状态码非 2xx 时抛出
 *
 * @example
 *   const stats = await get<DashboardStats>('/api/dashboard/stats')
 */
export async function get<T>(url: string): Promise<T> {
  const res = await fetch(BASE + url)
  // 非 2xx 状态码视为错误
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}