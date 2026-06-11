<!--
  TagsPage.vue — 标签管理页面

  支持标签的创建、编辑、删除，颜色选择。
-->

<template>
  <div class="tags-page">
    <div class="page-header">
      <h1 class="page-title">标签管理</h1>
      <button class="btn btn-primary" @click="openCreate">+ 新建标签</button>
    </div>

    <!-- 标签网格 -->
    <div class="tag-grid" v-if="tags.length">
      <div
        v-for="t in tags"
        :key="t.id"
        class="tag-card"
        :style="{ borderLeftColor: t.color }"
      >
        <div class="tag-dot" :style="{ background: t.color }"></div>
        <span class="tag-name">{{ t.name }}</span>
        <div class="tag-actions">
          <button class="btn btn-sm" @click="openEdit(t)">编辑</button>
          <button class="btn btn-sm btn-danger" @click="confirmDelete(t)">删除</button>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div class="empty-state" v-else>
      <span class="empty-icon">🏷️</span>
      <span class="empty-text">还没有标签，点击上方按钮创建</span>
    </div>

    <!-- 新建/编辑弹窗 -->
    <div class="modal-overlay" v-if="showModal" @click.self="closeModal">
      <div class="modal modal-sm">
        <div class="modal-header">
          <h2>{{ editing ? '编辑标签' : '新建标签' }}</h2>
          <button class="btn-close" @click="closeModal">✕</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>标签名称 <span class="required">*</span></label>
            <input v-model="form.name" class="form-input" placeholder="输入标签名称" @keyup.enter="save" />
          </div>
          <div class="form-group">
            <label>颜色</label>
            <div class="color-picker">
              <button
                v-for="c in colors"
                :key="c"
                class="color-swatch"
                :class="{ selected: form.color === c }"
                :style="{ background: c }"
                @click="form.color = c"
              />
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn" @click="closeModal">取消</button>
          <button class="btn btn-primary" @click="save" :disabled="!form.name">保存</button>
        </div>
      </div>
    </div>

    <!-- 删除确认 -->
    <div class="modal-overlay" v-if="showDeleteConfirm" @click.self="showDeleteConfirm = false">
      <div class="modal modal-sm">
        <div class="modal-header"><h2>确认删除</h2></div>
        <div class="modal-body">
          <p>确定要删除标签「{{ deleteTarget?.name }}」吗？关联的文件标签也会被移除。</p>
        </div>
        <div class="modal-footer">
          <button class="btn" @click="showDeleteConfirm = false">取消</button>
          <button class="btn btn-danger" @click="doDelete">删除</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import type { TagItem } from '@/features/tags/types'
import { fetchTags, createTag, updateTag, deleteTag } from '@/features/tags/api'

const colors = [
  '#ef4444', '#f97316', '#eab308', '#22c55e', '#14b8a6',
  '#3b82f6', '#6366f1', '#a855f7', '#ec4899', '#6b7280',
]

const tags = ref<TagItem[]>([])
const showModal = ref(false)
const editing = ref<TagItem | null>(null)
const form = reactive({ name: '', color: colors[5] }) // 默认蓝色

const showDeleteConfirm = ref(false)
const deleteTarget = ref<TagItem | null>(null)

onMounted(loadTags)

async function loadTags() {
  try {
    tags.value = await fetchTags()
  } catch (e) {
    console.error('加载标签失败', e)
  }
}

function openCreate() {
  editing.value = null
  form.name = ''
  form.color = colors[5]
  showModal.value = true
}

function openEdit(t: TagItem) {
  editing.value = t
  form.name = t.name
  form.color = t.color
  showModal.value = true
}

function closeModal() {
  showModal.value = false
}

async function save() {
  if (!form.name) return
  try {
    if (editing.value) {
      await updateTag(editing.value.id, { name: form.name, color: form.color })
    } else {
      await createTag({ name: form.name, color: form.color })
    }
    showModal.value = false
    await loadTags()
  } catch (e) {
    console.error('保存标签失败', e)
  }
}

function confirmDelete(t: TagItem) {
  deleteTarget.value = t
  showDeleteConfirm.value = true
}

async function doDelete() {
  if (!deleteTarget.value) return
  try {
    await deleteTag(deleteTarget.value.id)
    showDeleteConfirm.value = false
    await loadTags()
  } catch (e) {
    console.error('删除标签失败', e)
  }
}
</script>

<style scoped>
.tags-page { max-width: 800px; }

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}
.page-title { font-size: 22px; font-weight: 600; }

.tag-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 12px;
}

.tag-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-left: 4px solid;
  border-radius: 8px;
}

.tag-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  flex-shrink: 0;
}

.tag-name {
  flex: 1;
  font-size: 15px;
  font-weight: 500;
}

.tag-actions {
  display: flex;
  gap: 4px;
}

.color-picker {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.color-swatch {
  width: 28px;
  height: 28px;
  border-radius: 6px;
  border: 2px solid transparent;
  cursor: pointer;
  transition: all .15s;
}
.color-swatch:hover {
  transform: scale(1.15);
}
.color-swatch.selected {
  border-color: var(--text-primary);
  box-shadow: 0 0 0 2px var(--bg-primary);
}

/* 复用公共样式 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 60px 0;
  color: var(--text-muted);
}
.empty-icon { font-size: 48px; margin-bottom: 12px; }
.empty-text { font-size: 15px; }
</style>