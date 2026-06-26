// 文件 API 类型定义
export interface FileItem {
  id: string
  folder_id: string | null
  name: string
  paths: string | null
  file_type: string
  file_size: number
  description: string | null
  deleted_at: string | null
  created_at: string
  updated_at: string
  tags: TagItem[]
}

export interface TagItem {
  id: string
  name: string
  color: string
  created_at: string
  updated_at: string
}

export interface FileQuery {
  folder_id?: string
  keyword?: string
  type?: string
  tag_id?: string
  order_by?: string
  order_dir?: string
  page?: number
  page_size?: number
  trash?: boolean
}
