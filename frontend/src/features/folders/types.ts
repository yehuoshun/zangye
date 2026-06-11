/**
 * features/folders/types.ts — 文件夹类型定义
 */

/**
 * 文件夹 API 响应
 */
export interface FolderItem {
  id: string
  name: string
  icon: string
  parent_id: string | null
  description: string | null
  sort_order: number
  created_at: string
  updated_at: string
}

/**
 * 创建文件夹请求体
 */
export interface FolderCreateRequest {
  name: string
  icon?: string
  parent_id?: string | null
  description?: string | null
  sort_order?: number
}

/**
 * 更新文件夹请求体
 */
export interface FolderUpdateRequest {
  name?: string | null
  icon?: string | null
  parent_id?: string | null
  description?: string | null
  sort_order?: number | null
}