/**
 * features/files/api.ts — 文件管理 API 调用
 */

import type { FileItem, FileCreateRequest, FileUpdateRequest, FilePreview } from './types'

const BASE = '/api/files'

/**
 * 预览文件信息（不写入数据库）
 */
export async function previewFile(path: string): Promise<FilePreview> {
  const res = await fetch(`${BASE}/preview`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ path }),
  })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

/**
 * 获取文件列表
 */
export async function fetchFiles(): Promise<FileItem[]> {
  const res = await fetch(BASE)
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