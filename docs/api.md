# API 文档

所有端点前缀：`/api`

## 健康检查

```
GET /api/health
```

响应：`{"status":"ok","service":"zangye"}`

## 仪表盘

```
GET /api/dashboard/stats
```

响应：
```json
{
  "folder_count": 5,
  "file_count": 100,
  "image_count": 30,
  "video_count": 20,
  "audio_count": 10,
  "other_count": 40,
  "total_size": 1048576,
  "size_text": "1.00 MB"
}
```

## 设置

```
GET  /api/settings
PUT  /api/settings
```

PUT 请求体：`{"key": "value", ...}`

## 文件夹

```
GET    /api/folders              # 获取文件夹树
POST   /api/folders              # 创建文件夹
GET    /api/folders/{id}         # 获取文件夹详情
PUT    /api/folders/{id}         # 更新文件夹
DELETE /api/folders/{id}         # 删除文件夹
GET    /api/folders/{id}/stats   # 文件夹统计
```

### 创建文件夹

```json
{"name": "照片", "parent_id": null}
```

### 文件夹统计

```json
{"total_folders": 3, "total_files": 50, "total_size": 52428800}
```

## 文件

```
GET    /api/files                    # 查询文件列表
POST   /api/files                    # 创建文件
GET    /api/files/{id}               # 获取文件详情
PUT    /api/files/{id}               # 更新文件
DELETE /api/files/{id}               # 删除文件（软删除）
GET    /api/files/{id}/content       # 文件内容（支持 Range）
GET    /api/files/{id}/tags          # 获取文件标签
PUT    /api/files/{id}/tags          # 设置文件标签
POST   /api/files/{id}/tags          # 添加文件标签
DELETE /api/files/{id}/tags          # 移除文件标签
```

### 查询参数

| 参数 | 说明 | 示例 |
|------|------|------|
| folder_id | 按文件夹筛选 | `?folder_id=xxx` |
| keyword | 关键字搜索 | `?keyword=照片` |
| type | 按类型筛选 | `?type=jpg` |
| tag_id | 按标签筛选 | `?tag_id=xxx` |
| order_by | 排序字段 | `name/file_type/file_size/created_at` |
| order_dir | 排序方向 | `asc/desc` |
| page | 页码 | `?page=1` |
| page_size | 每页数量 | `?page_size=50` |

### 创建文件

```json
{
  "folder_id": null,
  "name": "photo.jpg",
  "paths": "[\"Z:\\\\backup\\\\photo.jpg\", \"115:/backup/photo.jpg\"]",
  "file_size": 102400,
  "description": "旅游照片"
}
```

## 标签

```
GET    /api/tags           # 获取所有标签
POST   /api/tags           # 创建标签
GET    /api/tags/{id}      # 获取标签详情
PUT    /api/tags/{id}      # 更新标签
DELETE /api/tags/{id}      # 删除标签
```

## 回收站

```
GET    /api/trash/files                # 回收站文件列表
POST   /api/trash/files/{id}/restore   # 恢复文件
DELETE /api/trash/files/{id}           # 彻底删除
```

## 错误响应

所有错误统一格式：

```json
{"error": "错误描述"}
```

HTTP 状态码：
- 400: 请求参数错误
- 404: 资源不存在
- 405: 不支持的请求方法
- 500: 服务器内部错误
