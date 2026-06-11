/**
 * features/files/types.ts — 文件管理类型定义
 */

/**
 * 文件 API 响应
 * 对应后端 handler.FileResponse
 */
export interface FileItem {
  id: string
  folder_id: string
  path: string
  display_name: string | null
  file_size: number
  mime_type: string | null
  file_mtime: string | null
  sort_order: number
  created_at: string
}

/**
 * 创建文件请求体
 */
export interface FileCreateRequest {
  folder_id: string
  path: string
  display_name?: string | null
  file_size?: number
  mime_type?: string | null
  file_mtime?: string | null
  sort_order?: number
}

/**
 * 更新文件请求体（所有字段可选）
 */
export interface FileUpdateRequest {
  folder_id?: string | null
  path?: string | null
  display_name?: string | null
  file_size?: number | null
  mime_type?: string | null
  file_mtime?: string | null
  sort_order?: number | null
}