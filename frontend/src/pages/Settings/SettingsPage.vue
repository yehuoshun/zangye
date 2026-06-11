<!--
  SettingsPage.vue — 系统设置页面

  支持修改：
    - 主题（深色/浅色）
    - 布局模式
    - 数据库版本号（只读）

  数据流：
    页面加载 → fetchSettings() → GET /api/settings → 渲染表单
    用户修改 → updateSettings() → PUT /api/settings → 保存到数据库
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
        <select v-model="form.theme" class="setting-input">
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
        <select v-model="form.layout" class="setting-input">
          <option value="sidebar">侧边栏</option>
          <option value="topbar">顶部导航</option>
        </select>
      </div>

      <!-- 数据库版本（只读） -->
      <div class="setting-item">
        <div class="setting-label">
          <span class="setting-name">数据库版本</span>
          <span class="setting-desc">当前数据库结构版本号</span>
        </div>
        <span class="setting-value">{{ form.version }}</span>
      </div>

      <!-- 保存按钮 -->
      <div class="setting-actions">
        <button class="btn-save" @click="save" :disabled="saving">
          {{ saving ? '保存中…' : '保存设置' }}
        </button>
        <span v-if="saved" class="saved-tip">✅ 已保存</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { fetchSettings, updateSettings } from '@/features/settings/api'

const form = reactive({ theme: 'dark', layout: 'sidebar', version: '' })
const saving = ref(false)
const saved = ref(false)

onMounted(async () => {
  try {
    const data = await fetchSettings()
    form.theme = data.theme || 'dark'
    form.layout = data.layout || 'sidebar'
    form.version = data.version || ''
  } catch (e) {
    console.error('加载设置失败', e)
  }
})

async function save() {
  saving.value = true
  saved.value = false
  try {
    await updateSettings({ theme: form.theme, layout: form.layout })
    saved.value = true
    setTimeout(() => saved.value = false, 2000)
  } catch (e) {
    console.error('保存设置失败', e)
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.settings-page { max-width: 640px; }
.page-title { font-size: 24px; font-weight: 600; margin-bottom: 24px; }

.settings-card {
  background: #16213e;
  border: 1px solid #0f3460;
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
  border-bottom: 1px solid #0f3460;
}
.setting-item:last-of-type { border-bottom: none; }

.setting-label { display: flex; flex-direction: column; gap: 4px; }
.setting-name { font-size: 15px; font-weight: 500; }
.setting-desc { font-size: 12px; color: #606080; }

.setting-input, .setting-value {
  background: #0f3460;
  border: 1px solid #0f3460;
  color: #e0e0e0;
  padding: 8px 12px;
  border-radius: 8px;
  font-size: 14px;
  outline: none;
  min-width: 140px;
}
.setting-input:focus { border-color: #60a5fa; }
.setting-value { color: #606080; }

.setting-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  padding-top: 4px;
}

.btn-save {
  background: #60a5fa;
  color: #0a0a1a;
  border: none;
  padding: 10px 24px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: opacity .2s;
}
.btn-save:hover { opacity: .85; }
.btn-save:disabled { opacity: .5; cursor: not-allowed; }

.saved-tip { font-size: 14px; color: #4ade80; }
</style>