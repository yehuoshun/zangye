<!--
  FilesPage.vue — 文件管理页面（文件 + 文件夹浏览）

  支持：
    - 文件夹网格展示 + 点击进入子目录
    - 面包屑导航 + 点击跳回
    - 新建文件夹 / 新建文件
    - 文件表格 + 编辑/删除
-->

<template>
  <div class="files-page">
    <!-- 面包屑导航 -->
    <div class="breadcrumb">
      <span class="breadcrumb-item" @click="goRoot">📂 根目录</span>
      <template v-for="(seg, i) in breadcrumbs" :key="i">
        <span class="breadcrumb-sep">›</span>
        <span
          class="breadcrumb-item"
          :class="{ active: i === breadcrumbs.length - 1 }"
          @click="goBreadcrumb(i)"
        >{{ seg.name }}</span>
      </template>
    </div>

    <!-- 操作栏 -->
    <div class="page-header">
      <h1 class="page-title">{{ currentFolder ? currentFolder.name : '文件管理' }}</h1>
      <div class="header-actions">
        <button class="btn" @click="openCreateFolder">+ 新建文件夹</button>
        <button class="btn btn-primary" @click="openCreateFile">+ 新建文件</button>
      </div>
    </div>

    <!-- 文件夹网格 -->
    <div class="folder-section" v-if="folders.length">
      <div class="folder-grid">
        <div class="folder-card" v-for="f in folders" :key="f.id" @dblclick="enterFolder(f)">
          <span class="folder-icon">{{ f.icon }}</span>
          <span class="folder-name">{{ f.name }}</span>
          <div class="folder-actions">
            <button class="btn-icon" title="重命名" @click.stop="openEditFolder(f)">✏️</button>
            <button class="btn-icon" title="删除" @click.stop="confirmDeleteFolder(f)">🗑️</button>
          </div>
        </div>
      </div>
    </div>

    <!-- 文件表格 -->
    <div class="table-card" v-if="files.length">
      <table class="data-table">
        <thead>
          <tr>
            <th>文件</th>
            <th>大小</th>
            <th>类型</th>
            <th>创建时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="f in files" :key="f.id">
            <td class="cell-file">
              <div class="file-name">{{ f.display_name || filenameFromPath(f.path) }}</div>
              <div class="file-path-breadcrumb">
                <span class="crumb-seg" v-for="(seg, i) in pathSegments(f.path)" :key="i">
                  <span v-if="i > 0" class="crumb-sep">›</span>
                  <span class="crumb-text">{{ seg }}</span>
                </span>
              </div>
            </td>
            <td>{{ formatFileSize(f.file_size) }}</td>
            <td>
              <span class="mime-badge" v-if="f.mime_type">{{ f.mime_type }}</span>
              <span class="text-muted" v-else>—</span>
            </td>
            <td class="cell-date">{{ formatDate(f.created_at) }}</td>
            <td class="cell-actions">
              <button class="btn btn-sm" @click="openEditFile(f)">编辑</button>
              <button class="btn btn-sm btn-danger" @click="confirmDeleteFile(f)">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 空状态 -->
    <div class="empty-state" v-if="!folders.length && !files.length">
      <span class="empty-icon">📂</span>
      <span class="empty-text">此目录为空，点击上方按钮新建文件夹或文件</span>
    </div>

    <!-- 文件夹弹窗（新建/编辑） -->
    <div class="modal-overlay" v-if="showFolderModal" @click.self="closeFolderModal">
      <div class="modal modal-sm">
        <div class="modal-header">
          <h2>{{ editingFolder ? '重命名文件夹' : '新建文件夹' }}</h2>
          <button class="btn-close" @click="closeFolderModal">✕</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>文件夹名称 <span class="required">*</span></label>
            <input v-model="folderForm.name" class="form-input" placeholder="输入文件夹名称" @keyup.enter="saveFolder" />
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn" @click="closeFolderModal">取消</button>
          <button class="btn btn-primary" @click="saveFolder" :disabled="!folderForm.name">
            保存
          </button>
        </div>
      </div>
    </div>

    <!-- 文件弹窗（新建/编辑） -->
    <div class="modal-overlay" v-if="showFileModal" @click.self="closeFileModal">
      <div class="modal">
        <div class="modal-header">
          <h2>{{ editingFile ? '编辑文件' : '新建文件' }}</h2>
          <button class="btn-close" @click="closeFileModal">✕</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>文件路径 <span class="required">*</span></label>
            <input v-model="fileForm.path" class="form-input" placeholder="E:\Game\YGO\MDPro3\Expansions\1.cdb" />
          </div>
          <div class="form-group">
            <label>显示名称</label>
            <input v-model="fileForm.display_name" class="form-input" placeholder="可选" />
          </div>
          <div class="form-row">
            <div class="form-group">
              <label>文件大小（字节）</label>
              <input v-model.number="fileForm.file_size" type="number" class="form-input" placeholder="0" />
            </div>
            <div class="form-group">
              <label>MIME 类型</label>
              <input v-model="fileForm.mime_type" class="form-input" placeholder="如 image/png" />
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn" @click="closeFileModal">取消</button>
          <button class="btn btn-primary" @click="saveFile" :disabled="!fileForm.path">
            保存
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
          <p>{{ deleteMessage }}</p>
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
import type { FileItem, FileCreateRequest, FileUpdateRequest } from '@/features/files/types'
import type { FolderItem, FolderCreateRequest, FolderUpdateRequest } from '@/features/folders/types'
import { fetchFiles, createFile, updateFile, deleteFile } from '@/features/files/api'
import { fetchFolders, createFolder, updateFolder, deleteFolder } from '@/features/folders/api'

// 导航状态
const currentFolderId = ref<string | null>(null)
const currentFolder = ref<FolderItem | null>(null)
const breadcrumbs = ref<{ id: string; name: string }[]>([])

// 数据
const folders = ref<FolderItem[]>([])
const files = ref<FileItem[]>([])

// 文件夹弹窗
const showFolderModal = ref(false)
const editingFolder = ref<FolderItem | null>(null)
const folderForm = reactive({ name: '' })

// 文件弹窗
const showFileModal = ref(false)
const editingFile = ref<FileItem | null>(null)
const fileForm = reactive<FileCreateRequest>({
  folder_id: '',
  path: '',
  display_name: null,
  file_size: 0,
  mime_type: null,
})

// 删除确认
const showDeleteConfirm = ref(false)
const deleteMessage = ref('')
const deleteType = ref<'file' | 'folder'>('file')
const deleteTargetId = ref('')

onMounted(loadData)

async function loadData() {
  await Promise.all([loadFolders(), loadFiles()])
}

async function loadFolders() {
  try {
    folders.value = await fetchFolders(currentFolderId.value || undefined)
  } catch (e) {
    console.error('加载文件夹失败', e)
  }
}

async function loadFiles() {
  try {
    // TODO: 等 files API 支持 folder_id 过滤后改这里
    files.value = await fetchFiles()
  } catch (e) {
    console.error('加载文件失败', e)
  }
}

// ===== 文件夹导航 =====

function goRoot() {
  currentFolderId.value = null
  currentFolder.value = null
  breadcrumbs.value = []
  loadData()
}

function enterFolder(f: FolderItem) {
  currentFolderId.value = f.id
  currentFolder.value = f
  breadcrumbs.value.push({ id: f.id, name: f.name })
  loadData()
}

function goBreadcrumb(i: number) {
  if (i === breadcrumbs.value.length - 1) return
  breadcrumbs.value = breadcrumbs.value.slice(0, i + 1)
  const last = breadcrumbs.value[i]
  currentFolderId.value = last.id
  currentFolder.value = { id: last.id, name: last.name } as FolderItem
  loadData()
}

// ===== 文件夹操作 =====

function openCreateFolder() {
  editingFolder.value = null
  folderForm.name = ''
  showFolderModal.value = true
}

function openEditFolder(f: FolderItem) {
  editingFolder.value = f
  folderForm.name = f.name
  showFolderModal.value = true
}

function closeFolderModal() {
  showFolderModal.value = false
}

async function saveFolder() {
  if (!folderForm.name) return
  try {
    if (editingFolder.value) {
      await updateFolder(editingFolder.value.id, { name: folderForm.name })
    } else {
      await createFolder({
        name: folderForm.name,
        parent_id: currentFolderId.value || null,
      })
    }
    showFolderModal.value = false
    await loadFolders()
  } catch (e) {
    console.error('保存文件夹失败', e)
  }
}

function confirmDeleteFolder(f: FolderItem) {
  deleteType.value = 'folder'
  deleteTargetId.value = f.id
  deleteMessage.value = `确定要删除文件夹「${f.name}」吗？其中的文件和子文件夹也会被删除，此操作不可撤销。`
  showDeleteConfirm.value = true
}

// ===== 文件操作 =====

function openCreateFile() {
  editingFile.value = null
  fileForm.folder_id = currentFolderId.value || ''
  fileForm.path = ''
  fileForm.display_name = null
  fileForm.file_size = 0
  fileForm.mime_type = null
  showFileModal.value = true
}

function openEditFile(f: FileItem) {
  editingFile.value = f
  fileForm.folder_id = f.folder_id
  fileForm.path = f.path
  fileForm.display_name = f.display_name
  fileForm.file_size = f.file_size
  fileForm.mime_type = f.mime_type
  showFileModal.value = true
}

function closeFileModal() {
  showFileModal.value = false
}

async function saveFile() {
  if (!fileForm.path) return
  try {
    if (editingFile.value) {
      const data: FileUpdateRequest = {
        folder_id: fileForm.folder_id || null,
        path: fileForm.path || null,
        display_name: fileForm.display_name || null,
        file_size: fileForm.file_size || null,
        mime_type: fileForm.mime_type || null,
      }
      await updateFile(editingFile.value.id, data)
    } else {
      await createFile(fileForm)
    }
    showFileModal.value = false
    await loadFiles()
  } catch (e) {
    console.error('保存文件失败', e)
  }
}

function confirmDeleteFile(f: FileItem) {
  deleteType.value = 'file'
  deleteTargetId.value = f.id
  deleteMessage.value = `确定要删除文件「${f.display_name || f.path}」吗？此操作不可撤销。`
  showDeleteConfirm.value = true
}

// ===== 删除执行 =====

async function doDelete() {
  try {
    if (deleteType.value === 'folder') {
      await deleteFolder(deleteTargetId.value)
      showDeleteConfirm.value = false
      await loadFolders()
    } else {
      await deleteFile(deleteTargetId.value)
      showDeleteConfirm.value = false
      await loadFiles()
    }
  } catch (e) {
    console.error('删除失败', e)
  }
}

// ===== 工具函数 =====

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

function formatDate(ts: string): string {
  const d = new Date(ts)
  return d.toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
}

function filenameFromPath(p: string): string {
  const normalized = p.replace(/\\/g, '/')
  return normalized.split('/').pop() || p
}

function pathSegments(p: string): string[] {
  const normalized = p.replace(/\\/g, '/')
  const parts = normalized.split('/')
  parts.pop()
  return parts.filter(Boolean)
}
</script>

<style scoped>
.files-page { max-width: 1200px; }

/* 面包屑 */
.breadcrumb {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 12px 16px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  margin-bottom: 20px;
  font-size: 14px;
  flex-wrap: wrap;
}
.breadcrumb-item {
  color: var(--text-secondary);
  cursor: pointer;
  padding: 2px 6px;
  border-radius: 4px;
  transition: all .15s;
}
.breadcrumb-item:hover { color: var(--accent); background: var(--bg-tertiary); }
.breadcrumb-item.active { color: var(--text-primary); font-weight: 500; cursor: default; }
.breadcrumb-item.active:hover { background: transparent; }
.breadcrumb-sep { color: var(--text-muted); user-select: none; }

/* 页面头部 */
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}
.page-title { font-size: 24px; font-weight: 600; }
.header-actions { display: flex; gap: 8px; }

/* 文件夹网格 */
.folder-section { margin-bottom: 20px; }
.folder-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 12px;
}
.folder-card {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  cursor: pointer;
  transition: all .15s;
  position: relative;
}
.folder-card:hover { border-color: var(--accent); background: var(--bg-tertiary); }
.folder-icon { font-size: 24px; flex-shrink: 0; }
.folder-name {
  flex: 1;
  font-size: 14px;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.folder-actions {
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity .15s;
}
.folder-card:hover .folder-actions { opacity: 1; }
.btn-icon {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 14px;
  padding: 2px;
  border-radius: 4px;
}
.btn-icon:hover { background: var(--bg-primary); }

/* 表格卡片 */
.table-card {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  overflow: hidden;
}
.data-table { width: 100%; border-collapse: collapse; }
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

.cell-file { max-width: 400px; }
.file-name { font-weight: 500; margin-bottom: 4px; }
.file-path-breadcrumb {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  font-size: 12px;
  color: var(--text-muted);
}
.crumb-sep { margin: 0 4px; color: var(--text-muted); }
.crumb-text { white-space: nowrap; }

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
.btn-primary { background: var(--accent); color: var(--accent-text); border-color: var(--accent); }
.btn-primary:hover { opacity: .9; }
.btn-danger { color: #ef4444; border-color: #ef4444; }
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
.form-group { display: flex; flex-direction: column; gap: 6px; margin-bottom: 16px; }
.form-group label { font-size: 13px; font-weight: 500; color: var(--text-secondary); }
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
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.text-muted { color: var(--text-muted); }
</style>