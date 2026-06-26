<template>
  <div class="layout">
    <!-- 侧边栏 -->
    <aside class="sidebar">
      <!-- Logo -->
      <div class="logo" @click="$router.push('/')">
        <span class="logo-icon">🍃</span>
        <span class="logo-text">藏叶</span>
      </div>

      <!-- 导航 -->
      <nav class="nav">
        <router-link to="/" class="nav-item" :class="{ active: $route.path === '/' }">
          <span class="nav-icon">📊</span>
          <span>仪表盘</span>
        </router-link>
        <router-link to="/files" class="nav-item" :class="{ active: $route.path === '/files' }">
          <span class="nav-icon">📁</span>
          <span>文件</span>
        </router-link>
        <router-link to="/tags" class="nav-item" :class="{ active: $route.path === '/tags' }">
          <span class="nav-icon">🏷️</span>
          <span>标签</span>
        </router-link>
        <router-link to="/trash" class="nav-item" :class="{ active: $route.path === '/trash' }">
          <span class="nav-icon">🗑️</span>
          <span>回收站</span>
        </router-link>
        <router-link to="/settings" class="nav-item" :class="{ active: $route.path === '/settings' }">
          <span class="nav-icon">⚙️</span>
          <span>设置</span>
        </router-link>
      </nav>

      <!-- 文件夹树 -->
      <div class="folder-section">
        <div class="folder-header">
          <span>文件夹</span>
          <button class="btn btn-icon btn-sm" @click="showCreateFolder = true" title="新建文件夹">+</button>
        </div>
        <div class="folder-tree" v-if="folderTree.length > 0">
          <FolderTreeNode
            v-for="node in folderTree"
            :key="node.id"
            :node="node"
            :depth="0"
            :selected-id="selectedFolderId"
            @select="onFolderSelect"
          />
        </div>
        <div v-else class="folder-empty">暂无文件夹</div>
      </div>
    </aside>

    <!-- 主区域 -->
    <div class="main">
      <!-- 顶部栏 -->
      <header class="topbar">
        <div class="search-box">
          <span class="search-icon">🔍</span>
          <input
            v-model="searchKeyword"
            class="input"
            placeholder="搜索文件..."
            @keyup.enter="onSearch"
          />
        </div>
        <div class="topbar-actions">
          <button class="btn btn-secondary btn-sm" @click="viewMode = 'grid'" :class="{ active: viewMode === 'grid' }">📇</button>
          <button class="btn btn-secondary btn-sm" @click="viewMode = 'list'" :class="{ active: viewMode === 'list' }">📋</button>
        </div>
      </header>

      <!-- 内容区 -->
      <main class="content">
        <router-view
          :view-mode="viewMode"
          :search-keyword="searchKeyword"
          :selected-folder-id="selectedFolderId"
          @search="onSearch"
          @folder-select="onFolderSelect"
        />
      </main>
    </div>

    <!-- 新建文件夹弹窗 -->
    <div v-if="showCreateFolder" class="modal-overlay" @click.self="showCreateFolder = false">
      <div class="modal-content">
        <div class="modal-header">
          <h3 class="modal-title">新建文件夹</h3>
          <button class="modal-close" @click="showCreateFolder = false">&times;</button>
        </div>
        <div style="margin-bottom: 12px;">
          <input v-model="newFolderName" class="input" placeholder="文件夹名称" @keyup.enter="handleCreateFolder" />
        </div>
        <div style="display: flex; gap: 8px; justify-content: flex-end;">
          <button class="btn btn-secondary" @click="showCreateFolder = false">取消</button>
          <button class="btn btn-primary" @click="handleCreateFolder">创建</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getFolderTree, createFolder } from '@/features/folders/api'
import type { FolderItem } from '@/features/folders/types'
import FolderTreeNode from './FolderTreeNode.vue'

const props = defineProps<{
  viewMode?: string
  searchKeyword?: string
  selectedFolderId?: string
}>()

const emit = defineEmits<{
  search: [keyword: string]
  folderSelect: [folderId: string | null]
}>()

const folderTree = ref<FolderItem[]>([])
const selectedFolderId = ref<string | null>(null)
const searchKeyword = ref('')
const viewMode = ref('list')
const showCreateFolder = ref(false)
const newFolderName = ref('')

onMounted(async () => {
  try {
    folderTree.value = await getFolderTree()
  } catch (e) {
    console.error('加载文件夹树失败:', e)
  }
})

function onFolderSelect(id: string | null) {
  selectedFolderId.value = id
  emit('folderSelect', id)
}

function onSearch() {
  emit('search', searchKeyword.value)
}

async function handleCreateFolder() {
  if (!newFolderName.value.trim()) return
  try {
    await createFolder(newFolderName.value.trim())
    newFolderName.value = ''
    showCreateFolder.value = false
    folderTree.value = await getFolderTree()
  } catch (e) {
    console.error('创建文件夹失败:', e)
  }
}
</script>

<script lang="ts">
// 需要导出组件定义以便递归使用
export default {
  name: 'MainLayout',
}
</script>

<style scoped>
.layout {
  display: flex;
  height: 100vh;
  overflow: hidden;
}

.sidebar {
  width: var(--sidebar-width);
  background: var(--bg-secondary);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.logo {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 16px 20px;
  cursor: pointer;
  border-bottom: 1px solid var(--border-color);
}

.logo-icon {
  font-size: 24px;
}

.logo-text {
  font-size: 18px;
  font-weight: 700;
  color: var(--accent-primary);
}

.nav {
  padding: 8px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: var(--border-radius-sm);
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 14px;
  transition: all var(--transition-fast);
}

.nav-item:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.nav-item.active {
  background: var(--bg-active);
  color: var(--accent-primary);
}

.nav-icon {
  font-size: 16px;
  width: 20px;
  text-align: center;
}

.folder-section {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.folder-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  font-size: 12px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.folder-empty {
  padding: 20px 12px;
  text-align: center;
  color: var(--text-muted);
  font-size: 12px;
}

.main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.topbar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 20px;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
  height: var(--topbar-height);
}

.search-box {
  flex: 1;
  position: relative;
  max-width: 400px;
}

.search-icon {
  position: absolute;
  left: 10px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 14px;
}

.search-box .input {
  padding-left: 32px;
}

.topbar-actions {
  display: flex;
  gap: 4px;
}

.topbar-actions .btn.active {
  background: var(--accent-primary);
  color: #fff;
}

.content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}
</style>
