// 文件夹 API 类型定义
export interface FolderItem {
  id: string
  name: string
  parent_id: string | null
  sort_order: number
  created_at: string
  updated_at: string
  children?: FolderItem[]
  file_count?: number
}

export interface FolderStats {
  total_folders: number
  total_files: number
  total_size: number
}
