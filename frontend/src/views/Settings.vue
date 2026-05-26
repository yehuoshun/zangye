<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { settings, type PrefixConfig } from '../api'

const themeSetting = ref('dark')
const layoutSetting = ref('sidebar')
const defaultView = ref('grid')
const prefixes = ref<PrefixConfig[]>([])
const activeTab = ref<'prefixes' | 'display'>('prefixes')

onMounted(async () => {
  try {
    const [t, l, v, p] = await Promise.all([
      settings.get('theme').catch(() => ({ key: 'theme', value: 'dark' })),
      settings.get('layout').catch(() => ({ key: 'layout', value: 'sidebar' })),
      settings.get('default_view').catch(() => ({ key: 'default_view', value: 'grid' })),
      settings.prefixes(),
    ])
    themeSetting.value = t.value; layoutSetting.value = l.value; defaultView.value = v.value; prefixes.value = p
  } catch (e) { console.error(e) }
})

async function saveTheme() { await settings.set('theme', themeSetting.value) }
async function saveLayout() { await settings.set('layout', layoutSetting.value) }
async function saveDefaultView() { await settings.set('default_view', defaultView.value) }
async function updatePrefix(prefix: PrefixConfig) {
  await settings.updatePrefix(prefix.prefix, { type: prefix.type, map_path: prefix.map_path, url_template: prefix.url_template })
}
</script>

<template>
  <div class="settings-page">
    <div class="page-header"><h2>⚙️ 设置</h2></div>
    <div class="settings-tabs">
      <button :class="['tab', { active: activeTab === 'prefixes' }]" @click="activeTab = 'prefixes'">前缀映射</button>
      <button :class="['tab', { active: activeTab === 'display' }]" @click="activeTab = 'display'">显示</button>
    </div>
    <div v-if="activeTab === 'prefixes'" class="settings-content">
      <table class="prefix-table" v-if="prefixes.length">
        <thead><tr><th>前缀</th><th>类型</th><th>本地映射路径</th><th>URL 模板</th><th>操作</th></tr></thead>
        <tbody>
          <tr v-for="p in prefixes" :key="p.prefix">
            <td><code class="prefix">{{ p.prefix }}</code></td>
            <td><select v-model="p.type" class="input-sm"><option value="local">本地</option><option value="web">网页</option></select></td>
            <td><input v-if="p.type === 'local'" v-model="p.map_path" class="input-sm" placeholder="如 Z:\115" /><span v-else class="text-dim">—</span></td>
            <td><input v-if="p.type === 'web'" v-model="p.url_template" class="input-sm" placeholder="如 https://t.me/{path}" /><span v-else class="text-dim">—</span></td>
            <td><button class="btn btn-sm" @click="updatePrefix(p)">保存</button></td>
          </tr>
        </tbody>
      </table>
      <div v-else class="status-text">暂无前缀配置</div>
    </div>
    <div v-if="activeTab === 'display'" class="settings-content">
      <div class="setting-row"><span class="setting-label">主题</span><select v-model="themeSetting" class="input" @change="saveTheme"><option value="dark">暗色</option><option value="light">亮色</option></select></div>
      <div class="setting-row"><span class="setting-label">布局</span><select v-model="layoutSetting" class="input" @change="saveLayout"><option value="sidebar">侧边栏</option><option value="top">顶部导航</option></select></div>
      <div class="setting-row"><span class="setting-label">默认视图</span><select v-model="defaultView" class="input" @change="saveDefaultView"><option value="grid">网格</option><option value="list">列表</option></select></div>
    </div>
  </div>
</template>

<style scoped>
.settings-page { padding: 20px; overflow-y: auto; height: 100%; }
.page-header { margin-bottom: 16px; border-bottom: 1px solid var(--border); padding-bottom: 12px; }
.page-header h2 { font-size: 18px; font-weight: 500; }
.settings-tabs { display: flex; border-bottom: 1px solid var(--border-light); gap: 0; }
.tab { padding: 10px 20px; font-size: 13px; border: none; background: none; color: var(--text-muted); cursor: pointer; border-bottom: 2px solid transparent; }
.tab.active { color: var(--accent); border-bottom-color: var(--accent); }
.settings-content { padding: 20px 0; }
.setting-row { display: flex; align-items: center; padding: 12px 0; border-bottom: 1px solid var(--border-light); gap: 16px; }
.setting-label { width: 100px; font-size: 13px; color: var(--text-secondary); flex-shrink: 0; }
.input { background: var(--bg-input); border: 1px solid var(--border); border-radius: 6px; padding: 6px 12px; color: var(--text-primary); font-size: 13px; outline: none; min-width: 160px; }
.prefix-table { width: 100%; border-collapse: collapse; font-size: 13px; }
.prefix-table th { text-align: left; padding: 8px 12px; color: var(--text-muted); border-bottom: 1px solid var(--border); font-weight: normal; font-size: 12px; }
.prefix-table td { padding: 8px 12px; border-bottom: 1px solid var(--border-light); vertical-align: middle; }
.prefix { color: var(--accent); font-family: 'Consolas', 'Courier New', monospace; font-size: 13px; }
.input-sm { background: var(--bg-input); border: 1px solid var(--border); border-radius: 4px; padding: 4px 8px; color: var(--text-primary); font-size: 12px; outline: none; width: 100%; }
select.input-sm { width: auto; }
.text-dim { color: var(--text-dim); font-size: 12px; }
.btn { padding: 6px 14px; border-radius: 6px; border: 1px solid var(--border); background: var(--bg-card); color: var(--text-secondary); cursor: pointer; font-size: 13px; }
.btn:hover { background: var(--bg-card-hover); }
.btn-sm { padding: 4px 10px; font-size: 12px; }
.status-text { color: var(--text-muted); text-align: center; padding: 40px; }
</style>
