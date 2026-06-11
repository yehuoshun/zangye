<!--
  MainLayout.vue — 藏叶主布局组件

  支持两种布局模式：
    - sidebar（侧边栏）：左侧固定 220px 导航 + 右侧内容区
    - topbar（顶部导航）：顶部 56px 导航栏 + 下方内容区

  布局状态从 App.vue inject，修改后即时响应无需刷新。
-->

<template>
  <!-- ===== 侧边栏布局 ===== -->
  <div v-if="layout === 'sidebar'" class="app-layout app-layout--sidebar">
    <aside class="sidebar">
      <div class="logo">
        <span class="logo-icon">🦞</span>
        <span class="logo-text">藏叶</span>
      </div>

      <nav class="nav">
        <router-link to="/" class="nav-item" exact-active-class="router-link-active">
          <span class="nav-icon">📊</span>
          <span class="nav-label">仪表盘</span>
        </router-link>
        <router-link to="/files" class="nav-item">
          <span class="nav-icon">📁</span>
          <span class="nav-label">文件管理</span>
        </router-link>
        <router-link to="/settings" class="nav-item">
          <span class="nav-icon">⚙️</span>
          <span class="nav-label">设置</span>
        </router-link>
      </nav>

      <div class="sidebar-footer">
        <span class="status-dot" :class="{ online }" />
        <span class="status-text">{{ online ? '已连接' : '离线' }}</span>
      </div>
    </aside>

    <main class="main">
      <router-view />
    </main>
  </div>

  <!-- ===== 顶部导航布局 ===== -->
  <div v-else class="app-layout app-layout--topbar">
    <header class="topbar">
      <div class="topbar-left">
        <span class="logo-icon">🦞</span>
        <span class="logo-text">藏叶</span>
      </div>

      <nav class="topbar-nav">
        <router-link to="/" class="nav-item" exact-active-class="router-link-active">
          <span class="nav-icon">📊</span>
          <span class="nav-label">仪表盘</span>
        </router-link>
        <router-link to="/files" class="nav-item">
          <span class="nav-icon">📁</span>
          <span class="nav-label">文件管理</span>
        </router-link>
        <router-link to="/settings" class="nav-item">
          <span class="nav-icon">⚙️</span>
          <span class="nav-label">设置</span>
        </router-link>
      </nav>

      <div class="topbar-right">
        <span class="status-dot" :class="{ online }" />
        <span class="status-text">{{ online ? '已连接' : '离线' }}</span>
      </div>
    </header>

    <main class="main">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, inject, onMounted, type Ref } from 'vue'

// 从 App.vue 注入全局布局状态
const layout = inject<Ref<string>>('layout', ref('sidebar'))

// 后端连接状态
const online = ref(false)

onMounted(async () => {
  try {
    const res = await fetch('/api/health')
    online.value = res.ok
  } catch {
    online.value = false
  }
})
</script>

<style scoped>
/* ===== 通用 ===== */
.app-layout { display: flex; min-height: 100vh; height: 100vh; }

/* ===== 侧边栏布局 ===== */
.app-layout--sidebar { flex-direction: row; }

.sidebar {
  width: 220px;
  background: var(--bg-secondary);
  display: flex;
  flex-direction: column;
  border-right: 1px solid var(--border-color);
  flex-shrink: 0;
}

/* ===== 顶部导航布局 ===== */
.app-layout--topbar { flex-direction: column; }

.topbar {
  height: 56px;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  padding: 0 16px;
  flex-shrink: 0;
}

.topbar-left {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-right: 32px;
}

.topbar-nav {
  display: flex;
  align-items: center;
  gap: 4px;
  flex: 1;
}

.topbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: var(--text-muted);
}

/* ===== 共享 ===== */

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 20px 16px;
  border-bottom: 1px solid var(--border-color);
}
.logo-icon { font-size: 28px; }
.logo-text { font-size: 20px; font-weight: 700; color: var(--accent); }

.nav {
  flex: 1;
  padding: 12px 8px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

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

.sidebar-footer {
  padding: 16px;
  border-top: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: var(--text-muted);
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #666;
}
.status-dot.online { background: #4ade80; }

.main {
  flex: 1;
  overflow: auto;
  padding: 24px;
}
</style>