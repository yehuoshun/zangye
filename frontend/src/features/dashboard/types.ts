/**
 * features/dashboard/types.ts — 仪表盘类型定义
 *
 * 定义仪表盘模块的 TypeScript 类型，与后端 Go 结构体保持一致。
 * JSON 字段名使用 snake_case，与后端 json tag 对齐。
 */

/**
 * 仪表盘统计数据
 *
 * 对应后端 handler.DashboardStats 结构体。
 * 字段名使用 snake_case 以匹配后端 JSON 序列化。
 */
export interface DashboardStats {
  file_count: number        // 文件总数
  collection_count: number  // 集合总数
  tag_count: number         // 标签总数
  storage_bytes: number     // 存储空间总字节数
  storage_display: string   // 存储空间的人类可读格式（如 "1.5 GB"）
}