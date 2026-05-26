<template>
  <div class="home">
    <h1 class="title">藏叶</h1>
    <p class="subtitle">个人文件管理器</p>
    <div class="status-card">
      <span class="status-label">服务状态</span>
      <span class="status-value" :class="statusClass">{{ statusText }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const statusText = ref('检测中…')
const statusClass = ref('')

onMounted(async () => {
  try {
    const res = await fetch('/api/health')
    const data = await res.json()
    if (data.status === 'ok') {
      statusText.value = '运行正常'
      statusClass.value = 'ok'
    } else {
      statusText.value = '数据库异常'
      statusClass.value = 'error'
    }
  } catch {
    statusText.value = '无法连接'
    statusClass.value = 'error'
  }
})
</script>

<style scoped>
.home {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  gap: 16px;
}

.title {
  font-size: 48px;
  font-weight: 800;
  color: #e94560;
}

.subtitle {
  font-size: 18px;
  color: #606080;
}

.status-card {
  margin-top: 20px;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 24px;
  background: #16213e;
  border-radius: 12px;
  border: 1px solid #0f3460;
}

.status-label {
  color: #606080;
  font-size: 14px;
}

.status-value {
  font-size: 14px;
  font-weight: 600;
}

.status-value.ok {
  color: #4ade80;
}

.status-value.error {
  color: #f87171;
}
</style>