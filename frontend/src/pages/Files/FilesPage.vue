<!--
  FilesPage.vue — 文件管理页面（文件 + 文件夹浏览）

  支持：
    - 文件夹网格展示 + 点击进入子目录
    - 面包屑导航 + 点击跳回
    - 新建文件夹 / 新建文件
    - 文件表格 + 编辑/删除
-->

<template>
  <div class="files-page" @click="closeContextMenu" @contextmenu.prevent>
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
        <div class="view-toggle">
          <button class="btn btn-sm" :class="{ active: viewMode === 'grid' }" @click="viewMode = 'grid'">▦ 网格</button>
          <button class="btn btn-sm" :class="{ active: viewMode === 'list' }" @click="viewMode = 'list'">☰ 列表</button>
        </div>
        <button class="btn" @click="openCreateFolder">+ 新建文件夹</button>
        <button class="btn btn-primary" @click="openCreateFile">+ 新建文件</button>
      </div>
    </div>

    <!-- 网格视图 -->
    <template v-if="viewMode === 'grid'">
      <div class="item-grid" v-if="folders.length || files.length">
        <div class="item-card folder-card" v-for="f in folders" :key="f.id" @dblclick="enterFolder(f)" @contextmenu.prevent="openContextMenu($event, f)">
          <span class="item-icon">{{ f.icon }}</span>
          <span class="item-name">{{ f.name }}</span>
        </div>
        <div class="item-card file-card" v-for="f in files" :key="f.id" @contextmenu.prevent="openContextMenu($event, null, f)">
          <span class="item-icon">{{ fileIcon(f.mime_type) }}</span>
          <span class="item-name">{{ f.display_name || filenameFromPath(f.path) }}</span>
          <div class="item-meta">{{ formatFileSize(f.file_size) }}</div>
        </div>
      </div>
    </template>

    <!-- 列表视图 -->
    <div class="table-card" v-if="viewMode === 'list' && (folders.length || files.length)">
      <table class="data-table">
        <thead>
          <tr>
            <th>名称</th>
            <th>大小</th>
            <th>类型</th>
            <th>创建时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="f in folders" :key="f.id" class="folder-row" @dblclick="enterFolder(f)" @contextmenu.prevent="openContextMenu($event, f)">
            <td>
              <span class="row-icon">{{ f.icon }}</span>
              <span class="row-name">{{ f.name }}</span>
            </td>
            <td class="text-muted">—</td>
            <td><span class="type-badge folder-badge">文件夹</span></td>
            <td class="cell-date">{{ formatDate(f.created_at) }}</td>
            <td class="cell-actions">
              <button class="btn btn-sm" @click.stop="openEditFolder(f)">重命名</button>
              <button class="btn btn-sm btn-danger" @click.stop="confirmDeleteFolder(f)">删除</button>
            </td>
          </tr>
          <tr v-for="f in files" :key="f.id" @contextmenu.prevent="openContextMenu($event, null, f)">
            <td>
              <span class="row-icon">{{ fileIcon(f.mime_type) }}</span>
              <span class="row-name">{{ f.display_name || filenameFromPath(f.path) }}</span>
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
            <input v-model="fileForm.path" class="form-input" placeholder="E:\Game\YGO\MDPro3\Expansions\1.cdb" @input="autoFillDisplayName" />
          </div>
          <div class="form-group">
            <label>显示名称</label>
            <input v-model="fileForm.display_name" class="form-input" placeholder="可选" />
          </div>
          <div class="form-row">
            <div class="form-group">
              <label>文件大小（字节）</label>
              <input v-model="fileForm.file_size" class="form-input" placeholder="0" @input="sanitizeFileSize" />
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

    <!-- 右键菜单 -->
    <div
      class="context-menu"
      v-if="contextMenu.visible"
      :style="{ left: contextMenu.x + 'px', top: contextMenu.y + 'px' }"
    >
      <div class="context-item" @click="contextOpen">📂 打开</div>
      <div class="context-item" v-if="contextMenu.folder" @click="contextRename">✏️ 重命名</div>
      <div class="context-item" v-if="contextMenu.file" @click="contextEdit">✏️ 编辑</div>
      <div class="context-item" v-if="contextMenu.folder" @click="contextDetail">ℹ️ 详情</div>
      <div class="context-sep"></div>
      <div class="context-item context-danger" @click="contextDelete">🗑️ 删除</div>
    </div>

    <!-- 文件夹详情弹窗 -->
    <div class="modal-overlay" v-if="showDetailModal" @click.self="showDetailModal = false">
      <div class="modal modal-sm">
        <div class="modal-header">
          <h2>文件夹详情</h2>
          <button class="btn-close" @click="showDetailModal = false">✕</button>
        </div>
        <div class="modal-body">
          <div class="detail-grid">
            <div class="detail-item">
              <span class="detail-label">名称</span>
              <span class="detail-value">{{ detailFolder?.name }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">子文件夹</span>
              <span class="detail-value">{{ detailStats.subFolders }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">包含文件</span>
              <span class="detail-value">{{ detailStats.files }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">总大小</span>
              <span class="detail-value">{{ formatFileSize(detailStats.totalSize) }}</span>
            </div>
            <div class="detail-item" style="grid-column: 1 / -1">
              <span class="detail-label">创建时间</span>
              <span class="detail-value">{{ detailFolder ? formatDate(detailFolder.created_at) : '—' }}</span>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-primary" @click="showDetailModal = false">关闭</button>
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
import { fetchFolders, createFolder, updateFolder, deleteFolder, fetchFolderStats } from '@/features/folders/api'

// 导航状态
const currentFolderId = ref<string | null>(null)
const currentFolder = ref<FolderItem | null>(null)
const breadcrumbs = ref<{ id: string; name: string }[]>([])
// 视图模式
const viewMode = ref<'grid' | 'list'>('grid')

// 右键菜单
const contextMenu = reactive({
  visible: false,
  x: 0,
  y: 0,
  folder: null as FolderItem | null,
  file: null as FileItem | null,
})

// 文件夹详情
const showDetailModal = ref(false)
const detailFolder = ref<FolderItem | null>(null)
const detailStats = reactive({ subFolders: 0, files: 0, totalSize: 0 })

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

// ===== 右键菜单 =====

function openContextMenu(e: MouseEvent, folder?: FolderItem | null, file?: FileItem | null) {
  contextMenu.visible = true
  contextMenu.x = e.clientX
  contextMenu.y = e.clientY
  contextMenu.folder = folder || null
  contextMenu.file = file || null
}

function closeContextMenu() {
  contextMenu.visible = false
}

function contextOpen() {
  if (contextMenu.folder) enterFolder(contextMenu.folder)
  closeContextMenu()
}

function contextRename() {
  if (contextMenu.folder) openEditFolder(contextMenu.folder)
  closeContextMenu()
}

function contextEdit() {
  if (contextMenu.file) openEditFile(contextMenu.file)
  closeContextMenu()
}

async function contextDetail() {
  if (!contextMenu.folder) return
  detailFolder.value = contextMenu.folder
  try {
    const stats = await fetchFolderStats(contextMenu.folder.id)
    detailStats.subFolders = stats.folder_count
    detailStats.files = stats.file_count
    detailStats.totalSize = stats.total_size
  } catch {
    detailStats.subFolders = 0
    detailStats.files = 0
    detailStats.totalSize = 0
  }
  showDetailModal.value = true
  closeContextMenu()
}

function contextDelete() {
  if (contextMenu.folder) confirmDeleteFolder(contextMenu.folder)
  else if (contextMenu.file) confirmDeleteFile(contextMenu.file)
  closeContextMenu()
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

function autoFillDisplayName() {
  fileForm.display_name = filenameFromPath(fileForm.path)
}

function sanitizeFileSize(e: Event) {
  const input = e.target as HTMLInputElement
  const raw = input.value.replace(/\s/g, '').replace(/[^\d]/g, '')
  fileForm.file_size = raw ? parseInt(raw, 10) : 0
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
  if (bytes === 0) return '0 B'
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

// 根据 MIME 类型返回文件图标
function fileIcon(mime: string | null): string {
  if (!mime) return '📄'
  if (mime.startsWith('image/')) return '🖼️'
  if (mime.startsWith('video/')) return '🎬'
  if (mime.startsWith('audio/')) return '🎵'
  if (mime.startsWith('text/')) return '📝'
  if (mime.includes('pdf')) return '📕'
  if (mime.includes('zip') || mime.includes('rar') || mime.includes('tar')) return '📦'
  return '📄'
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
.header-actions { display: flex; gap: 8px; align-items: center; }

/* 视图切换 */
.view-toggle {
  display: flex;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  overflow: hidden;
  margin-right: 8px;
}
.view-toggle .btn {
  border: none;
  border-radius: 0;
  padding: 6px 12px;
  font-size: 13px;
}
.view-toggle .btn:first-child { border-right: 1px solid var(--border-color); }
.view-toggle .btn.active { background: var(--accent); color: var(--accent-text); }

/* 网格视图 */
.item-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 12px;
}
.item-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 20px 12px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  cursor: default;
  transition: all .15s;
  position: relative;
}
.item-card:hover { border-color: var(--accent); background: var(--bg-tertiary); }
.folder-card { cursor: pointer; }
.item-icon { font-size: 36px; }
.item-name {
  font-size: 13px;
  font-weight: 500;
  text-align: center;
  word-break: break-all;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  max-width: 100%;
}
.item-meta { font-size: 11px; color: var(--text-muted); }

/* 表格 */
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
.folder-row { cursor: pointer; }
.row-icon { font-size: 20px; margin-right: 10px; vertical-align: middle; }
.row-name { font-weight: 500; }
.type-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}
.folder-badge { background: var(--bg-tertiary); color: var(--accent); }

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

/* 右键菜单 */
.context-menu {
  position: fixed;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 4px 0;
  min-width: 160px;
  z-index: 200;
  box-shadow: 0 8px 24px rgba(0, 0, 0, .3);
}
.context-item {
  padding: 8px 16px;
  font-size: 14px;
  cursor: pointer;
  transition: background .1s;
}
.context-item:hover { background: var(--bg-tertiary); }
.context-danger { color: #ef4444; }
.context-sep {
  height: 1px;
  background: var(--border-color);
  margin: 4px 0;
}

/* 详情弹窗 */
.detail-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}
.detail-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.detail-label {
  font-size: 11px;
  color: var(--text-muted);
  text-transform: uppercase;
}
.detail-value {
  font-size: 16px;
  font-weight: 600;
}

.text-muted { color: var(--text-muted); }
</style>