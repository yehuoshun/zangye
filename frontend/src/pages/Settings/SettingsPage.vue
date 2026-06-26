<template>
  <div class="settings-page">
    <h2 class="page-title">设置</h2>

    <div v-if="loading" class="loading">加载中...</div>

    <div v-else class="settings-sections">
      <div class="card settings-section">
        <h3 class="section-title">显示设置</h3>
        <div class="setting-item">
          <div class="setting-label">
            <div class="setting-name">默认视图</div>
            <div class="setting-desc">文件列表默认显示方式</div>
          </div>
          <select v-model="settings.default_view" class="input" style="width: 160px;">
            <option value="list">列表视图</option>
            <option value="grid">网格视图</option>
          </select>
        </div>
        <div class="setting-item">
          <div class="setting-label">
            <div class="setting-name">每页数量</div>
            <div class="setting-desc">文件列表每页显示条数</div>
          </div>
          <input v-model.number="settings.items_per_page" class="input" type="number" style="width: 100px;" min="10" max="200" />
        </div>
      </div>

      <div class="card settings-section">
        <h3 class="section-title">打开方式</h3>
        <div class="setting-desc" style="margin-bottom: 12px; color: var(--text-muted);">
          配置不同文件类型的外部打开程序
        </div>
        <div v-if="openWithPrograms.length === 0" class="empty-state" style="padding: 20px;">
          <div class="empty-state-text">暂无配置</div>
        </div>
        <div v-for="(prog, index) in openWithPrograms" :key="index" class="open-with-item">
          <input v-model="prog.extension" class="input" placeholder="扩展名（如 .txt）" style="width: 120px;" />
          <input v-model="prog.program" class="input" placeholder="程序路径" style="flex: 1;" />
          <button class="btn btn-sm btn-danger" @click="removeOpenWith(index)">删除</button>
        </div>
        <button class="btn btn-sm btn-secondary" style="margin-top: 8px;" @click="addOpenWith">+ 添加</button>
      </div>

      <div style="display: flex; gap: 8px; justify-content: flex-end; margin-top: 24px;">
        <button class="btn btn-primary" @click="saveSettings">保存设置</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getSettings, updateSettings } from '@/features/settings/api'

const loading = ref(false)
const settings = ref<Record<string, any>>({
  default_view: 'list',
  items_per_page: 50,
})

interface OpenWithProgram {
  extension: string
  program: string
}

const openWithPrograms = ref<OpenWithProgram[]>([])

onMounted(() => {
  loadSettings()
})

async function loadSettings() {
  loading.value = true
  try {
    const data = await getSettings()
    if (data.default_view) settings.value.default_view = data.default_view
    if (data.items_per_page) settings.value.items_per_page = parseInt(data.items_per_page)
    if (data.open_with_programs) {
      try {
        openWithPrograms.value = JSON.parse(data.open_with_programs)
      } catch {
        openWithPrograms.value = []
      }
    }
  } catch (e) {
    console.error('加载设置失败:', e)
  } finally {
    loading.value = false
  }
}

function addOpenWith() {
  openWithPrograms.value.push({ extension: '', program: '' })
}

function removeOpenWith(index: number) {
  openWithPrograms.value.splice(index, 1)
}

async function saveSettings() {
  try {
    await updateSettings({
      default_view: settings.value.default_view,
      items_per_page: String(settings.value.items_per_page),
      open_with_programs: JSON.stringify(openWithPrograms.value),
    })
    alert('设置已保存')
  } catch (e) {
    console.error('保存设置失败:', e)
  }
}
</script>

<style scoped>
.settings-page {
  max-width: 700px;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 24px;
}

.settings-section {
  margin-bottom: 20px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 16px;
}

.setting-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 0;
  border-bottom: 1px solid var(--border-color);
}

.setting-item:last-child {
  border-bottom: none;
}

.setting-name {
  font-size: 14px;
  font-weight: 500;
}

.setting-desc {
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 2px;
}

.open-with-item {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 8px;
}
</style>
