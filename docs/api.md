# API 文档

## 概述

藏叶后端提供 RESTful API，所有端点返回 JSON 格式数据。

- **Base URL**：`http://127.0.0.1:27138`
- **Content-Type**：`application/json`
- **字符编码**：UTF-8

---

## 端点速查

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/health` | 健康检查 |
| GET | `/api/dashboard/stats` | 仪表盘统计 |
| GET | `/api/settings` | 获取所有设置 |
| PUT | `/api/settings` | 批量更新设置 |
| GET | `/api/folders` | 文件夹列表（支持 `?parent_id=`） |
| GET | `/api/folders/{id}` | 文件夹详情 |
| GET | `/api/folders/{id}/stats` | 文件夹统计（递归，含分类） |
| POST | `/api/folders` | 创建文件夹 |
| PUT | `/api/folders/{id}` | 更新文件夹 |
| DELETE | `/api/folders/{id}` | 删除文件夹 |
| GET | `/api/files` | 文件列表（支持 `?folder_id=`） |
| GET | `/api/files/{id}` | 文件详情 |
| POST | `/api/files` | 创建文件 |
| PUT | `/api/files/{id}` | 更新文件 |
| DELETE | `/api/files/{id}` | 删除文件 |
| GET | `/api/files/{id}/tags` | 获取文件标签 |
| PUT | `/api/files/{id}/tags` | 全量设置文件标签 |
| POST | `/api/files/{id}/tags` | 添加单个标签 |
| DELETE | `/api/files/{id}/tags/{tagId}` | 移除单个标签 |
| GET | `/api/tags` | 标签列表 |
| POST | `/api/tags` | 创建标签 |
| PUT | `/api/tags/{id}` | 更新标签 |
| DELETE | `/api/tags/{id}` | 删除标签 |

---

## 1. 健康检查

```
GET /api/health
```

**响应示例**（成功）：

```json
{ "status": "ok" }
```

**响应示例**（数据库不可用）：

```json
{ "status": "db_error" }
```

| 状态码 | 含义 |
|--------|------|
| 200 | 服务正常 |
| 503 | 数据库连接失败 |

---

## 2. 仪表盘统计

```
GET /api/dashboard/stats
```

**响应体**：

```typescript
{
  folder_count: number      // 文件夹总数
  file_count: number        // 文件总数
  image_count: number       // 图片数
  video_count: number       // 视频数
  audio_count: number       // 音频数
  other_count: number       // 其他文件数
  storage_bytes: number     // 总字节数
  storage_display: string   // 人类可读格式
}
```

**响应示例**：

```json
{
  "folder_count": 15,
  "file_count": 128,
  "image_count": 45,
  "video_count": 12,
  "audio_count": 8,
  "other_count": 63,
  "storage_bytes": 1073741824,
  "storage_display": "1.0 GB"
}
```

---

## 3. 设置

### 获取所有设置

```
GET /api/settings
```

**响应**：`{ "key": "value", ... }` 键值对对象。

### 批量更新设置

```
PUT /api/settings
```

**请求体**：`{ "key1": "value1", "key2": "value2" }`

**响应**：`{ "status": "ok" }`

---

## 4. 文件夹管理

### 获取文件夹列表

```
GET /api/folders
GET /api/folders?parent_id=xxx
```

不传 `parent_id` 返回根目录文件夹（`parent_id IS NULL`），传入则返回指定父文件夹下的子文件夹。

### 获取文件夹详情

```
GET /api/folders/{id}
```

### 获取文件夹统计（递归）

```
GET /api/folders/{id}/stats
```

**响应体**：

```typescript
{
  folder_count: number   // 子文件夹数
  file_count: number     // 文件总数（递归）
  total_size: number     // 总字节数（递归）
  image_count: number    // 图片数
  video_count: number    // 视频数
  audio_count: number    // 音频数
  other_count: number    // 其他文件数
}
```

### 创建文件夹

```
POST /api/folders
```

**请求体**：

```typescript
{
  name: string          // 必填
  icon?: string         // 默认 📁
  parent_id?: string    // null = 根目录
  description?: string  // 描述
  sort_order?: number   // 默认 0
}
```

### 更新文件夹

```
PUT /api/folders/{id}
```

**请求体**：所有字段可选（`name`, `icon`, `parent_id`, `description`, `sort_order`），只更新传入的字段。

### 删除文件夹

```
DELETE /api/folders/{id}
```

---

## 5. 文件管理

### 获取文件列表

```
GET /api/files
GET /api/files?folder_id=xxx
```

不传 `folder_id` 返回所有文件，传入则按文件夹过滤。

### 获取文件详情

```
GET /api/files/{id}
```

### 创建文件

```
POST /api/files
```

**请求体**：

```typescript
{
  folder_id: string      // 必填
  path: string           // 必填，文件路径
  display_name?: string  // 显示名称
  file_size?: number     // 字节
  mime_type?: string     // MIME 类型
  file_mtime?: string    // 文件修改时间
  sort_order?: number    // 默认 0
}
```

### 更新文件

```
PUT /api/files/{id}
```

**请求体**：所有字段可选，只更新传入的字段。

### 删除文件

```
DELETE /api/files/{id}
```

---

## 6. 文件-标签关联

### 获取文件标签

```
GET /api/files/{id}/tags
```

**响应**：`TagItem[]` 数组。

### 全量设置文件标签

```
PUT /api/files/{id}/tags
```

**请求体**：`{ "tag_ids": ["tag-uuid-1", "tag-uuid-2"] }`

先清空旧关联，再建立新关联（事务保证原子性）。

### 添加单个标签

```
POST /api/files/{id}/tags
```

**请求体**：`{ "tag_id": "tag-uuid" }`

使用 `INSERT IGNORE`，重复添加不报错。

### 移除单个标签

```
DELETE /api/files/{id}/tags/{tagId}
```

---

## 7. 标签管理

### 获取标签列表

```
GET /api/tags
```

### 创建标签

```
POST /api/tags
```

**请求体**：

```typescript
{
  name: string   // 必填，全局唯一
  color?: string // 默认 "gray"
}
```

### 更新标签

```
PUT /api/tags/{id}
```

**请求体**：`name`, `color` 可选。

### 删除标签

```
DELETE /api/tags/{id}
```

级联删除 `file_tags` 中的关联记录。

---

## 存储大小格式化规则

`storage_display` 使用 1024 进制：

| 范围 | 格式 | 示例 |
|------|------|------|
| 0 B | `0 B` | `0 B` |
| 1 KB ~ 1023 KB | `X.X KB` | `512.0 KB` |
| 1 MB ~ 1023 MB | `X.X MB` | `256.5 MB` |
| 1 GB ~ 1023 GB | `X.X GB` | `1.5 GB` |
| ≥ 1 TB | `X.X TB` | `2.0 TB` |

---

## 错误处理

所有错误响应格式：

```json
{ "error": "错误描述信息" }
```

常见 HTTP 状态码：

| 状态码 | 含义 |
|--------|------|
| 200 | 成功 |
| 201 | 创建成功 |
| 400 | 请求参数无效 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |
