# 数据库设计

## 概述

藏叶使用 MySQL 8.0 作为数据存储，数据库名为 `zang_ye`，字符集为 `utf8mb4`，所有表使用 InnoDB 引擎以支持事务和外键。

## 表结构

### 1. collections（集合表）

文件集合，支持树形结构（通过 `parent_id` 自引用）。

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | VARCHAR(36) PK | 集合唯一标识（UUID） |
| `name` | VARCHAR(255) NOT NULL | 集合名称 |
| `icon` | VARCHAR(10) DEFAULT '📁' | 集合图标（emoji） |
| `parent_id` | VARCHAR(36) DEFAULT NULL | 父集合 ID（NULL = 根集合） |
| `sort_order` | INT DEFAULT 0 | 排序序号 |
| `created_at` | TIMESTAMP | 创建时间 |
| `updated_at` | TIMESTAMP | 更新时间（自动更新） |

### 2. files（文件表）

记录文件的元信息，每个文件必须属于一个集合。

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | VARCHAR(36) PK | 文件唯一标识（UUID） |
| `collection_id` | VARCHAR(36) FK NOT NULL | 所属集合 ID |
| `path` | VARCHAR(1024) NOT NULL | 文件在本地文件系统中的路径 |
| `display_name` | VARCHAR(512) DEFAULT NULL | 显示名称 |
| `file_size` | BIGINT DEFAULT 0 | 文件大小（字节） |
| `mime_type` | VARCHAR(255) DEFAULT NULL | MIME 类型 |
| `sort_order` | INT DEFAULT 0 | 排序序号 |
| `created_at` | TIMESTAMP | 创建时间 |

**外键**：`collection_id → collections(id) ON DELETE CASCADE`

### 3. tags（标签表）

标签用于分类和检索，名称全局唯一。

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | VARCHAR(36) PK | 标签唯一标识（UUID） |
| `name` | VARCHAR(128) UNIQUE NOT NULL | 标签名称 |
| `color` | VARCHAR(32) DEFAULT 'gray' | 标签颜色 |
| `created_at` | TIMESTAMP | 创建时间 |

### 4. file_tags（文件-标签关联表）

多对多关联表，连接 files 和 tags。

| 字段 | 类型 | 说明 |
|------|------|------|
| `file_id` | VARCHAR(36) PK, FK | 文件 ID |
| `tag_id` | VARCHAR(36) PK, FK | 标签 ID |

**联合主键**：`(file_id, tag_id)` — 防止重复关联

**外键**：
- `file_id → files(id) ON DELETE CASCADE`
- `tag_id → tags(id) ON DELETE CASCADE`

### 5. settings（设置表）

系统级配置项，键值对存储。

| 字段 | 类型 | 说明 |
|------|------|------|
| `key` | VARCHAR(128) PK | 设置键名 |
| `value` | TEXT NOT NULL | 设置值 |

**默认数据**：

| key | value | 说明 |
|-----|-------|------|
| `version` | `1` | 数据库版本号 |
| `theme` | `dark` | 默认主题（dark/light） |
| `layout` | `sidebar` | 默认布局模式 |

## ER 图

```
┌──────────────┐       ┌──────────────┐
│  collections │       │     tags     │
│──────────────│       │──────────────│
│ id (PK)      │       │ id (PK)      │
│ name         │       │ name (UNIQUE)│
│ icon         │       │ color        │
│ parent_id ───┼─┐     │ created_at   │
│ sort_order   │ │     └──────┬───────┘
│ created_at   │ │            │
│ updated_at   │ │            │
└──────┬───────┘ │            │
       │         │            │
       │ 1:N     │ 自引用     │
       ▼         │            │
┌──────────────┐ │            │
│    files     │ │            │
│──────────────│ │            │
│ id (PK)      │ │            │
│ collection_id│◄┘            │
│ path         │              │
│ display_name │              │
│ file_size    │              │
│ mime_type    │              │
│ sort_order   │              │
│ created_at   │              │
└──────┬───────┘              │
       │                      │
       │ N:M                  │
       ▼                      ▼
┌──────────────────────────────────────┐
│             file_tags                │
│──────────────────────────────────────│
│ file_id (PK, FK → files.id)          │
│ tag_id  (PK, FK → tags.id)           │
└──────────────────────────────────────┘

┌──────────────┐
│   settings   │
│──────────────│
│ key (PK)     │
│ value        │
└──────────────┘
```

## 级联删除策略

| 操作 | 影响 |
|------|------|
| 删除集合 | 级联删除该集合下所有文件 |
| 删除文件 | 级联删除该文件的所有标签关联 |
| 删除标签 | 级联删除该标签的所有文件关联 |

## 初始化

数据库初始化在程序启动时自动完成（`internal/db/mysql.go` 的 `New()` 函数）。所有 `CREATE TABLE` 使用 `IF NOT EXISTS`，确保幂等执行。默认设置通过 `INSERT IGNORE` 插入，不会覆盖已有数据。