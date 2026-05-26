<template>
  <div class="dashboard">
    <h1 class="page-title">仪表盘</h1>

    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">📁</div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.file_count }}</div>
          <div class="stat-label">文件总数</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">📂</div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.collection_count }}</div>
          <div class="stat-label">集合数</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">🏷️</div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.tag_count }}</div>
          <div class="stat-label">标签数</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">💾</div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.storage_display }}</div>
          <div class="stat-label">存储空间</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, onMounted } from 'vue'

interface DashboardStats {
  file_count: number
  collection_count: number
  tag_count: number
  storage_bytes: number
  storage_display: string
}

const stats = reactive<DashboardStats>({
  file_count: 0,
  collection_count: 0,
  tag_count: 0,
  storage_bytes: 0,
  storage_display: '0 B',
})

onMounted(async () => {
  try {
    const res = await fetch('/api/dashboard/stats')
    const data = await res.json()
    Object.assign(stats, data)
  } catch (e) {
    console.error('仪表盘数据加载失败', e)
  }
})
</script>

<style scoped>
.dashboard {
  max-width: 1200px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #e0e0e0;
  margin-bottom: 24px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 16px;
}

.stat-card {
  background: #16213e;
  border: 1px solid #0f3460;
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  transition: border-color 0.2s;
}

.stat-card:hover {
  border-color: #e94560;
}

.stat-icon {
  font-size: 36px;
  width: 56px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #0f3460;
  border-radius: 12px;
}

.stat-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #e94560;
}

.stat-label {
  font-size: 13px;
  color: #606080;
}
</style>