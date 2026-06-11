<!--
  MainLayout.vue — 藏叶主布局组件

  提供应用的整体布局结构：
    - 左侧固定侧边栏（220px）：Logo + 导航菜单 + 连接状态
    - 右侧主内容区：由 <router-view /> 渲染当前路由页面

  侧边栏功能：
    - Logo 展示（🦞 藏叶）
    - 导航菜单（使用 router-link 实现 SPA 导航）
    - 底部连接状态指示器（通过 /api/health 检测后端是否在线）

  样式采用暗色主题，与 global.css 保持一致。
-->

<template>
  <div class="app-layout">
    <!-- ===== 左侧边栏 ===== -->
    <aside class="sidebar">
      <!-- Logo 区域 -->
      <div class="logo">
        <span class="logo-icon">🦞</span>
        <span class="logo-text">藏叶</span>
      </div>

      <!-- 导航菜单 -->
      <nav class="nav">
        <!-- 仪表盘导航项 -->
        <!-- router-link-active 类名用于高亮当前激活的导航项 -->
        <router-link to="/" class="nav-item">
          <span class="nav-icon">📊</span>
          <span class="nav-label">仪表盘</span>
        </router-link>
        <router-link to="/settings" class="nav-item">
          <span class="nav-icon">⚙️</span>
          <span class="nav-label">设置</span>
        </router-link>
        <!-- TODO: 后续添加更多导航项（文件管理、标签管理等） -->
      </nav>

      <!-- 底部连接状态 -->
      <div class="sidebar-footer">
        <!-- 状态指示灯：绿色 = 在线，灰色 = 离线 -->
        <span class="status-dot" :class="{ online }" />
        <span class="status-text">{{ online ? '已连接' : '离线' }}</span>
      </div>
    </aside>

    <!-- ===== 右侧主内容区 ===== -->
    <main class="main">
      <!-- 路由出口：渲染当前激活的子路由页面 -->
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { fetchSettings } from '@/features/settings/api'

// 后端连接状态：true = 在线，false = 离线
const online = ref(false)

// 组件挂载后检测后端健康状态并加载主题
onMounted(async () => {
  // 加载主题设置
  try {
    const settings = await fetchSettings()
    document.documentElement.setAttribute('data-theme', settings.theme || 'dark')
  } catch { /* 使用默认深色主题 */ }

  // 检测后端健康状态
  try {
    // 调用健康检查 API 确认后端是否在线
    const res = await fetch('/api/health')
    online.value = res.ok
  } catch {
    // 请求失败（网络错误、后端未启动等）视为离线
    online.value = false
  }
})
</script>

<style scoped>
/* ===== 布局 ===== */

/* Flex 布局：侧边栏 + 主内容区并排 */
.app-layout { display: flex; height: 100vh; }

/* ===== 侧边栏 ===== */

.sidebar {
  width: 220px;
  background: var(--bg-secondary);
  display: flex;
  flex-direction: column;
  border-right: 1px solid var(--border-color);
  flex-shrink: 0;
}

/* Logo 区域 */
.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 20px 16px;
  border-bottom: 1px solid var(--border-color);
}
.logo-icon { font-size: 28px; }
.logo-text { font-size: 20px; font-weight: 700; color: var(--accent); }

/* 导航菜单 */
.nav {
  flex: 1;
  padding: 12px 8px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

/* 导航项 */
.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 8px;
  color: var(--text-secondary);
  text-decoration: none;
  transition: all .2s;
  font-size: 15px;
}
.nav-item:hover { background: var(--bg-tertiary); color: var(--text-primary); }
.nav-item.router-link-active { background: var(--bg-tertiary); color: var(--accent); }

.nav-icon { font-size: 18px; }

/* 底部状态栏 */
.sidebar-footer {
  padding: 16px;
  border-top: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: var(--text-muted);
}

/* 状态指示灯 */
.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;              /* 圆形 */
  background: #666;                /* 默认灰色（离线） */
}
.status-dot.online { background: #4ade80; } /* 绿色（在线） */

/* ===== 主内容区 ===== */

.main {
  flex: 1;                         /* 占据剩余空间 */
  overflow: auto;                  /* 内容溢出时滚动 */
  padding: 24px;
}
</style>