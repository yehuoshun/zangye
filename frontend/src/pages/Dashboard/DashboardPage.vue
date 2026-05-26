<template>
  <div class="dashboard">
    <h1 class="page-title">仪表盘</h1>
    <div class="stats-grid">
      <div class="stat-card" v-for="card in cards" :key="card.label">
        <div class="stat-icon">{{ card.icon }}</div>
        <div class="stat-info">
          <div class="stat-value">{{ card.value }}</div>
          <div class="stat-label">{{ card.label }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { DashboardStats } from '@/features/dashboard/types'
import { fetchDashboardStats } from '@/features/dashboard/api'

const stats = ref<DashboardStats>({
  file_count: 0, collection_count: 0, tag_count: 0,
  storage_bytes: 0, storage_display: '0 B',
})

const cards = computed(() => [
  { icon: '📁', value: stats.value.file_count, label: '文件总数' },
  { icon: '📂', value: stats.value.collection_count, label: '集合数' },
  { icon: '🏷️', value: stats.value.tag_count, label: '标签数' },
  { icon: '💾', value: stats.value.storage_display, label: '存储空间' },
])

onMounted(async () => {
  try { stats.value = await fetchDashboardStats() }
  catch (e) { console.error('仪表盘数据加载失败', e) }
})
</script>

<style scoped>
.dashboard { max-width: 1200px; }
.page-title { font-size: 24px; font-weight: 600; margin-bottom: 24px; }
.stats-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(240px, 1fr)); gap: 16px; }
.stat-card {
  background: #16213e; border: 1px solid #0f3460; border-radius: 12px;
  padding: 20px; display: flex; align-items: center; gap: 16px;
}
.stat-card:hover { border-color: #e94560; }
.stat-icon {
  font-size: 36px; width: 56px; height: 56px; display: flex;
  align-items: center; justify-content: center; background: #0f3460; border-radius: 12px;
}
.stat-info { display: flex; flex-direction: column; gap: 4px; }
.stat-value { font-size: 28px; font-weight: 700; color: #e94560; }
.stat-label { font-size: 13px; color: #606080; }
</style>