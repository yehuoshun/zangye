# 前端文档

## 技术栈

- **框架**: Vue 3 (Composition API + `<script setup>`)
- **语言**: TypeScript
- **构建**: Vite 6
- **路由**: Vue Router 4 (history 模式)
- **样式**: 纯 CSS 变量（暗色 Material 风格）
- **无第三方 UI 库**

## 项目结构

```
frontend/src/
├── app/               # 入口
│   ├── main.ts        # 应用启动
│   └── App.vue        # 根组件
├── router/            # 路由配置
├── api/               # HTTP 请求封装
├── features/          # 功能模块（API + 类型定义）
│   ├── dashboard/
│   ├── files/
│   ├── folders/
│   ├── tags/
│   └── settings/
├── pages/             # 页面组件
│   ├── Dashboard/
│   ├── Files/
│   ├── Tags/
│   └── Settings/
├── layouts/           # 布局组件
│   ├── MainLayout.vue
│   └── FolderTreeNode.vue
└── styles/            # 全局样式
    └── global.css
```

## 路由

| 路径 | 页面 | 说明 |
|------|------|------|
| `/` | Dashboard | 仪表盘 |
| `/files` | Files | 文件管理 |
| `/tags` | Tags | 标签管理 |
| `/trash` | Files (trash) | 回收站 |
| `/settings` | Settings | 设置 |

## 功能模块

### 文件管理 (FilesPage)
- 网格/列表视图切换
- 表头点击排序
- 关键字搜索
- 文件夹筛选
- 标签筛选
- 分页
- 预览弹窗（图片/视频/音频/文本）
- 编辑/删除
- 右键菜单
- 标签管理

### 标签管理 (TagsPage)
- 标签卡片展示
- 创建/编辑/删除
- 颜色选择

### 仪表盘 (DashboardPage)
- 7 张统计卡片
- 文件夹数/文件数/图片/视频/音频/其他/存储空间

### 设置 (SettingsPage)
- 默认视图
- 每页数量
- 打开方式配置

## 样式系统

使用 CSS 变量实现暗色 Material 风格：

```css
--bg-primary: #121212;
--bg-secondary: #1e1e1e;
--accent-primary: #7c4dff;
--text-primary: #e0e0e0;
```

## 开发

```bash
cd frontend
npm install
npm run dev     # 开发模式
npm run build   # 构建
```
