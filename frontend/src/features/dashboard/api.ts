/**
 * features/dashboard/api.ts — 仪表盘 API 调用
 *
 * 封装仪表盘相关的后端 API 请求，返回类型安全的 Promise。
 * 页面组件通过此模块获取数据，无需直接操作 HTTP 请求。
 */

import type { DashboardStats } from './types'
import { get } from '@/api/request'

/**
 * 获取仪表盘统计数据。
 *
 * 调用 GET /api/dashboard/stats 获取文件数、集合数、标签数和存储空间等信息。
 * 后端由 DashboardHandler.Stats 处理。
 *
 * @returns Promise<DashboardStats> - 仪表盘统计数据
 */
export function fetchDashboardStats(): Promise<DashboardStats> {
  return get<DashboardStats>('/api/dashboard/stats')
}