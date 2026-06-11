import type { Settings } from './types'
import { get } from '@/api/request'

export function fetchSettings(): Promise<Settings> {
  return get<Settings>('/api/settings')
}

export async function updateSettings(data: Partial<Settings>): Promise<void> {
  const res = await fetch('/api/settings', {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
}