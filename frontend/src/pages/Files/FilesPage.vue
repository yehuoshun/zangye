<template>
  <div class="files-page">
    <div class="page-header">
      <h2 class="page-title">{{ isTrash ? '回收站' : '文件' }}</h2>
      <button v-if="!isTrash" class="btn btn-primary" @click="showCreateModal = true">+ 新建文件</button>
    </div>

    <!-- 视图切换 -->
    <div v-if="viewMode === 'list'" class="list-view">
      <table class="table">
        <thead>
          <tr>
            <th style="width: 30px;">
              <input type="checkbox" :checked="allSelected" @change="toggleAll" />
            </th>
            <th @click="toggleSort('name')">名称 {{ sortIcon('name') }}</th>
            <th @click="toggleSort('file_type')">类型 {{ sortIcon('file_type') }}</th>
            <th @click="toggleSort('file_size')">大小 {{ sortIcon('file_size') }}</th>
            <th @click="toggleSort('created_at')">创建时间 {{ sortIcon('created_at') }}</th>
            <th>标签</th>
            <th style="width: 180px;">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="file in files" :key="file.id"
            @contextmenu.prevent="showContextMenu($event, file)"
          >
            <td><input type="checkbox" :checked="selectedIds.has(file.id)" @change="toggleSelect(file.id)" /></td>
            <td class="file-name-cell">
              <span class="file-icon">{{ getFileIcon(file.file_type) }}</span>
              {{ file.name }}
            </td>
            <td><span class="file-type-badge">{{ file.file_type }}</span></td>
            <td>{{ formatSize(file.file_size) }}</td>
            <td>{{ formatDate(file.created_at) }}</td>
            <td>
              <span v-for="tag in file.tags" :key="tag.id"
                class="tag"
                :style="{ background: tag.color + '22', color: tag.color, border: '1px solid ' + tag.color + '44' }"
              >{{ tag.name }}</span>
              <span v-if="!file.tags || file.tags.length === 0" class="no-tags">无</span>
            </td>
            <td class="actions">
              <button class="btn btn-sm btn-secondary" @click="previewFile(file)">预览</button>
              <button class="btn btn-sm btn-secondary" @click="editFile(file)">编辑</button>
              <button class="btn btn-sm btn-secondary" @click="manageTags(file)">标签</button>
              <button v-if="isTrash" class="btn btn-sm btn-secondary" @click="restoreFile(file.id)">恢复</button>
              <button v-if="isTrash" class="btn btn-sm btn-danger" @click="confirmHardDelete(file)">彻底删除</button>
              <button v-else class="btn btn-sm btn-danger" @click="deleteFile(file.id)">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 网格视图 -->
    <div v-else class="grid-view">
      <div v-for="file in files" :key="file.id" class="grid-card"
        @contextmenu.prevent="showContextMenu($event, file)"
      >
        <div class="grid-icon">{{ getFileIcon(file.file_type) }}</div>
        <div class="grid-name" :title="file.name">{{ file.name }}</div>
        <div class="grid-meta">
          <span class="file-type-badge">{{ file.file_type }}</span>
          <span>{{ formatSize(file.file_size) }}</span>
        </div>
        <div class="grid-tags">
          <span v-for="tag in (file.tags || [])" :key="tag.id"
            class="tag"
            :style="{ background: tag.color + '22', color: tag.color }"
          >{{ tag.name }}</span>
        </div>
        <div class="grid-actions">
          <button class="btn btn-sm btn-secondary" @click="previewFile(file)">预览</button>
          <button class="btn btn-sm btn-secondary" @click="editFile(file)">编辑</button>
          <button v-if="!isTrash" class="btn btn-sm btn-danger" @click="deleteFile(file.id)">删除</button>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-if="files.length === 0 && !loading" class="empty-state">
      <div class="empty-state-icon">{{ isTrash ? '🗑️' : '📁' }}</div>
      <div class="empty-state-text">{{ isTrash ? '回收站为空' : '暂无文件，点击上方按钮创建' }}</div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading">加载中...</div>

    <!-- 分页 -->
    <div v-if="total > pageSize" class="pagination">
      <button :disabled="currentPage <= 1" @click="goPage(currentPage - 1)">上一页</button>
      <button v-for="p in pageNumbers" :key="p"
        :class="{ active: p === currentPage }"
        @click="goPage(p)"
      >{{ p }}</button>
      <button :disabled="currentPage >= totalPages" @click="goPage(currentPage + 1)">下一页</button>
      <span class="pagination-info">共 {{ total }} 条</span>
    </div>

    <!-- 右键菜单 -->
    <div v-if="contextMenu.visible" class="context-menu"
      :style="{ left: contextMenu.x + 'px', top: contextMenu.y + 'px' }"
    >
      <div class="context-menu-item" @click="previewFile(contextMenu.file!)">👁️ 预览</div>
      <div class="context-menu-item" @click="editFile(contextMenu.file!)">✏️ 编辑</div>
      <div class="context-menu-item" @click="manageTags(contextMenu.file!)">🏷️ 管理标签</div>
      <div class="context-menu-item" @click="openWith(contextMenu.file!)">🔗 打开方式</div>
      <div class="context-menu-item danger" @click="deleteFile(contextMenu.file!.id)">🗑️ 删除</div>
    </div>

    <!-- 预览弹窗 -->
    <div v-if="previewFileData" class="modal-overlay" @click.self="previewFileData = null">
      <div class="modal-content preview-modal">
        <div class="modal-header">
          <h3 class="modal-title">{{ previewFileData.name }}</h3>
          <button class="modal-close" @click="previewFileData = null">&times;</button>
        </div>
        <div class="preview-content">
          <img v-if="isImage(previewFileData.file_type)"
            :src="`/api/files/${previewFileData.id}/content`"
            class="preview-image"
          />
          <video v-else-if="isVideo(previewFileData.file_type)"
            :src="`/api/files/${previewFileData.id}/content`"
            controls class="preview-video"
          ></video>
          <audio v-else-if="isAudio(previewFileData.file_type)"
            :src="`/api/files/${previewFileData.id}/content`"
            controls class="preview-audio"
          ></audio>
          <pre v-else-if="isText(previewFileData.file_type)" class="preview-text">
            <iframe :src="`/api/files/${previewFileData.id}/content`" class="preview-iframe"></iframe>
          </pre>
          <div v-else class="preview-unsupported">
            该类型暂不支持预览
          </div>
        </div>
      </div>
    </div>

    <!-- 编辑/新建弹窗 -->
    <div v-if="showCreateModal || editingFile" class="modal-overlay" @click.self="closeEditModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3 class="modal-title">{{ editingFile ? '编辑文件' : '新建文件' }}</h3>
          <button class="modal-close" @click="closeEditModal">&times;</button>
        </div>
        <div class="form">
          <div class="form-group">
            <label>文件名</label>
            <input v-model="form.name" class="input" placeholder="文件名（含扩展名）" />
          </div>
          <div class="form-group">
            <label>文件路径 (JSON 数组)</label>
            <textarea v-model="form.paths" class="input" rows="3" placeholder='["C:\\backup\\file.jpg"]'></textarea>
          </div>
          <div class="form-group">
            <label>文件大小 (字节)</label>
            <input v-model.number="form.file_size" class="input" type="number" placeholder="0" />
          </div>
          <div class="form-group">
            <label>描述</label>
            <textarea v-model="form.description" class="input" rows="2" placeholder="可选描述"></textarea>
          </div>
          <div style="display: flex; gap: 8px; justify-content: flex-end; margin-top: 16px;">
            <button class="btn btn-secondary" @click="closeEditModal">取消</button>
            <button class="btn btn-primary" @click="saveFile">{{ editingFile ? '保存' : '创建' }}</button>
          </div>
        </div>
      </div>
    </div>

    <!-- 标签管理弹窗 -->
    <div v-if="tagManagerFile" class="modal-overlay" @click.self="tagManagerFile = null">
      <div class="modal-content">
        <div class="modal-header">
          <h3 class="modal-title">管理标签 - {{ tagManagerFile.name }}</h3>
          <button class="modal-close" @click="tagManagerFile = null">&times;</button>
        </div>
        <div class="tag-manager">
          <div class="tag-list">
            <div v-for="tag in allTags" :key="tag.id" class="tag-item"
              :class="{ selected: fileTagIds.has(tag.id) }"
              @click="toggleFileTag(tag.id)"
            >
              <span class="tag-dot" :style="{ background: tag.color }"></span>
              {{ tag.name }}
              <span v-if="fileTagIds.has(tag.id)" class="tag-check">✓</span>
            </div>
          </div>
          <div class="tag-add" style="margin-top: 12px;">
            <div style="display: flex; gap: 8px;">
              <input v-model="newTagName" class="input" placeholder="新建标签..." @keyup.enter="addNewTag" />
              <button class="btn btn-primary btn-sm" @click="addNewTag">添加</button>
            </div>
          </div>
          <div style="display: flex; gap: 8px; justify-content: flex-end; margin-top: 16px;">
            <button class="btn btn-primary" @click="saveTags">保存</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getFiles, createFile, updateFile, deleteFile as apiDeleteFile,
  restoreFile as apiRestoreFile, hardDeleteFile as apiHardDeleteFile,
  getFileTags, setFileTags } from '@/features/files/api'
import { getTags, createTag } from '@/features/tags/api'
import type { FileItem } from '@/features/files/types'
import type { TagItem } from '@/features/tags/types'

const route = useRoute()

const props = defineProps<{
  viewMode?: string
  searchKeyword?: string
  selectedFolderId?: string
}>()

const files = ref<FileItem[]>([])
const loading = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(50)
const sortField = ref('created_at')
const sortDir = ref('desc')
const selectedIds = ref(new Set<string>())

// 弹窗状态
const showCreateModal = ref(false)
const editingFile = ref<FileItem | null>(null)
const previewFileData = ref<FileItem | null>(null)
const tagManagerFile = ref<FileItem | null>(null)
const allTags = ref<TagItem[]>([])
const fileTagIds = ref(new Set<string>())
const newTagName = ref('')

// 右键菜单
const contextMenu = ref<{ visible: boolean; x: number; y: number; file: FileItem | null }>({
  visible: false,
  x: 0,
  y: 0,
  file: null,
})

// 表单
const form = ref({
  name: '',
  paths: '',
  file_size: 0,
  description: '',
})

const isTrash = computed(() => route.path === '/trash')

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

const allSelected = computed(() => {
  return files.value.length > 0 && files.value.every(f => selectedIds.value.has(f.id))
})

const pageNumbers = computed(() => {
  const pages: number[] = []
  const start = Math.max(1, currentPage.value - 2)
  const end = Math.min(totalPages.value, currentPage.value + 2)
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
})

// 监听路由变化
watch(() => route.path, () => {
  loadFiles()
})

// 监听筛选条件
watch([() => props.searchKeyword, () => props.selectedFolderId], () => {
  currentPage.value = 1
  loadFiles()
})

onMounted(() => {
  loadFiles()
  loadTags()
  // 点击其他地方关闭右键菜单
  document.addEventListener('click', () => {
    contextMenu.value.visible = false
  })
})

async function loadFiles() {
  loading.value = true
  try {
    const result = await getFiles({
      folder_id: props.selectedFolderId || undefined,
      keyword: props.searchKeyword || undefined,
      order_by: sortField.value,
      order_dir: sortDir.value,
      page: currentPage.value,
      page_size: pageSize.value,
      trash: isTrash.value,
    })
    files.value = result.data
    total.value = result.total
  } catch (e) {
    console.error('加载文件列表失败:', e)
  } finally {
    loading.value = false
  }
}

async function loadTags() {
  try {
    allTags.value = await getTags()
  } catch (e) {
    console.error('加载标签列表失败:', e)
  }
}

function toggleSort(field: string) {
  if (sortField.value === field) {
    sortDir.value = sortDir.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortField.value = field
    sortDir.value = 'asc'
  }
  loadFiles()
}

function sortIcon(field: string) {
  if (sortField.value !== field) return '↕'
  return sortDir.value === 'asc' ? '↑' : '↓'
}

function toggleAll() {
  if (allSelected.value) {
    selectedIds.value.clear()
  } else {
    files.value.forEach(f => selectedIds.value.add(f.id))
  }
}

function toggleSelect(id: string) {
  if (selectedIds.value.has(id)) {
    selectedIds.value.delete(id)
  } else {
    selectedIds.value.add(id)
  }
}

function goPage(page: number) {
  currentPage.value = page
  loadFiles()
}

function formatSize(size: number): string {
  if (size < 1024) return size + ' B'
  if (size < 1024 * 1024) return (size / 1024).toFixed(1) + ' KB'
  if (size < 1024 * 1024 * 1024) return (size / (1024 * 1024)).toFixed(1) + ' MB'
  return (size / (1024 * 1024 * 1024)).toFixed(2) + ' GB'
}

function formatDate(dateStr: string): string {
  const d = new Date(dateStr)
  return d.toLocaleDateString('zh-CN') + ' ' + d.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}

function getFileIcon(type: string): string {
  const t = type.toLowerCase()
  if (['jpg', 'jpeg', 'png', 'gif', 'webp', 'bmp', 'svg'].includes(t)) return '🖼️'
  if (['mp4', 'avi', 'mov', 'mkv', 'webm'].includes(t)) return '🎬'
  if (['mp3', 'wav', 'ogg', 'flac', 'aac'].includes(t)) return '🎵'
  if (['pdf'].includes(t)) return '📕'
  if (['zip', 'rar', '7z', 'tar', 'gz'].includes(t)) return '📦'
  if (['txt', 'md'].includes(t)) return '📝'
  if (['doc', 'docx'].includes(t)) return '📄'
  if (['xls', 'xlsx'].includes(t)) return '📊'
  return '📄'
}

function isImage(type: string) { return ['jpg', 'jpeg', 'png', 'gif', 'webp', 'bmp', 'svg'].includes(type) }
function isVideo(type: string) { return ['mp4', 'avi', 'mov', 'mkv', 'webm'].includes(type) }
function isAudio(type: string) { return ['mp3', 'wav', 'ogg', 'flac', 'aac'].includes(type) }
function isText(type: string) { return ['txt', 'md', 'json', 'xml', 'html', 'css', 'js', 'ts', 'go', 'py', 'java', 'sql', 'yaml', 'yml', 'log'].includes(type) }

function previewFile(file: FileItem) {
  previewFileData.value = file
  contextMenu.value.visible = false
}

function editFile(file: FileItem) {
  editingFile.value = file
  form.value = {
    name: file.name,
    paths: file.paths || '',
    file_size: file.file_size,
    description: file.description || '',
  }
  contextMenu.value.visible = false
}

function closeEditModal() {
  showCreateModal.value = false
  editingFile.value = null
  form.value = { name: '', paths: '', file_size: 0, description: '' }
}

async function saveFile() {
  if (!form.value.name.trim()) return
  try {
    if (editingFile.value) {
      await updateFile(editingFile.value.id, form.value)
    } else {
      await createFile(form.value)
    }
    closeEditModal()
    loadFiles()
  } catch (e) {
    console.error('保存文件失败:', e)
  }
}

async function deleteFile(id: string) {
  if (!confirm('确定删除此文件？')) return
  try {
    await apiDeleteFile(id)
    loadFiles()
  } catch (e) {
    console.error('删除文件失败:', e)
  }
  contextMenu.value.visible = false
}

async function restoreFile(id: string) {
  try {
    await apiRestoreFile(id)
    loadFiles()
  } catch (e) {
    console.error('恢复文件失败:', e)
  }
}

async function confirmHardDelete(file: FileItem) {
  if (!confirm(`确定彻底删除"${file.name}"？此操作不可恢复！`)) return
  try {
    await apiHardDeleteFile(file.id)
    loadFiles()
  } catch (e) {
    console.error('彻底删除文件失败:', e)
  }
}

function showContextMenu(event: MouseEvent, file: FileItem) {
  contextMenu.value = {
    visible: true,
    x: event.clientX,
    y: event.clientY,
    file,
  }
}

function openWith(file: FileItem) {
  // 打开方式 - 使用系统默认
  contextMenu.value.visible = false
  // 这里可以扩展为选择外部程序
  alert(`打开方式：${file.name}\n（使用系统默认程序打开）`)
}

async function manageTags(file: FileItem) {
  tagManagerFile.value = file
  fileTagIds.value = new Set()
  try {
    const tags = await getFileTags(file.id)
    tags.forEach((t: any) => fileTagIds.value.add(t.id))
  } catch (e) {
    console.error('加载文件标签失败:', e)
  }
  contextMenu.value.visible = false
}

function toggleFileTag(tagId: string) {
  if (fileTagIds.value.has(tagId)) {
    fileTagIds.value.delete(tagId)
  } else {
    fileTagIds.value.add(tagId)
  }
}

async function addNewTag() {
  if (!newTagName.value.trim()) return
  try {
    const tag = await createTag(newTagName.value.trim())
    newTagName.value = ''
    await loadTags()
    fileTagIds.value.add(tag.id)
  } catch (e) {
    console.error('创建标签失败:', e)
  }
}

async function saveTags() {
  if (!tagManagerFile.value) return
  try {
    await setFileTags(tagManagerFile.value.id, Array.from(fileTagIds.value))
    tagManagerFile.value = null
    loadFiles()
  } catch (e) {
    console.error('保存标签失败:', e)
  }
}
</script>

<style scoped>
.files-page {
  max-width: 100%;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
}

.file-name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.file-icon {
  font-size: 16px;
}

.file-type-badge {
  font-size: 11px;
  padding: 2px 6px;
  border-radius: 4px;
  background: var(--bg-tertiary);
  color: var(--text-secondary);
  text-transform: uppercase;
}

.no-tags {
  color: var(--text-muted);
  font-size: 12px;
}

.actions {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

/* 网格视图 */
.grid-view {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 12px;
}

.grid-card {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius);
  padding: 16px;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.grid-card:hover {
  border-color: var(--accent-primary);
  box-shadow: var(--shadow-sm);
}

.grid-icon {
  font-size: 36px;
  text-align: center;
  margin-bottom: 8px;
}

.grid-name {
  font-size: 13px;
  font-weight: 500;
  text-align: center;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-bottom: 4px;
}

.grid-meta {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 11px;
  color: var(--text-muted);
  margin-bottom: 8px;
}

.grid-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  justify-content: center;
  margin-bottom: 8px;
}

.grid-actions {
  display: flex;
  gap: 4px;
  justify-content: center;
}

/* 预览 */
.preview-modal {
  width: 80vw;
  max-width: 900px;
}

.preview-content {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 200px;
}

.preview-image {
  max-width: 100%;
  max-height: 70vh;
  object-fit: contain;
}

.preview-video {
  max-width: 100%;
  max-height: 70vh;
}

.preview-audio {
  width: 100%;
}

.preview-iframe {
  width: 100%;
  height: 60vh;
  border: none;
  background: #fff;
}

.preview-unsupported {
  color: var(--text-muted);
  padding: 40px;
}

/* 表单 */
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

.form-group textarea {
  resize: vertical;
  font-family: var(--font-mono);
  font-size: 12px;
}

/* 标签管理 */
.tag-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border: 1px solid var(--border-color);
  border-radius: 16px;
  cursor: pointer;
  font-size: 13px;
  transition: all var(--transition-fast);
}

.tag-item:hover {
  border-color: var(--accent-primary);
}

.tag-item.selected {
  border-color: var(--accent-primary);
  background: var(--accent-primary) + '15';
}

.tag-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.tag-check {
  color: var(--accent-success);
  font-weight: bold;
}
</style>
