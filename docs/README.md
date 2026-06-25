# 🦞 藏叶 — 项目文档

> 藏 = 收藏收纳，叶 = 文件如叶
> 把散落各处的文件收纳到一处，像叶子藏在树荫下。

藏叶是一个个人文件管理器，提供文件集合管理、标签系统和统计仪表盘。

## 技术栈

| 层级 | 技术 | 说明 |
|------|------|------|
| 后端 | Go 1.22 + net/http | 标准库 HTTP 服务，无第三方框架 |
| 数据库 | MySQL 8.0 | 关系型数据库，InnoDB 引擎 |
| 前端 | Vue 3 + Vite + TypeScript | 响应式 SPA 前端 |
| 部署 | 单二进制文件 | 前端构建产物通过 embed 嵌入 |

## 文档导航

| 文档 | 内容 |
|------|------|
| [架构设计](./architecture.md) | 系统架构、模块划分、数据流 |
| [API 文档](./api.md) | 后端 API 端点说明 |
| [数据库设计](./database.md) | 表结构、字段说明、关系图 |
| [前端结构](./frontend.md) | 前端目录结构、组件说明、路由设计 |
| [构建部署](./deploy.md) | 编译、打包、部署流程 |

## 快速开始

```bash
# 编译前端
cd frontend && npm install && npm run build && cd ..

# 编译运行
go build -o zangye .
./zangye
```

启动后访问 `http://127.0.0.1:27138`

## 项目结构

```
zangye/
├── main.go                 # 入口：HTTP 服务 + SPA 回退
├── go.mod / go.sum         # Go 模块依赖
├── internal/
│   ├── db/
│   │   ├── mysql.go        # MySQL 连接 + 自动建表
│   │   └── schema.sql      # 数据库表结构
│   └── handler/
│       ├── dashboard.go    # 仪表盘 API
│       ├── settings.go     # 设置 API
│       ├── folders.go      # 文件夹 CRUD + 递归统计
│       ├── files.go        # 文件 CRUD + folder_id 过滤
│       ├── tags.go         # 标签 CRUD
│       ├── file_tags.go    # 文件-标签关联
│       └── response.go     # 公共工具函数
├── frontend/
│   ├── package.json        # 前端依赖
│   ├── vite.config.ts      # Vite 构建配置
│   ├── tsconfig.json       # TypeScript 配置
│   └── src/
│       ├── app/
│       │   ├── main.ts     # 前端入口
│       │   └── App.vue     # 根组件
│       ├── router/
│       │   └── index.ts    # 路由配置
│       ├── api/
│       │   └── request.ts  # HTTP 请求封装
│       ├── features/
│       │   ├── dashboard/  # 仪表盘模块
│       │   ├── files/      # 文件模块（含标签关联）
│       │   ├── folders/    # 文件夹模块
│       │   ├── tags/       # 标签模块
│       │   └── settings/   # 设置模块
│       ├── layouts/
│       │   └── MainLayout.vue  # 主布局
│       ├── pages/
│       │   ├── Dashboard/  # 仪表盘页面
│       │   ├── Files/      # 文件管理页面
│       │   ├── Tags/       # 标签管理页面
│       │   └── Settings/   # 设置页面
│       └── styles/
│           └── global.css  # 全局样式
└── docs/                   # 项目文档
    ├── README.md           # 文档索引（本文件）
    ├── architecture.md     # 架构设计
    ├── api.md              # API 文档
    ├── database.md         # 数据库设计
    ├── frontend.md         # 前端结构
    └── deploy.md           # 构建部署
```

## 核心功能

- **文件夹管理**：树形结构组织文件，支持嵌套子文件夹，右键菜单操作
- **文件管理**：网格/列表双视图，面包屑导航，按文件夹过滤，多路径批量创建
- **标签系统**：为文件打标签/去标签，支持多对多关联，10 色选择器
- **文件夹详情**：BFS 递归统计文件数、总大小，以及图片/视频/音频/其他分类统计
- **仪表盘**：7 张统计卡片，存储空间支持 TB/GB/MB/KB/B 5 维度拆分展示