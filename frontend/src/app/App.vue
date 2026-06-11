<!--
  App.vue — 藏叶前端根组件

  管理全局状态（主题、布局），通过 provide 下发到子组件。
  设置页修改后无需刷新即可生效。
-->

<template>
  <router-view />
</template>

<script setup lang="ts">
import { ref, provide, onMounted } from 'vue'
import { fetchSettings } from '@/features/settings/api'

// 全局主题状态
const theme = ref(localStorage.getItem('zangye-theme') || 'dark')
// 全局布局状态
const layout = ref(localStorage.getItem('zangye-layout') || 'sidebar')

// 提供给后代组件
provide('theme', theme)
provide('layout', layout)

// 主题变更时立即应用到 DOM
function applyTheme(t: string) {
  document.documentElement.setAttribute('data-theme', t)
  localStorage.setItem('zangye-theme', t)
}

// 布局变更时缓存
function persistLayout(l: string) {
  localStorage.setItem('zangye-layout', l)
}

// 提供修改方法给子组件
provide('setTheme', (t: string) => {
  theme.value = t
  applyTheme(t)
})
provide('setLayout', (l: string) => {
  layout.value = l
  persistLayout(l)
})

onMounted(async () => {
  applyTheme(theme.value)
  persistLayout(layout.value)
  try {
    const settings = await fetchSettings()
    const t = settings.theme || 'dark'
    const l = settings.layout || 'sidebar'
    theme.value = t
    layout.value = l
    applyTheme(t)
    persistLayout(l)
  } catch { /* 使用缓存值 */ }
})
</script>