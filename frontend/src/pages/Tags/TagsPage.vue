<template>
  <div class="tags-page">
    <div class="page-header">
      <h2 class="page-title">标签管理</h2>
      <div class="tag-create" style="display: flex; gap: 8px;">
        <input v-model="newTagName" class="input" placeholder="新标签名称..." @keyup.enter="handleCreate" />
        <button class="btn btn-primary" @click="handleCreate">创建</button>
      </div>
    </div>

    <div v-if="loading" class="loading">加载中...</div>

    <div v-else class="tags-grid">
      <div v-for="tag in tags" :key="tag.id" class="tag-card">
        <div class="tag-color-bar" :style="{ background: tag.color }"></div>
        <div class="tag-card-body">
          <div class="tag-card-header">
            <span class="tag-name">{{ tag.name }}</span>
            <span class="tag-count">{{ tag.file_count || 0 }} 个文件</span>
          </div>
          <div class="tag-card-actions">
            <button class="btn btn-sm btn-secondary" @click="editTag(tag)">编辑</button>
            <button class="btn btn-sm btn-danger" @click="handleDelete(tag.id)">删除</button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="tags.length === 0 && !loading" class="empty-state">
      <div class="empty-state-icon">🏷️</div>
      <div class="empty-state-text">暂无标签，输入名称创建</div>
    </div>

    <!-- 编辑弹窗 -->
    <div v-if="editingTag" class="modal-overlay" @click.self="editingTag = null">
      <div class="modal-content">
        <div class="modal-header">
          <h3 class="modal-title">编辑标签</h3>
          <button class="modal-close" @click="editingTag = null">&times;</button>
        </div>
        <div class="form">
          <div class="form-group">
            <label>标签名称</label>
            <input v-model="editForm.name" class="input" />
          </div>
          <div class="form-group">
            <label>颜色</label>
            <div style="display: flex; gap: 8px; align-items: center;">
              <input v-model="editForm.color" class="input" style="width: 100px;" />
              <span class="color-preview" :style="{ background: editForm.color }"></span>
            </div>
          </div>
          <div style="display: flex; gap: 8px; justify-content: flex-end; margin-top: 16px;">
            <button class="btn btn-secondary" @click="editingTag = null">取消</button>
            <button class="btn btn-primary" @click="handleUpdate">保存</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getTags, createTag, updateTag, deleteTag } from '@/features/tags/api'
import type { TagItem } from '@/features/tags/types'

const tags = ref<TagItem[]>([])
const loading = ref(false)
const newTagName = ref('')
const editingTag = ref<TagItem | null>(null)
const editForm = ref({ name: '', color: '' })

onMounted(() => {
  loadTags()
})

async function loadTags() {
  loading.value = true
  try {
    tags.value = await getTags()
  } catch (e) {
    console.error('加载标签失败:', e)
  } finally {
    loading.value = false
  }
}

async function handleCreate() {
  if (!newTagName.value.trim()) return
  try {
    await createTag(newTagName.value.trim())
    newTagName.value = ''
    loadTags()
  } catch (e) {
    console.error('创建标签失败:', e)
  }
}

function editTag(tag: TagItem) {
  editingTag.value = tag
  editForm.value = { name: tag.name, color: tag.color }
}

async function handleUpdate() {
  if (!editingTag.value || !editForm.value.name.trim()) return
  try {
    await updateTag(editingTag.value.id, editForm.value.name, editForm.value.color)
    editingTag.value = null
    loadTags()
  } catch (e) {
    console.error('更新标签失败:', e)
  }
}

async function handleDelete(id: string) {
  if (!confirm('确定删除此标签？')) return
  try {
    await deleteTag(id)
    loadTags()
  } catch (e) {
    console.error('删除标签失败:', e)
  }
}
</script>

<style scoped>
.tags-page {
  max-width: 800px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 12px;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
}

.tags-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 12px;
}

.tag-card {
  display: flex;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius);
  overflow: hidden;
  transition: all var(--transition-fast);
}

.tag-card:hover {
  border-color: var(--accent-primary);
  box-shadow: var(--shadow-sm);
}

.tag-color-bar {
  width: 6px;
  flex-shrink: 0;
}

.tag-card-body {
  flex: 1;
  padding: 16px;
}

.tag-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.tag-name {
  font-size: 15px;
  font-weight: 600;
}

.tag-count {
  font-size: 12px;
  color: var(--text-muted);
}

.tag-card-actions {
  display: flex;
  gap: 4px;
}

.form-group {
  margin-bottom: 12px;
}

.form-group label {
  display: block;
  font-size: 12px;
  font-weight: 600;
  color: var(--text-secondary);
  margin-bottom: 4px;
}

.color-preview {
  width: 32px;
  height: 32px;
  border-radius: var(--border-radius-sm);
  border: 1px solid var(--border-color);
}
</style>
