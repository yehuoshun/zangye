import { request, requestList } from '@/api/request'
import type { FileItem, FileQuery } from './types'

// 查询文件列表
export function getFiles(query: FileQuery) {
  const params = new URLSearchParams()
  if (query.folder_id) params.set('folder_id', query.folder_id)
  if (query.keyword) params.set('keyword', query.keyword)
  if (query.type) params.set('type', query.type)
  if (query.tag_id) params.set('tag_id', query.tag_id)
  if (query.order_by) params.set('order_by', query.order_by)
  if (query.order_dir) params.set('order_dir', query.order_dir)
  if (query.page) params.set('page', String(query.page))
  if (query.page_size) params.set('page_size', String(query.page_size))
  if (query.trash) params.set('trash', 'true')

  const qs = params.toString()
  return requestList<FileItem>(`/files${qs ? '?' + qs : ''}`)
}

// 获取文件详情
export function getFile(id: string): Promise<FileItem> {
  return request<FileItem>(`/files/${id}`)
}

// 创建文件
export function createFile(data: Partial<FileItem> & { name: string }): Promise<FileItem> {
  return request<FileItem>('/files', { method: 'POST', body: data })
}

// 更新文件
export function updateFile(id: string, data: Partial<FileItem>): Promise<FileItem> {
  return request<FileItem>(`/files/${id}`, { method: 'PUT', body: data })
}

// 删除文件（软删除）
export function deleteFile(id: string): Promise<void> {
  return request<void>(`/files/${id}`, { method: 'DELETE' })
}

// 获取文件标签
export function getFileTags(id: string) {
  return request<any[]>(`/files/${id}/tags`)
}

// 设置文件标签
export function setFileTags(id: string, tagIds: string[]): Promise<void> {
  return request<void>(`/files/${id}/tags`, { method: 'PUT', body: { tag_ids: tagIds } })
}

// 添加文件标签
export function addFileTag(id: string, tagId: string): Promise<void> {
  return request<void>(`/files/${id}/tags`, { method: 'POST', body: { tag_id: tagId } })
}

// 移除文件标签
export function removeFileTag(id: string, tagId: string): Promise<void> {
  return request<void>(`/files/${id}/tags?tag_id=${tagId}`, { method: 'DELETE' })
}

// 获取回收站文件
export function getTrashFiles() {
  return requestList<FileItem>('/trash/files')
}

// 恢复文件
export function restoreFile(id: string): Promise<void> {
  return request<void>(`/trash/files/${id}/restore`, { method: 'POST' })
}

// 彻底删除文件
export function hardDeleteFile(id: string): Promise<void> {
  return request<void>(`/trash/files/${id}`, { method: 'DELETE' })
}
