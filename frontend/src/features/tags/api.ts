import { request } from '@/api/request'
import type { TagItem } from './types'

// 获取所有标签
export function getTags(): Promise<TagItem[]> {
  return request<TagItem[]>('/tags')
}

// 获取标签详情
export function getTag(id: string): Promise<TagItem> {
  return request<TagItem>(`/tags/${id}`)
}

// 创建标签
export function createTag(name: string): Promise<TagItem> {
  return request<TagItem>('/tags', { method: 'POST', body: { name } })
}

// 更新标签
export function updateTag(id: string, name: string, color?: string): Promise<TagItem> {
  return request<TagItem>(`/tags/${id}`, { method: 'PUT', body: { name, color } })
}

// 删除标签
export function deleteTag(id: string): Promise<void> {
  return request<void>(`/tags/${id}`, { method: 'DELETE' })
}
