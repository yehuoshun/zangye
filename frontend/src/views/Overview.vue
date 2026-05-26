<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { overview, formatSize, type OverviewStats } from '../api'

const stats = ref<OverviewStats | null>(null)
const loading = ref(true)

onMounted(async () => {
  try { stats.value = await overview.stats() } catch (e) { console.error(e) }
  finally { loading.value = false }
})
</script>

<template>
  <div class="overview-page">
    <div class="page-header"><h2>📊 总览</h2></div>
    <div v-if="loading" class="status-text">加载中…</div>
    <div v-else-if="stats" class="overview-content">
      <div class="stat-cards">
        <div class="stat-card"><div class="stat-value">{{ stats.total_files.toLocaleString() }}</div><div class="stat-label">总文件数</div></div>
        <div class="stat-card"><div class="stat-value">{{ formatSize(stats.total_size) }}</div><div class="stat-label">总大小</div></div>
        <div class="stat-card"><div class="stat-value">{{ stats.collection_count }}</div><div class="stat-label">文件夹数</div></div>
      </div>
      <div class="category-cards">
        <div v-for="cat in stats.categories" :key="cat.name" class="category-card">
          <div class="cat-label">{{ cat.label }}</div>
          <div class="cat-count">{{ cat.count.toLocaleString() }}</div>
          <div class="cat-size">{{ formatSize(cat.size) }}</div>
          <div class="cat-bar"><div class="cat-bar-fill" :style="{ width: cat.percent.toFixed(1) + '%' }"></div></div>
        </div>
      </div>
    </div>
    <div v-else class="status-text">暂无数据</div>
  </div>
</template>

<style scoped>
.overview-page { padding: 20px; overflow-y: auto; height: 100%; }
.page-header { margin-bottom: 20px; border-bottom: 1px solid var(--border); padding-bottom: 12px; }
.page-header h2 { font-size: 18px; font-weight: 500; }
.stat-cards { display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px; margin-bottom: 24px; }
.stat-card { background: var(--bg-card); border-radius: 10px; padding: 20px; text-align: center; }
.stat-value { font-size: 28px; font-weight: bold; color: var(--accent); }
.stat-label { font-size: 12px; color: var(--text-muted); margin-top: 4px; }
.category-cards { display: grid; grid-template-columns: repeat(auto-fill, minmax(240px, 1fr)); gap: 12px; }
.category-card { background: var(--bg-card); border: 1px solid var(--border-light); border-radius: 10px; padding: 16px; }
.cat-label { font-size: 14px; margin-bottom: 6px; }
.cat-count { font-size: 20px; font-weight: bold; }
.cat-size { font-size: 12px; color: var(--text-muted); }
.cat-bar { margin-top: 8px; background: var(--bg-secondary); border-radius: 3px; height: 4px; overflow: hidden; }
.cat-bar-fill { height: 100%; background: var(--accent); border-radius: 3px; transition: width 0.3s; }
.status-text { color: var(--text-muted); text-align: center; padding: 60px; font-size: 14px; }
</style>
