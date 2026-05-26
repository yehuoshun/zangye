<template>
  <div class="app-layout">
    <aside class="sidebar">
      <div class="logo">
        <span class="logo-icon">🦞</span>
        <span class="logo-text">藏叶</span>
      </div>
      <nav class="nav">
        <router-link to="/" class="nav-item">
          <span class="nav-icon">📁</span>
          <span class="nav-label">文件</span>
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
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

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

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  background: #1a1a2e;
  color: #e0e0e0;
  min-height: 100vh;
}

.app-layout {
  display: flex;
  height: 100vh;
}

.sidebar {
  width: 220px;
  background: #16213e;
  display: flex;
  flex-direction: column;
  border-right: 1px solid #0f3460;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 20px 16px;
  border-bottom: 1px solid #0f3460;
}

.logo-icon {
  font-size: 28px;
}

.logo-text {
  font-size: 20px;
  font-weight: 700;
  color: #e94560;
}

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
  color: #a0a0b0;
  text-decoration: none;
  transition: all 0.2s;
  font-size: 15px;
}

.nav-item:hover {
  background: #0f3460;
  color: #e0e0e0;
}

.nav-item.router-link-active {
  background: #0f3460;
  color: #e94560;
}

.nav-icon {
  font-size: 18px;
}

.sidebar-footer {
  padding: 16px;
  border-top: 1px solid #0f3460;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #606080;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #666;
}

.status-dot.online {
  background: #4ade80;
}

.main {
  flex: 1;
  overflow: auto;
  padding: 24px;
}
</style>