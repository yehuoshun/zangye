# 藏叶 (ZangYe) 虚拟文件管理器

> 藏叶于林，藏名于器。

## 简介

藏叶是一个虚拟文件管理器，用于管理用户在互联网上的备份文件清单。

**核心功能：**
- 轻量级目录树 + 文件记录
- 标签系统（多对多）
- 文件预览（图片/视频/音频/文本）
- 打开方式配置
- 回收站（软删除）
- 仪表盘统计

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go 1.22+ (net/http 标准库) |
| 数据库 | MySQL 8.0 |
| 前端 | Vue 3 + TypeScript + Vite 6 |
| 路由 | Vue Router 4 (history 模式) |
| 部署 | 单 exe (go:embed) |

## 快速开始

### 前置条件

- Go 1.22+
- Node.js 18+
- MySQL 8.0

### 1. 初始化数据库

```bash
mysql -u root -p < internal/db/schema.sql
```

### 2. 启动后端

```bash
# 设置数据库连接（可选，默认 root:root@tcp(127.0.0.1:3306)/zangye）
export ZANGYE_DSN="user:password@tcp(127.0.0.1:3306)/zangye?charset=utf8mb4&parseTime=True&loc=Local"

# 启动服务
go run main.go
```

服务启动在 `http://127.0.0.1:27138`

### 3. 启动前端（开发模式）

```bash
cd frontend
npm install
npm run dev
```

前端开发服务器启动在 `http://127.0.0.1:5173`，自动代理 API 到后端。

### 4. 构建单 exe

```bash
cd frontend && npm install && npm run build && cd ..
go build -o zangye.exe -ldflags="-s -w" .
```

## 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `ZANGYE_ADDR` | HTTP 监听地址 | `127.0.0.1:27138` |
| `ZANGYE_DSN` | MySQL 数据源 | `root:root@tcp(127.0.0.1:3306)/zangye?charset=utf8mb4&parseTime=True&loc=Local` |

## 项目结构

```
zangye/
├── main.go              # 主入口
├── go.mod               # Go 模块定义
├── internal/
│   ├── config/          # 配置管理
│   ├── db/              # 数据库连接 + 建表 SQL
│   ├── model/           # 数据模型
│   ├── repository/      # 数据访问层 (DAO)
│   ├── service/         # 业务逻辑层
│   ├── handler/         # HTTP 处理器 (Controller)
│   ├── middleware/      # 中间件 (CORS)
│   └── util/            # 工具类
├── frontend/            # Vue 3 前端
│   ├── src/
│   │   ├── app/         # 入口
│   │   ├── router/      # 路由
│   │   ├── api/         # HTTP 请求封装
│   │   ├── features/    # 功能模块
│   │   ├── pages/       # 页面组件
│   │   ├── layouts/     # 布局组件
│   │   └── styles/      # 全局样式
│   └── dist/            # 构建输出 (go:embed)
├── docs/                # 文档
└── .github/workflows/   # CI
```
