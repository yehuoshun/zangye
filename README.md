# 🦞 藏叶 — 个人文件管理器

> 藏 = 收藏收纳，叶 = 文件如叶
> 把散落各处的文件收纳到一处，像叶子藏在树荫下。

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go 1.22 + net/http |
| 数据库 | MySQL 8.0 |
| 前端 | Vue 3 + Vite + TypeScript |
| 部署 | 单二进制（前端嵌入） |

## 快速开始

### 环境要求

- Go 1.22+
- MySQL 8.0+
- Node.js 18+（仅开发前端时）

### 初始化数据库

```sql
CREATE DATABASE IF NOT EXISTS zang_ye CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER IF NOT EXISTS 'zangye'@'localhost' IDENTIFIED BY 'zangye123';
GRANT ALL PRIVILEGES ON zang_ye.* TO 'zangye'@'localhost';
FLUSH PRIVILEGES;
```

### 编译运行

```bash
# 编译前端
cd frontend && npm install && npm run build && cd ..

# 编译运行
go build -o zangye .
./zangye
```

启动后访问 `http://127.0.0.1:27138`

### 开发模式

```bash
# 终端 1：启动后端
go run .

# 终端 2：启动前端开发服务器（自动代理 /api）
cd frontend && npm run dev
```

## 项目结构

```
zangye/
├── main.go                 # 入口，HTTP 服务 + SPA 回退
├── internal/
│   └── db/
│       ├── mysql.go        # MySQL 连接 + 自动建表
│       └── schema.sql      # 表结构定义
├── frontend/
│   ├── index.html
│   ├── vite.config.ts
│   └── src/
│       ├── main.ts         # Vue 入口 + Router
│       ├── App.vue         # 侧边栏布局
│       └── views/          # 页面组件
└── .gitignore
```

## 数据库配置

默认连接 `127.0.0.1:3306`，数据库 `zang_ye`，用户 `zangye`。

可通过环境变量覆盖：

```
# TODO: 后续支持
```

## 许可

MIT