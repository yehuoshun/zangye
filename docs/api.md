# API 文档

## 概述

藏叶后端提供 RESTful API，所有端点返回 JSON 格式数据。

- **Base URL**：`http://127.0.0.1:27138`
- **Content-Type**：`application/json`
- **字符编码**：UTF-8

---

## 端点列表

### 1. 健康检查

检查服务是否正常运行，包括数据库连接状态。

```
GET /api/health
```

**响应示例**（成功）：

```json
{
  "status": "ok"
}
```

**响应示例**（数据库不可用）：

```json
{
  "status": "db_error"
}
```

**HTTP 状态码**：

| 状态码 | 含义 |
|--------|------|
| 200 | 服务正常，数据库连接正常 |
| 503 | 服务运行中，但数据库连接失败 |

---

### 2. 仪表盘统计

获取文件管理系统的统计数据概览。

```
GET /api/dashboard/stats
```

**响应体**：

```typescript
{
  "file_count": number,        // 文件总数
  "collection_count": number,  // 集合总数
  "tag_count": number,         // 标签总数
  "storage_bytes": number,     // 存储空间总字节数
  "storage_display": string    // 存储空间的人类可读格式
}
```

**响应示例**：

```json
{
  "file_count": 128,
  "collection_count": 15,
  "tag_count": 42,
  "storage_bytes": 1073741824,
  "storage_display": "1.0 GB"
}
```

**HTTP 状态码**：

| 状态码 | 含义 |
|--------|------|
| 200 | 成功返回统计数据 |

**前端调用**：

```typescript
import { fetchDashboardStats } from '@/features/dashboard/api'

const stats = await fetchDashboardStats()
// stats: { file_count: 128, collection_count: 15, ... }
```

---

## 存储大小格式化规则

`storage_display` 字段使用 1024 进制（二进制前缀）计算：

| 范围 | 格式 | 示例 |
|------|------|------|
| 0 B | `0 B` | `0 B` |
| 1 KB ~ 1023 KB | `X.X KB` | `512.0 KB` |
| 1 MB ~ 1023 MB | `X.X MB` | `256.5 MB` |
| 1 GB ~ 1023 GB | `X.X GB` | `1.5 GB` |
| ≥ 1 TB | `X.X TB` | `2.0 TB` |

---

## 错误处理

当前 API 设计较为简单，后续版本计划添加：

- 统一的错误响应格式
- 请求参数验证
- 分页支持
- 认证/授权