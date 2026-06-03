# 前端结构

## 技术栈

| 技术 | 版本 | 用途 |
|------|------|------|
| Vue | 3.x | 响应式 UI 框架 |
| TypeScript | 5.x | 类型安全 |
| Vite | 5.x | 构建工具和开发服务器 |
| Vue Router | 4.x | 客户端路由 |

## 目录结构

```
frontend/
├── index.html              # HTML 入口（Vite 模板）
├── package.json            # 依赖和脚本
├── vite.config.ts          # Vite 配置（代理、别名）
├── tsconfig.json           # TypeScript 配置
└── src/
    ├── env.d.ts            # 环境类型声明（.vue 模块）
    ├── app/
    │   ├── main.ts         # 应用入口：创建 Vue 实例
    │   └── App.vue         # 根组件（仅 <router-view />）
    ├── router/
    │   └── index.ts        # 路由配置
    ├── api/
    │   └── request.ts      # 通用 HTTP 请求封装
    ├── features/
    │   └── dashboard/
    │       ├── api.ts      # 仪表盘 API 调用
    │       └── types.ts    # 仪表盘类型定义
    ├── layouts/
    │   └── MainLayout.vue  # 主布局（侧边栏 + 内容区）
    ├── pages/
    │   └── Dashboard/
    │       └── DashboardPage.vue  # 仪表盘页面
    └── styles/
        └── global.css      # 全局样式（CSS Reset + 基础）
```

## 模块架构

```
┌─────────────────────────────────────┐
│            main.ts                  │
│  创建 Vue 应用 → 注册 Router        │
│  → 加载全局样式 → 挂载到 DOM        │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│            App.vue                  │
│  根组件：<router-view />            │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│         MainLayout.vue              │
│  侧边栏（Logo + 导航 + 状态）       │
│  主内容区（<router-view />）        │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│       DashboardPage.vue             │
│  调用 fetchDashboardStats()         │
│  → 渲染统计卡片                     │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│    features/dashboard/api.ts        │
│  → api/request.ts                   │
│    → fetch('/api/dashboard/stats')  │
└─────────────────────────────────────┘
```

## 路由设计

| 路径 | 组件 | 说明 |
|------|------|------|
| `/` | MainLayout > DashboardPage | 仪表盘（默认首页） |
| `/:pathMatch(.*)*` | 重定向到 `/` | 404 兜底 |

路由使用 `createWebHistory` 模式（HTML5 History），无 `#` 号。

## 数据流

### API 请求

```
组件 onMounted
  → features/dashboard/api.ts: fetchDashboardStats()
    → api/request.ts: get<DashboardStats>('/api/dashboard/stats')
      → fetch('/api/dashboard/stats')
        → 后端 Go Handler
          → MySQL 查询
        ← JSON 响应
      ← 解析 JSON
    ← 类型安全的 DashboardStats
  ← stats.value = {...}
  → 视图更新（卡片渲染）
```

### 开发代理

在开发模式下（`npm run dev`），Vite 开发服务器运行在 `:5173`，通过 `vite.config.ts` 中的 proxy 配置将 `/api` 请求代理到 Go 后端 `127.0.0.1:27138`，避免跨域问题。

## 样式方案

- **全局样式**：`global.css` 提供 CSS Reset 和暗色主题基础
- **组件样式**：使用 Vue 的 `<style scoped>` 实现组件级样式隔离
- **主题色**：
  - 背景：`#1a1a2e`（深蓝）
  - 卡片：`#16213e`（稍浅蓝）
  - 强调：`#e94560`（红色）
  - 文字：`#e0e0e0`（浅灰）

## 构建流程

```bash
npm run build
# 1. TypeScript 类型检查
# 2. Vite 打包（Rollup）
# 3. 输出到 frontend/dist/
# 4. Go build 时通过 embed 嵌入二进制
```