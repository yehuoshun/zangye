import type { DashboardStats } from './types'
import { get } from '@/api/request'

export function fetchDashboardStats(): Promise<DashboardStats> {
  return get<DashboardStats>('/api/dashboard/stats')
}