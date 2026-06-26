# 数据库设计

## 表结构

### folders (文件夹表)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | VARCHAR(36) PK | UUID 主键 |
| name | VARCHAR(255) | 文件夹名称 |
| parent_id | VARCHAR(36) NULL | 父文件夹 ID |
| sort_order | INT | 排序序号 |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |

### files (文件表)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | VARCHAR(36) PK | UUID 主键 |
| folder_id | VARCHAR(36) NULL | 所属文件夹 ID |
| name | VARCHAR(255) | 文件名 |
| paths | TEXT NULL | 路径 JSON 数组 |
| file_type | VARCHAR(32) | 文件类型 |
| file_size | BIGINT | 文件大小 |
| description | TEXT NULL | 描述 |
| deleted_at | TIMESTAMP NULL | 软删除时间 |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |

### tags (标签表)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | VARCHAR(36) PK | UUID 主键 |
| name | VARCHAR(64) UNIQUE | 标签名称 |
| color | VARCHAR(7) | 十六进制颜色 |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |

### file_tags (文件-标签关联)

| 字段 | 类型 | 说明 |
|------|------|------|
| file_id | VARCHAR(36) PK | 文件 ID |
| tag_id | VARCHAR(36) PK | 标签 ID |
| created_at | TIMESTAMP | 创建时间 |

外键：file_id → files(id) ON DELETE CASCADE, tag_id → tags(id) ON DELETE CASCADE

### settings (设置表)

| 字段 | 类型 | 说明 |
|------|------|------|
| setting_key | VARCHAR(128) PK | 键名 |
| setting_value | TEXT | 值（JSON 格式） |
| updated_at | TIMESTAMP | 更新时间 |

## ER 图

```
folders ──1:N── files ──N:M── tags
                    │
                    └── file_tags (中间表)
```

## 索引

- folders: parent_id, sort_order
- files: folder_id, file_type, deleted_at, name
- tags: name (UNIQUE)
- file_tags: file_id (PK), tag_id (INDEX)
