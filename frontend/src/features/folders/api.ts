/**
 * features/folders/api.ts — 文件夹管理 API 调用
 */

import type { FolderItem, FolderCreateRequest, FolderUpdateRequest } from './types'

/**
 * 文件夹统计信息
 */
export interface FolderStats {
  folder_count: number
  file_count: number
  total_size: number
  image_count: number
  video_count: number
  audio_count: number
  other_count: number
}

const BASE = '/api/folders'

/**
 * 获取文件夹列表（可选 parent_id 查询子文件夹）
 */
export async function fetchFolders(parentId?: string): Promise<FolderItem[]> {
  const url = parentId ? `${BASE}?parent_id=${parentId}` : BASE
  const res = await fetch(url)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

/**
 * 获取单个文件夹详情
 */
export async function fetchFolder(id: string): Promise<FolderItem> {
  const res = await fetch(`${BASE}/${id}`)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

/**
 * 创建文件夹
 */
export async function createFolder(data: FolderCreateRequest): Promise<FolderItem> {
  const res = await fetch(BASE, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

/**
 * 更新文件夹
 */
export async function updateFolder(id: string, data: FolderUpdateRequest): Promise<FolderItem> {
  const res = await fetch(`${BASE}/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

/**
 * 删除文件夹
 */
export async function deleteFolder(id: string): Promise<void> {
  const res = await fetch(`${BASE}/${id}`, { method: 'DELETE' })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
}

/**
 * 获取文件夹统计（递归子文件夹）
 */
export async function fetchFolderStats(id: string): Promise<FolderStats> {
  const res = await fetch(`${BASE}/${id}/stats`)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}