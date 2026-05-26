const BASE = '/api'

async function get<T = any>(path: string): Promise<T> {
  const res = await fetch(BASE + path)
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

async function post<T = any>(path: string, body: any): Promise<T> {
  const res = await fetch(BASE + path, {
    method: 'POST', headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

async function put<T = any>(path: string, body: any): Promise<T> {
  const res = await fetch(BASE + path, {
    method: 'PUT', headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

async function del(path: string): Promise<void> {
  const res = await fetch(BASE + path, { method: 'DELETE' })
  if (!res.ok) throw new Error(await res.text())
}

export interface Collection {
  id: string; name: string; icon: string
  parent_id: string | null; sort_order: number
  created_at: string; updated_at: string
}

export interface BrowseItem {
  type: 'file' | 'virtual' | 'link'
  name: string; path: string; size: number
  mod_time: string; mime_type: string; source: string
  display_name: string | null; is_dir: boolean
}

export interface OverviewStats {
  total_files: number; total_size: number
  categories: { name: string; label: string; count: number; size: number; percent: number }[]
  collection_count: number
}

export interface Tag { id: string; name: string; color: string }

export interface PrefixConfig {
  prefix: string; type: string
  map_path: string | null; url_template: string | null
}

export const collections = {
  listRoot: () => get<Collection[]>('/collections'),
  get: (id: string) => get<Collection>(`/collections/${id}`),
  children: (id: string) => get<Collection[]>(`/collections/${id}/children`),
  create: (data: Partial<Collection>) => post<Collection>('/collections', data),
  update: (id: string, data: Partial<Collection>) => put(`/collections/${id}`, data),
  delete: (id: string) => del(`/collections/${id}`),
  reorder: (order: { id: string; sort_order: number }[]) => put('/collections/reorder', order),
  browse: (id: string, params?: { search?: string; sort?: string }) => {
    const q = new URLSearchParams()
    if (params?.search) q.set('search', params.search)
    if (params?.sort) q.set('sort', params.sort)
    return get<BrowseItem[]>(`/collections/${id}/browse?${q}`)
  },
}

export const overview = { stats: () => get<OverviewStats>('/overview/stats') }

export const tags = {
  search: (q: string) => get<Tag[]>(`/tags/search?q=${encodeURIComponent(q)}`),
  create: (data: { name: string; color?: string }) => post<Tag>('/tags', data),
  delete: (id: string) => del(`/tags/${id}`),
  forFile: (fileId: string) => get<Tag[]>(`/files/${fileId}/tags`),
  updateFileTags: (fileId: string, tagIds: string[]) => put(`/files/${fileId}/tags`, tagIds),
  forCollection: (colId: string) => get<Tag[]>(`/collections/${colId}/tags`),
  updateCollectionTags: (colId: string, tagIds: string[]) => put(`/collections/${colId}/tags`, tagIds),
}

export const settings = {
  get: (key: string) => get<{ key: string; value: string }>(`/settings/${key}`),
  set: (key: string, value: string) => put(`/settings/${key}`, { value }),
  prefixes: () => get<PrefixConfig[]>('/prefixes'),
  updatePrefix: (prefix: string, data: Partial<PrefixConfig>) => put(`/prefixes/${prefix}`, data),
}

export function formatSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

export function formatTime(ts: string): string {
  if (!ts) return ''
  const d = new Date(ts)
  return d.toLocaleDateString('zh-CN') + ' ' + d.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}

export function getIcon(item: BrowseItem): string {
  if (item.type === 'link') return '🔗'
  if (item.mime_type?.startsWith('video/')) return '🎬'
  if (item.mime_type?.startsWith('audio/')) return '🎵'
  if (item.mime_type?.startsWith('image/')) return '🖼️'
  if (item.mime_type?.includes('pdf') || item.mime_type?.includes('document')) return '📄'
  if (item.mime_type?.startsWith('text/')) return '📝'
  return '📦'
}
