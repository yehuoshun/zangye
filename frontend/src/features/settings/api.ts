import { request } from '@/api/request'
import type { Settings } from './types'

// 获取所有设置
export function getSettings(): Promise<Settings> {
  return request<Settings>('/settings')
}

// 更新设置
export function updateSettings(settings: Settings): Promise<void> {
  return request<void>('/settings', { method: 'PUT', body: settings })
}
