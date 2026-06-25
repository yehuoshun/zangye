# 前端结构

## 技术栈

| 技术 | 版本 | 用途 |
|------|------|------|
| Vue | 3.x | 响应式 UI 框架 |
| TypeScript | 5.x | 类型安全 |
| Vite | 6.x | 构建工具和开发服务器 |
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
    │   └── App.vue         # 根组件（theme + layout 全局状态）
    ├── router/
    │   └── index.ts        # 路由配置
    ├── api/
    │   └── request.ts      # 通用 HTTP 请求封装
    ├── features/
    │   ├── dashboard/
    │   │   ├── api.ts      # 仪表盘 API 调用
    │   │   └── types.ts    # 仪表盘类型定义
    │   ├── files/
    │   │   ├── api.ts      # 文件管理 API（含标签关联）
    │   │   └── types.ts    # 文件类型定义
    │   ├── folders/
    │   │   ├── api.ts      # 文件夹管理 API
    │   │   └── types.ts    # 文件夹类型定义
    │   ├── tags/
    │   │   ├── api.ts      # 标签管理 API
    │   │   └── types.ts    # 标签类型定义
    │   └── settings/
    │       ├── api.ts      # 设置 API
    │       └── types.ts    # 设置类型定义
    ├── layouts/
    │   └── MainLayout.vue  # 主布局（侧边栏 + 顶部栏双布局）
    ├── pages/
    │   ├── Dashboard/
    │   │   └── DashboardPage.vue  # 仪表盘页面（7 张统计卡片）
    │   ├── Files/
    │   │   └── FilesPage.vue       # 文件管理页面（文件夹+文件浏览）
    │   ├── Tags/
    │   │   └── TagsPage.vue        # 标签管理页面
    │   └── Settings/
    │       └── SettingsPage.vue    # 设置页面
    └── styles/
        └── global.css      # 全局样式（CSS 变量主题）
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
│  根组件：theme + layout 状态        │
│  <router-view />                    │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│         MainLayout.vue              │
│  侧边栏（Logo + 导航）              │
│  顶部栏（面包屑 + 操作）             │
│  主内容区（<router-view />）        │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│       各页面组件                     │
│  Dashboard / Files / Tags / Settings│
│  调用 features/api.ts               │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│    features/*/api.ts                │
│  → fetch('/api/...')                │
└─────────────────────────────────────┘
```

## 路由设计

| 路径 | 组件 | 说明 |
|------|------|------|
| `/` | MainLayout > DashboardPage | 仪表盘（默认首页） |
| `/files` | MainLayout > FilesPage | 文件管理 |
| `/tags` | MainLayout > TagsPage | 标签管理 |
| `/settings` | MainLayout > SettingsPage | 系统设置 |

路由使用 `createWebHistory` 模式（HTML5 History），无 `#` 号。

## 数据流

### API 请求

```
组件 onMounted
  → features/*/api.ts: fetchXxx()
    → fetch('/api/xxx')
      → 后端 Go Handler
        → MySQL 查询
      ← JSON 响应
    ← 解析 JSON
  ← 类型安全的响应
  → 视图更新
```

### 文件标签数据流

```
FilesPage 加载
  → fetchFiles(folderId) → 获取文件列表
  → fetchFileTags(fileId) → 逐个获取文件标签
  → fileTagMap 缓存标签 → 列表/网格视图渲染 tag-chip
  → 右键菜单 → 管理标签 → 弹窗展示所有标签 checkbox
  → 保存 → setFileTags(fileId, tagIds) → 更新本地缓存
```

### 开发代理

在开发模式下（`npm run dev`），Vite 开发服务器运行在 `:5173`，通过 `vite.config.ts` 中的 proxy 配置将 `/api` 请求代理到 Go 后端 `127.0.0.1:27138`，避免跨域问题。

## 样式方案

- **全局样式**：`global.css` 提供 CSS 变量主题（`--bg-primary`, `--bg-secondary`, `--accent` 等），支持亮色/暗色切换
- **组件样式**：使用 Vue 的 `<style scoped>` 实现组件级样式隔离
- **CSS 变量**：所有颜色使用 `var(--bg-primary)` 等变量，不硬编码颜色值

## 构建流程

```bash
npm run build
# 1. TypeScript 类型检查（vue-tsc --noEmit）
# 2. Vite 打包（Rollup）
# 3. 输出到 frontend/dist/
# 4. Go build 时通过 embed 嵌入二进制
```