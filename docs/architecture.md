# 架构设计

## 系统架构总览

```
┌──────────────────────────────────────────────────┐
│                    浏览器                         │
│              http://127.0.0.1:27138               │
└──────────────────────┬───────────────────────────┘
                       │
                       ▼
┌──────────────────────────────────────────────────┐
│              Go HTTP Server (net/http)            │
│                                                    │
│  ┌──────────────┐  ┌───────────────────────────┐ │
│  │  API 路由     │  │  SPA 回退处理器           │ │
│  │  /api/*      │  │  (所有非 API 请求)         │ │
│  └──────┬───────┘  └───────────┬───────────────┘ │
│         │                       │                  │
│         ▼                       ▼                  │
│  ┌──────────────┐  ┌───────────────────────────┐ │
│  │  Handler 层   │  │  嵌入的前端 SPA           │ │
│  │  (dashboard, │  │  (frontend/dist/)          │ │
│  │   files,     │  │                            │ │
│  │   folders,   │  └───────────────────────────┘ │
│  │   tags,      │                                 │
│  │   file_tags, │                                 │
│  │   settings)  │                                 │
│  └──────┬───────┘                                 │
│         │                                          │
└─────────┼──────────────────────────────────────────┘
          │
          ▼
┌──────────────────────────────────────────────────┐
│              MySQL 8.0 (zang_ye)                  │
│  ┌──────────┐ ┌──────┐ ┌───────┐ ┌──────────┐  │
│  │ folders  │ │files │ │ tags  │ │ settings │  │
│  └──────────┘ └──────┘ └───┬───┘ └──────────┘  │
│                             │                    │
│                       ┌─────┴─────┐              │
│                       │ file_tags │              │
│                       └───────────┘              │
└──────────────────────────────────────────────────┘
```

## 分层架构

### 1. 入口层 (`main.go`)

- 初始化数据库连接
- 注册 HTTP 路由（API + SPA）
- 配置 HTTP 服务器参数
- 处理优雅关闭

### 2. 数据库层 (`internal/db/`)

- **`mysql.go`**：MySQL 连接管理
  - 连接池配置（最大 10 连接，空闲 5 连接）
  - 自动执行建表语句（幂等）
  - 连接健康检查
- **`schema.sql`**：表结构定义
  - 5 张表：folders, files, tags, file_tags, settings
  - 使用 InnoDB 引擎，支持外键和级联删除

### 3. Handler 层 (`internal/handler/`)

- **`dashboard.go`**：仪表盘 API — `GET /api/dashboard/stats`
- **`settings.go`**：设置 API — `GET/PUT /api/settings`
- **`folders.go`**：文件夹 CRUD + 递归统计（含分类）
- **`files.go`**：文件 CRUD + folder_id 过滤
- **`tags.go`**：标签 CRUD
- **`file_tags.go`**：文件-标签关联（获取/设置/添加/移除）
- **`response.go`**：公共工具（writeJSON, writeError, scanRow, closeRows, buildInQuery）

### 4. 前端层 (`frontend/`)

- Vue 3 + TypeScript + Vite 构建
- 编译后通过 `//go:embed` 嵌入到 Go 二进制
- 开发时通过 Vite proxy 代理 API 到后端

## 数据流

### API 请求流

```
浏览器 → GET /api/dashboard/stats
    → http.ServeMux 路由匹配
    → DashboardHandler.Stats()
    → MySQL 查询
    → JSON 编码响应
    → 浏览器解析 JSON
    → 组件更新视图
```

### SPA 请求流

```
浏览器 → GET /files
    → http.ServeMux 匹配 "/" 路由
    → spaHandler.ServeHTTP()
    → 尝试打开 /files 文件 → 不存在
    → 回退到 index.html
    → 浏览器加载 Vue SPA
    → Vue Router 解析 /files → 渲染 FilesPage
```

## 设计决策

| 决策 | 理由 |
|------|------|
| 使用标准库 net/http | 无外部依赖，简单可靠，适合小型项目 |
| 前端嵌入后端 | 单文件部署，无需 Nginx 反向代理 |
| SPA 回退到 index.html | 支持 Vue Router history 模式 |
| 自动建表 | 零配置启动，无需手动执行 SQL |
| 连接池配置 | 控制数据库资源，防止连接泄漏 |
| 优雅关闭 | 等待当前请求处理完毕再退出 |
| 动态 SQL 更新 | 只更新传入的字段，避免覆盖未传字段 |
| 文件标签事务 | 全量设置标签使用事务（删+插），保证原子性 |
