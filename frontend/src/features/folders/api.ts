import { request } from '@/api/request'
import type { FolderItem, FolderStats } from './types'

// 获取文件夹树
export function getFolderTree(): Promise<FolderItem[]> {
  return request<FolderItem[]>('/folders')
}

// 获取文件夹详情
export function getFolder(id: string): Promise<FolderItem> {
  return request<FolderItem>(`/folders/${id}`)
}

// 创建文件夹
export function createFolder(name: string, parentId?: string | null): Promise<FolderItem> {
  return request<FolderItem>('/folders', {
    method: 'POST',
    body: { name, parent_id: parentId || null },
  })
}

// 更新文件夹
export function updateFolder(id: string, name: string, parentId?: string | null): Promise<FolderItem> {
  return request<FolderItem>(`/folders/${id}`, {
    method: 'PUT',
    body: { name, parent_id: parentId || null },
  })
}

// 删除文件夹
export function deleteFolder(id: string): Promise<void> {
  return request<void>(`/folders/${id}`, { method: 'DELETE' })
}

// 获取文件夹统计
export function getFolderStats(id: string): Promise<FolderStats> {
  return request<FolderStats>(`/folders/${id}/stats`)
}
