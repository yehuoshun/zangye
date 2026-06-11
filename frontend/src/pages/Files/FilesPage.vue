<!--
  FilesPage.vue — 文件管理页面

  功能：
    - 文件列表展示（表格形式）
    - 新建文件（弹窗表单）
    - 编辑文件（弹窗表单）
    - 删除文件（确认后删除）
    - 文件大小格式化显示
-->

<template>
  <div class="files-page">
    <!-- 页面标题 + 操作栏 -->
    <div class="page-header">
      <h1 class="page-title">文件管理</h1>
      <button class="btn btn-primary" @click="openCreate">+ 新建文件</button>
    </div>

    <!-- 文件表格 -->
    <div class="table-card">
      <table class="data-table" v-if="files.length">
        <thead>
          <tr>
            <th>文件名</th>
            <th>路径</th>
            <th>大小</th>
            <th>类型</th>
            <th>创建时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="f in files" :key="f.id">
            <td class="cell-name">{{ f.display_name || f.path.split('/').pop() || f.path }}</td>
            <td class="cell-path">{{ f.path }}</td>
            <td>{{ formatFileSize(f.file_size) }}</td>
            <td>
              <span class="mime-badge" v-if="f.mime_type">{{ f.mime_type }}</span>
              <span class="text-muted" v-else>—</span>
            </td>
            <td class="cell-date">{{ formatDate(f.created_at) }}</td>
            <td class="cell-actions">
              <button class="btn btn-sm" @click="openEdit(f)">编辑</button>
              <button class="btn btn-sm btn-danger" @click="confirmDelete(f)">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
      <div class="empty-state" v-else>
        <span class="empty-icon">📂</span>
        <span class="empty-text">暂无文件，点击右上角「新建文件」添加</span>
      </div>
    </div>

    <!-- 新建/编辑弹窗 -->
    <div class="modal-overlay" v-if="showModal" @click.self="closeModal">
      <div class="modal">
        <div class="modal-header">
          <h2>{{ isEditing ? '编辑文件' : '新建文件' }}</h2>
          <button class="btn-close" @click="closeModal">✕</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>文件路径 <span class="required">*</span></label>
            <input v-model="form.path" class="form-input" placeholder="/data/movies/example.mp4" />
            <span class="form-hint">输入服务器上的文件路径，大小和类型自动识别</span>
          </div>
          <div class="form-group">
            <label>显示名称</label>
            <input v-model="form.display_name" class="form-input" placeholder="可选，不填则使用文件名" />
          </div>
          <div class="form-group">
            <label>所属集合 ID <span class="required">*</span></label>
            <input v-model="form.collection_id" class="form-input" placeholder="集合 UUID" />
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn" @click="closeModal">取消</button>
          <button class="btn btn-primary" @click="save" :disabled="saving">
            {{ saving ? '保存中…' : '保存' }}
          </button>
        </div>
      </div>
    </div>

    <!-- 删除确认弹窗 -->
    <div class="modal-overlay" v-if="showDeleteConfirm" @click.self="showDeleteConfirm = false">
      <div class="modal modal-sm">
        <div class="modal-header">
          <h2>确认删除</h2>
        </div>
        <div class="modal-body">
          <p>确定要删除文件「{{ deleteTarget?.display_name || deleteTarget?.path }}」吗？此操作不可撤销。</p>
        </div>
        <div class="modal-footer">
          <button class="btn" @click="showDeleteConfirm = false">取消</button>
          <button class="btn btn-danger" @click="doDelete" :disabled="saving">删除</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import type { FileItem, FileCreateRequest, FileUpdateRequest } from '@/features/files/types'
import { fetchFiles, createFile, updateFile, deleteFile } from '@/features/files/api'

// 文件列表
const files = ref<FileItem[]>([])

// 弹窗状态
const showModal = ref(false)
const isEditing = ref(false)
const editingId = ref('')
const saving = ref(false)

// 表单
const form = reactive<FileCreateRequest>({
  collection_id: '',
  path: '',
  display_name: null,
})

// 删除确认
const showDeleteConfirm = ref(false)
const deleteTarget = ref<FileItem | null>(null)

// 加载文件列表
onMounted(loadFiles)

async function loadFiles() {
  try {
    files.value = await fetchFiles()
  } catch (e) {
    console.error('加载文件列表失败', e)
  }
}

// 新建
function openCreate() {
  isEditing.value = false
  editingId.value = ''
  form.collection_id = ''
  form.path = ''
  form.display_name = null
  showModal.value = true
}

// 编辑
function openEdit(f: FileItem) {
  isEditing.value = true
  editingId.value = f.id
  form.collection_id = f.collection_id
  form.path = f.path
  form.display_name = f.display_name
  showModal.value = true
}

// 关闭弹窗
function closeModal() {
  showModal.value = false
}

// 保存
async function save() {
  if (!form.path || !form.collection_id) return
  saving.value = true
  try {
    if (isEditing.value) {
      const data: FileUpdateRequest = {
        collection_id: form.collection_id || null,
        path: form.path || null,
        display_name: form.display_name || null,
      }
      await updateFile(editingId.value, data)
    } else {
      await createFile({
        collection_id: form.collection_id,
        path: form.path,
        display_name: form.display_name || null,
      })
    }
    showModal.value = false
    await loadFiles()
  } catch (e) {
    console.error('保存失败', e)
  } finally {
    saving.value = false
  }
}

// 删除确认
function confirmDelete(f: FileItem) {
  deleteTarget.value = f
  showDeleteConfirm.value = true
}

// 执行删除
async function doDelete() {
  if (!deleteTarget.value) return
  saving.value = true
  try {
    await deleteFile(deleteTarget.value.id)
    showDeleteConfirm.value = false
    deleteTarget.value = null
    await loadFiles()
  } catch (e) {
    console.error('删除失败', e)
  } finally {
    saving.value = false
  }
}

// 格式化文件大小
function formatFileSize(bytes: number): string {
  if (bytes === 0) return '—'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let i = 0
  let size = bytes
  while (size >= 1024 && i < units.length - 1) {
    size /= 1024
    i++
  }
  return size.toFixed(i === 0 ? 0 : 1) + ' ' + units[i]
}

// 格式化日期
function formatDate(ts: string): string {
  const d = new Date(ts)
  return d.toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
}
</script>

<style scoped>
.files-page { max-width: 1200px; }

/* 页面头部 */
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}
.page-title { font-size: 24px; font-weight: 600; }

/* 表格卡片 */
.table-card {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  overflow: hidden;
}

/* 表格 */
.data-table {
  width: 100%;
  border-collapse: collapse;
}
.data-table th {
  text-align: left;
  padding: 12px 16px;
  font-size: 12px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: .5px;
  background: var(--bg-tertiary);
  border-bottom: 1px solid var(--border-color);
}
.data-table td {
  padding: 12px 16px;
  font-size: 14px;
  border-bottom: 1px solid var(--border-color);
}
.data-table tr:last-child td { border-bottom: none; }
.data-table tr:hover td { background: var(--bg-tertiary); }

.cell-name { font-weight: 500; max-width: 200px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.cell-path { max-width: 300px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; color: var(--text-muted); font-family: monospace; font-size: 13px; }
.cell-date { white-space: nowrap; color: var(--text-muted); font-size: 13px; }
.cell-actions { white-space: nowrap; display: flex; gap: 8px; }

.mime-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 4px;
  background: var(--bg-tertiary);
  font-size: 12px;
  font-family: monospace;
  color: var(--text-secondary);
}

/* 空状态 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 60px 24px;
  color: var(--text-muted);
}
.empty-icon { font-size: 48px; }
.empty-text { font-size: 14px; }

/* 按钮 */
.btn {
  display: inline-flex;
  align-items: center;
  padding: 8px 16px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: var(--bg-secondary);
  color: var(--text-primary);
  font-size: 14px;
  cursor: pointer;
  transition: all .2s;
}
.btn:hover { background: var(--bg-tertiary); }
.btn:disabled { opacity: .5; cursor: not-allowed; }

.btn-primary {
  background: var(--accent);
  color: var(--accent-text);
  border-color: var(--accent);
}
.btn-primary:hover { opacity: .9; }

.btn-danger {
  color: #ef4444;
  border-color: #ef4444;
}
.btn-danger:hover { background: #ef4444; color: #fff; }

.btn-sm { padding: 4px 10px; font-size: 13px; }

.btn-close {
  background: none;
  border: none;
  color: var(--text-muted);
  font-size: 18px;
  cursor: pointer;
  padding: 4px;
}
.btn-close:hover { color: var(--text-primary); }

/* 弹窗 */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, .5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}
.modal {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  width: 520px;
  max-height: 80vh;
  overflow: auto;
}
.modal-sm { width: 400px; }

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px 0;
}
.modal-header h2 { font-size: 18px; font-weight: 600; }

.modal-body { padding: 20px 24px; }

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 0 24px 20px;
}

/* 表单 */
.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 16px;
}
.form-group label {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
}
.required { color: #ef4444; }

.form-input {
  padding: 8px 12px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 14px;
  outline: none;
}
.form-input:focus { border-color: var(--accent); }
.form-input::placeholder { color: var(--text-muted); }

.form-hint {
  font-size: 12px;
  color: var(--text-muted);
}

.text-muted { color: var(--text-muted); }
</style>