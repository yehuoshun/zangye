<!--
  DashboardPage.vue — 仪表盘页面

  展示藏叶的核心统计数据概览：
    - 文件总数
    - 集合数
    - 标签数
    - 存储空间

  数据通过 fetchDashboardStats() 从后端 API 获取，
  以卡片网格形式展示，支持响应式布局。

  数据流：
    DashboardPage → fetchDashboardStats() → GET /api/dashboard/stats
    → DashboardHandler.Stats → MySQL 查询 → JSON 响应
-->

<template>
  <div class="dashboard">
    <!-- 页面标题 -->
    <h1 class="page-title">仪表盘</h1>

    <!-- 统计卡片网格 -->
    <div class="stats-grid">
      <!--
        遍历 cards 计算属性渲染统计卡片。
        每张卡片包含：图标、数值、标签。
      -->
      <div class="stat-card" v-for="card in cards" :key="card.label">
        <!-- 图标区域 -->
        <div class="stat-icon">{{ card.icon }}</div>
        <!-- 信息区域：数值 + 标签 -->
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

// 统计数据响应式状态，初始值为 0
const stats = ref<DashboardStats>({
  folder_count: 0,
  file_count: 0,
  image_count: 0,
  video_count: 0,
  audio_count: 0,
  other_count: 0,
  storage_bytes: 0,
  storage_display: '0 B',
})

/**
 * 计算属性：将原始统计数据映射为卡片展示数据。
 *
 * 每张卡片包含：
 *   - icon: emoji 图标
 *   - value: 显示数值（数字或格式化字符串）
 *   - label: 中文标签
 */
const cards = computed(() => [
  { icon: '📂', value: stats.value.folder_count, label: '文件夹' },
  { icon: '📁', value: stats.value.file_count, label: '文件总数' },
  { icon: '🖼️', value: stats.value.image_count, label: '图片' },
  { icon: '🎬', value: stats.value.video_count, label: '视频' },
  { icon: '🎵', value: stats.value.audio_count, label: '音频' },
  { icon: '📄', value: stats.value.other_count, label: '其他' },
  { icon: '💾', value: stats.value.storage_display, label: '存储空间' },
])

// 组件挂载后异步加载仪表盘数据
onMounted(async () => {
  try {
    stats.value = await fetchDashboardStats()
  } catch (e) {
    // 加载失败时保持默认值（全 0），仅打印错误日志
    console.error('仪表盘数据加载失败', e)
  }
})
</script>

<style scoped>
/* 仪表盘容器：限制最大宽度，居中显示 */
.dashboard { max-width: 1200px; }

/* 页面标题 */
.page-title { font-size: 24px; font-weight: 600; margin-bottom: 24px; }

/* 统计卡片网格：响应式布局，最小列宽 240px */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 16px;
}

/* 单张统计卡片 */
.stat-card {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
}
.stat-card:hover { border-color: var(--accent); }

.stat-icon {
  font-size: 36px;
  width: 56px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-tertiary);
  border-radius: 12px;
}

.stat-info { display: flex; flex-direction: column; gap: 4px; }
.stat-value { font-size: 28px; font-weight: 700; color: var(--accent); }
.stat-label { font-size: 13px; color: var(--text-muted); }
</style>