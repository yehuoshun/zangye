<template>
  <div class="dashboard">
    <h2 class="page-title">仪表盘</h2>

    <div v-if="loading" class="loading">加载中...</div>

    <div v-else class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">📂</div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.folder_count }}</div>
          <div class="stat-label">文件夹</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">📄</div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.file_count }}</div>
          <div class="stat-label">文件总数</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">🖼️</div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.image_count }}</div>
          <div class="stat-label">图片</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">🎬</div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.video_count }}</div>
          <div class="stat-label">视频</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">🎵</div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.audio_count }}</div>
          <div class="stat-label">音频</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">📦</div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.other_count }}</div>
          <div class="stat-label">其他</div>
        </div>
      </div>

      <div class="stat-card stat-card-wide">
        <div class="stat-icon">💾</div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.size_text }}</div>
          <div class="stat-label">总存储空间</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getDashboardStats } from '@/features/dashboard/api'
import type { DashboardStats } from '@/features/dashboard/types'

const loading = ref(true)
const stats = ref<DashboardStats>({
  folder_count: 0,
  file_count: 0,
  image_count: 0,
  video_count: 0,
  audio_count: 0,
  other_count: 0,
  total_size: 0,
  size_text: '0 B',
})

onMounted(async () => {
  try {
    stats.value = await getDashboardStats()
  } catch (e) {
    console.error('加载仪表盘数据失败:', e)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.dashboard {
  max-width: 900px;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 24px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius);
  padding: 20px;
  transition: all var(--transition-fast);
}

.stat-card:hover {
  border-color: var(--accent-primary);
  box-shadow: var(--shadow-sm);
}

.stat-card-wide {
  grid-column: 1 / -1;
}

.stat-icon {
  font-size: 32px;
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-tertiary);
  border-radius: var(--border-radius);
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1.2;
}

.stat-label {
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 4px;
}
</style>
