<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

const navItems = [
  { icon: '📊', label: '总览', route: '/overview' },
  { icon: '📁', label: '文件', route: '/' },
  { icon: '⚙️', label: '设置', route: '/settings' },
]

const activeRoute = computed(() => route.path)
const isBrowse = computed(() => route.path.startsWith('/browse/'))

function navigate(item: typeof navItems[0]) { router.push(item.route) }

const theme = ref<'dark' | 'light'>('dark')

function toggleTheme() { theme.value = theme.value === 'dark' ? 'light' : 'dark' }

fetch('/api/settings/theme').then(r => r.json()).then(d => { if (d.value) theme.value = d.value }).catch(() => {})
</script>

<template>
  <div :class="['app-layout', theme]">
    <aside class="sidebar">
      <div class="sidebar-brand">🦞</div>
      <nav class="sidebar-nav">
        <button v-for="item in navItems" :key="item.route"
          :class="['nav-btn', { active: activeRoute === item.route || (isBrowse && item.route === '/') }]"
          :title="item.label" @click="navigate(item)">
          <span class="nav-icon">{{ item.icon }}</span>
        </button>
      </nav>
      <div class="sidebar-footer">
        <button class="nav-btn" title="切换主题" @click="toggleTheme">
          <span class="nav-icon">{{ theme === 'dark' ? '🌙' : '☀️' }}</span>
        </button>
      </div>
    </aside>
    <main class="main-content">
      <router-view :theme="theme" />
    </main>
  </div>
</template>

<style>
.dark {
  --bg-primary: #1a1a2e; --bg-secondary: #16162a; --bg-card: #252540; --bg-card-hover: #2a2a50;
  --bg-input: #1e1e30; --border: #333; --border-light: #2a2a3a; --text-primary: #e0e0e0;
  --text-secondary: #ccc; --text-muted: #888; --text-dim: #666; --accent: #7799cc;
  --accent-bg: #3a3a8a; --accent-hover: #4a4a9a; --danger: #cc5555; --green: #6c6;
  --tag-local-bg: #1a3a1a; --tag-local-text: #6c6; --tag-web-bg: #3a1a1a; --tag-web-text: #c66;
  --preview-bg: #0d0d1a;
}
.light {
  --bg-primary: #f5f5f5; --bg-secondary: #e8e8e8; --bg-card: #fff; --bg-card-hover: #f0f0ff;
  --bg-input: #fff; --border: #ddd; --border-light: #e8e8e8; --text-primary: #222;
  --text-secondary: #444; --text-muted: #888; --text-dim: #aaa; --accent: #4466aa;
  --accent-bg: #5577cc; --accent-hover: #6688dd; --danger: #cc3333; --green: #3a3;
  --tag-local-bg: #e8ffe8; --tag-local-text: #2a2; --tag-web-bg: #ffe8e8; --tag-web-text: #c44;
  --preview-bg: #f0f0f0;
}
* { margin: 0; padding: 0; box-sizing: border-box; }
body { overflow: hidden; }
.app-layout { display: flex; height: 100vh; background: var(--bg-primary); color: var(--text-primary); }
.sidebar { width: 56px; background: var(--bg-secondary); border-right: 1px solid var(--border); display: flex; flex-direction: column; align-items: center; padding: 12px 0; gap: 4px; flex-shrink: 0; }
.sidebar-brand { font-size: 28px; margin-bottom: 8px; cursor: default; }
.sidebar-nav { display: flex; flex-direction: column; gap: 4px; flex: 1; }
.sidebar-footer { margin-top: auto; }
.nav-btn { width: 38px; height: 38px; border-radius: 8px; display: flex; align-items: center; justify-content: center; border: none; background: transparent; cursor: pointer; font-size: 18px; transition: background 0.15s; }
.nav-btn:hover { background: var(--bg-card); }
.nav-btn.active { background: var(--accent-bg); }
.main-content { flex: 1; display: flex; flex-direction: column; overflow: hidden; }
</style>