/**
 * features/files/api.ts — 文件管理 API 调用
 */

import type { FileItem, FileCreateRequest, FileUpdateRequest } from './types'

const BASE = '/api/files'

/**
 * 获取文件列表
 * @param folderId 可选，按文件夹过滤
 */
export async function fetchFiles(folderId?: string): Promise<FileItem[]> {
  const url = folderId ? `${BASE}?folder_id=${encodeURIComponent(folderId)}` : BASE
  const res = await fetch(url)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

/**
 * 获取单个文件详情
 */
export async function fetchFile(id: string): Promise<FileItem> {
  const res = await fetch(`${BASE}/${id}`)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

/**
 * 创建文件
 */
export async function createFile(data: FileCreateRequest): Promise<FileItem> {
  const res = await fetch(BASE, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

/**
 * 更新文件
 */
export async function updateFile(id: string, data: FileUpdateRequest): Promise<FileItem> {
  const res = await fetch(`${BASE}/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

/**
 * 删除文件
 */
export async function deleteFile(id: string): Promise<void> {
  const res = await fetch(`${BASE}/${id}`, { method: 'DELETE' })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
}

// ===== 文件-标签关联 =====

import type { TagItem } from '@/features/tags/types'

/**
 * 获取文件的标签列表
 */
export async function fetchFileTags(fileId: string): Promise<TagItem[]> {
  const res = await fetch(`${BASE}/${fileId}/tags`)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

/**
 * 全量设置文件标签
 */
export async function setFileTags(fileId: string, tagIds: string[]): Promise<void> {
  const res = await fetch(`${BASE}/${fileId}/tags`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ tag_ids: tagIds }),
  })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
}

/**
 * 给文件添加单个标签
 */
export async function addFileTag(fileId: string, tagId: string): Promise<void> {
  const res = await fetch(`${BASE}/${fileId}/tags`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ tag_id: tagId }),
  })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
}

/**
 * 移除文件的单个标签
 */
export async function removeFileTag(fileId: string, tagId: string): Promise<void> {
  const res = await fetch(`${BASE}/${fileId}/tags/${tagId}`, { method: 'DELETE' })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
}