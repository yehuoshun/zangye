<!--
  SettingsPage.vue — 系统设置页面

  支持修改主题（深色/浅色）和布局模式（侧边栏/顶部导航）。
  修改后立即全局生效，无需刷新页面。
-->

<template>
  <div class="settings-page">
    <h1 class="page-title">系统设置</h1>

    <div class="settings-card">
      <!-- 主题设置 -->
      <div class="setting-item">
        <div class="setting-label">
          <span class="setting-name">主题</span>
          <span class="setting-desc">界面颜色方案</span>
        </div>
        <select v-model="form.theme" class="setting-input" @change="autoSave">
          <option value="dark">深色</option>
          <option value="light">浅色</option>
        </select>
      </div>

      <!-- 布局设置 -->
      <div class="setting-item">
        <div class="setting-label">
          <span class="setting-name">布局模式</span>
          <span class="setting-desc">侧边栏或顶部导航</span>
        </div>
        <select v-model="form.layout" class="setting-input" @change="autoSave">
          <option value="sidebar">侧边栏</option>
          <option value="topbar">顶部导航</option>
        </select>
      </div>

      <!-- 保存状态 -->
      <div class="setting-status" v-if="status">{{ status }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, inject } from 'vue'
import { fetchSettings, updateSettings } from '@/features/settings/api'

const form = reactive({ theme: 'dark', layout: 'sidebar' })
const status = ref('')

// 从 App.vue 注入全局状态修改方法
const setTheme = inject<(t: string) => void>('setTheme', () => {})
const setLayout = inject<(l: string) => void>('setLayout', () => {})

onMounted(async () => {
  try {
    const data = await fetchSettings()
    form.theme = data.theme || 'dark'
    form.layout = data.layout || 'sidebar'
  } catch (e) {
    console.error('加载设置失败', e)
  }
})

async function autoSave() {
  // 立即全局生效
  setTheme(form.theme)
  setLayout(form.layout)
  try {
    await updateSettings({ theme: form.theme, layout: form.layout })
    status.value = '✅ 已保存'
    setTimeout(() => status.value = '', 2000)
  } catch (e) {
    status.value = '❌ 保存失败'
    setTimeout(() => status.value = '', 2000)
  }
}
</script>

<style scoped>
.settings-page { max-width: 640px; }
.page-title { font-size: 24px; font-weight: 600; margin-bottom: 24px; }

.settings-card {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.setting-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--border-color);
}
.setting-item:last-of-type { border-bottom: none; }

.setting-label { display: flex; flex-direction: column; gap: 4px; }
.setting-name { font-size: 15px; font-weight: 500; }
.setting-desc { font-size: 12px; color: var(--text-muted); }

.setting-input {
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  color: var(--text-primary);
  padding: 8px 12px;
  border-radius: 8px;
  font-size: 14px;
  outline: none;
  min-width: 140px;
}
.setting-input:focus { border-color: var(--accent); }

.setting-status {
  font-size: 13px;
  color: #4ade80;
  padding-top: 4px;
}
</style>