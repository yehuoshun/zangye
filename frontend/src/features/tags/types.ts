/**
 * features/tags/types.ts — 标签类型定义
 */

export interface TagItem {
  id: string
  name: string
  color: string
  created_at: string
}

export interface TagCreateRequest {
  name: string
  color?: string
}

export interface TagUpdateRequest {
  name?: string | null
  color?: string | null
}
