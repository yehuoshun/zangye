<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { collections, formatSize, formatTime, getIcon, type Collection, type BrowseItem } from '../api'

const route = useRoute()
const router = useRouter()
const collection = ref<Collection | null>(null)
const items = ref<BrowseItem[]>([])
const loading = ref(true)
const viewMode = ref<'grid' | 'list'>('grid')
const searchQuery = ref('')
const sortBy = ref('name')
const selectedItem = ref<BrowseItem | null>(null)

const showPathManager = ref(false)
const boundPaths = ref<any[]>([])
const newPath = ref('')
const newPathAutoScan = ref(true)

const showVFileCreator = ref(false)
const newVFile = ref({ path: '', display_name: '', size: 0 })

const breadcrumb = ref<{ id: string; name: string }[]>([])

onMounted(() => loadData())
watch(() => route.params.id, () => loadData())

async function loadData() {
  const id = route.params.id as string
  loading.value = true
  try {
    collection.value = await collections.get(id)
    items.value = await collections.browse(id, { sort: sortBy.value })
    const crumbs: { id: string; name: string }[] = []
    let current: Collection | null = collection.value
    while (current) {
      crumbs.unshift({ id: current.id, name: current.name })
      if (current.parent_id) current = await collections.get(current.parent_id)
      else break
    }
    breadcrumb.value = crumbs
  } catch (e) { console.error(e) }
  finally { loading.value = false }
}

async function refresh() { await loadData() }
function goToRoot() { router.push('/') }
function goToCollection(id: string) { router.push(`/browse/${id}`) }
function goToBreadcrumb(index: number) {
  if (index < breadcrumb.value.length - 1) router.push(`/browse/${breadcrumb.value[index].id}`)
}
function selectItem(item: BrowseItem) { selectedItem.value = item }

async function openItem(item: BrowseItem) {
  const res = await fetch('/api/open-external', {
    method: 'POST', headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ path: item.path }),
  })
  const data = await res.json()
  if (data.type === 'web' && data.url) window.open(data.url, '_blank')
  selectItem(item)
}

function previewSrc(item: BrowseItem): string | null {
  const mime = item.mime_type || ''
  if (mime.startsWith('image/')) return `/api/preview/thumbnail?path=${encodeURIComponent(item.path)}`
  return null
}

async function loadPaths() {
  try { boundPaths.value = await fetch(`/api/collections/${route.params.id}/paths`).then(r => r.json()) } catch {}
}
async function openPathManager() { await loadPaths(); showPathManager.value = true }
async function addPath() {
  if (!newPath.value.trim()) return
  try {
    await fetch(`/api/collections/${route.params.id}/paths`, {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ path: newPath.value.trim(), auto_scan: newPathAutoScan.value }),
    })
    newPath.value = ''
    await loadPaths()
    if (newPathAutoScan.value && boundPaths.value.length) await scanPath(boundPaths.value[boundPaths.value.length - 1].id)
  } catch (e: any) { alert('添加失败: ' + e.message) }
}
async function scanPath(pathId: string) {
  try {
    const res = await fetch(`/api/paths/${pathId}/scan`, { method: 'POST' })
    const data = await res.json()
    if (data.status === 'ok') await loadData()
  } catch (e: any) { alert('扫描失败: ' + e.message) }
}
async function removePath(pathId: string) {
  if (!confirm('确定移除此路径？')) return
  try { await fetch(`/api/paths/${pathId}`, { method: 'DELETE' }); await loadPaths(); await loadData() } catch {}
}
async function addVFile() {
  if (!newVFile.value.path.trim()) return
  try {
    await fetch(`/api/collections/${route.params.id}/vfiles`, {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(newVFile.value),
    })
    newVFile.value = { path: '', display_name: '', size: 0 }
    showVFileCreator.value = false
    await loadData()
  } catch (e: any) { alert('添加失败: ' + e.message) }
}
</script>

<template>
  <div class="browse-page">
    <div class="breadcrumb-bar">
      <span class="crumb" @click="goToRoot">全部文件夹</span>
      <template v-for="(crumb, idx) in breadcrumb" :key="crumb.id">
        <span class="crumb-sep">›</span>
        <span :class="['crumb', { active: idx === breadcrumb.length - 1 }]" @click="goToBreadcrumb(idx)">{{ crumb.name }}</span>
      </template>
    </div>

    <div class="toolbar">
      <input v-model="searchQuery" class="search-input" placeholder="搜索文件…" />
      <select v-model="sortBy" class="sort-select" @change="loadData">
        <option value="name">名称排序</option>
        <option value="size">大小排序</option>
        <option value="time">时间排序</option>
      </select>
      <button class="btn" @click="refresh">🔄 刷新</button>
      <button class="btn" @click="openPathManager">📂 路径</button>
      <button class="btn" @click="showVFileCreator = true">📎 虚拟文件</button>
      <div class="view-toggle">
        <button :class="{ on: viewMode === 'grid' }" @click="viewMode = 'grid'">▦</button>
        <button :class="{ on: viewMode === 'list' }" @click="viewMode = 'list'">☰</button>
      </div>
    </div>

    <div class="content-area">
      <div v-if="loading" class="status-text">加载中…</div>

      <div v-else-if="viewMode === 'grid'" class="grid-view">
        <div v-if="items.length === 0" class="status-text">此文件夹为空。点 📂 路径 添加并扫描。</div>
        <div v-else class="grid">
          <div v-for="item in items" :key="item.path"
            :class="['grid-item', { selected: selectedItem?.path === item.path }]"
            @click="selectItem(item)" @dblclick="openItem(item)">
            <div class="grid-icon">{{ getIcon(item) }}</div>
            <div class="grid-name">{{ item.name }}</div>
            <div class="grid-meta">
              <span class="source-tag" :class="'source-' + item.source">{{ item.source }}</span>
              <span class="file-size">{{ formatSize(item.size) }}</span>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="list-view">
        <div class="list-header">
          <span class="list-col-name">名称</span><span class="list-col-source">来源</span>
          <span class="list-col-size">大小</span><span class="list-col-time">修改时间</span>
        </div>
        <div v-for="item in items" :key="item.path"
          :class="['list-item', { selected: selectedItem?.path === item.path }]"
          @click="selectItem(item)" @dblclick="openItem(item)">
          <span class="list-col-name"><span class="list-icon">{{ getIcon(item) }}</span>{{ item.name }}</span>
          <span class="list-col-source"><span class="source-tag" :class="'source-' + item.source">{{ item.source }}</span></span>
          <span class="list-col-size">{{ formatSize(item.size) }}</span>
          <span class="list-col-time">{{ formatTime(item.mod_time) }}</span>
        </div>
      </div>

      <div v-if="selectedItem" class="preview-panel">
        <div class="preview-header">
          <span>{{ selectedItem.name }}</span>
          <button class="close-btn" @click="selectedItem = null">✕</button>
        </div>
        <div class="preview-body">
          <template v-if="previewSrc(selectedItem)">
            <img :src="previewSrc(selectedItem)!" :alt="selectedItem.name" />
          </template>
          <template v-else-if="selectedItem.mime_type?.startsWith('video/')">
            <video :src="`/api/preview/thumbnail?path=${encodeURIComponent(selectedItem.path)}`" controls style="max-width:100%;max-height:100%"></video>
          </template>
          <template v-else>
            <div class="preview-placeholder">
              <div style="font-size:48px">{{ getIcon(selectedItem) }}</div>
              <div style="margin-top:8px;color:var(--text-muted)">{{ selectedItem.name }}</div>
              <div style="margin-top:4px;font-size:12px;color:var(--text-dim)">{{ formatSize(selectedItem.size) }} · {{ selectedItem.mime_type || '未知类型' }}</div>
            </div>
          </template>
        </div>
      </div>
    </div>

    <div class="statusbar">共 {{ items.length }} 项 · {{ viewMode === 'grid' ? '网格' : '列表' }}视图</div>

    <Teleport to="body">
      <div v-if="showPathManager" class="modal-overlay" @click.self="showPathManager = false">
        <div class="modal wider">
          <div class="modal-header">📂 管理「{{ collection?.name }}」的绑定路径</div>
          <div class="modal-body">
            <div v-if="boundPaths.length" class="path-list">
              <div v-for="p in boundPaths" :key="p.id" class="path-item">
                <code class="path-text">{{ p.path }}</code>
                <span class="path-badge" :class="p.auto_scan ? 'badge-auto' : 'badge-manual'">{{ p.auto_scan ? '自动扫描' : '手动' }}</span>
                <div class="path-actions">
                  <button class="btn btn-sm" @click="scanPath(p.id)">🔍 扫描</button>
                  <button class="btn btn-sm btn-danger" @click="removePath(p.id)">🗑️</button>
                </div>
              </div>
            </div>
            <div v-else class="no-paths">暂无绑定路径</div>
            <div class="add-path-row">
              <input v-model="newPath" class="input flex-1" placeholder="路径，如 115:\\学习\\Java\\ 或 D:\\Videos\\" @keyup.enter="addPath" />
              <label class="checkbox-label"><input type="checkbox" v-model="newPathAutoScan" /> 自动扫描</label>
              <button class="btn btn-primary btn-sm" @click="addPath">添加</button>
            </div>
          </div>
          <div class="modal-footer"><button class="btn" @click="showPathManager = false">关闭</button></div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showVFileCreator" class="modal-overlay" @click.self="showVFileCreator = false">
        <div class="modal">
          <div class="modal-header">📎 添加虚拟文件</div>
          <div class="modal-body">
            <label>路径（含前缀）</label>
            <input v-model="newVFile.path" class="input" placeholder="如 tg:\\Java频道 或 115:\\教程\\视频.mp4" />
            <label>显示名称（可选）</label>
            <input v-model="newVFile.display_name" class="input" placeholder="留空则使用路径" />
            <label>文件大小（字节，可选）</label>
            <input v-model.number="newVFile.size" class="input" type="number" placeholder="0" />
          </div>
          <div class="modal-footer">
            <button class="btn" @click="showVFileCreator = false">取消</button>
            <button class="btn btn-primary" @click="addVFile">添加</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<style scoped>
.browse-page { display: flex; flex-direction: column; height: 100%; }
.breadcrumb-bar { padding: 8px 16px; font-size: 13px; background: var(--bg-secondary); border-bottom: 1px solid var(--border-light); display: flex; align-items: center; gap: 4px; flex-shrink: 0; }
.crumb { color: var(--accent); cursor: pointer; }
.crumb:hover { text-decoration: underline; }
.crumb.active { color: var(--text-primary); cursor: default; font-weight: 500; }
.crumb.active:hover { text-decoration: none; }
.crumb-sep { color: var(--text-dim); }
.toolbar { display: flex; align-items: center; gap: 10px; padding: 10px 16px; border-bottom: 1px solid var(--border-light); flex-shrink: 0; }
.search-input { background: var(--bg-input); border: 1px solid var(--border); border-radius: 6px; padding: 6px 12px; color: var(--text-primary); font-size: 13px; width: 200px; outline: none; }
.search-input:focus { border-color: var(--accent); }
.sort-select { background: var(--bg-input); border: 1px solid var(--border); border-radius: 6px; padding: 6px 10px; color: var(--text-secondary); font-size: 13px; outline: none; }
.btn { padding: 6px 12px; border-radius: 6px; border: 1px solid var(--border); background: var(--bg-card); color: var(--text-secondary); cursor: pointer; font-size: 13px; }
.btn:hover { background: var(--bg-card-hover); }
.view-toggle { display: flex; gap: 2px; margin-left: auto; }
.view-toggle button { padding: 6px 10px; border: 1px solid var(--border); background: var(--bg-card); color: var(--text-secondary); cursor: pointer; font-size: 13px; border-radius: 4px; }
.view-toggle button.on { background: var(--accent-bg); color: #eef; }
.content-area { flex: 1; display: flex; overflow: hidden; }
.grid-view, .list-view { flex: 1; overflow-y: auto; padding: 16px; }
.grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(120px, 1fr)); gap: 12px; }
.grid-item { background: var(--bg-card); border-radius: 8px; padding: 12px; text-align: center; cursor: pointer; border: 1px solid transparent; transition: all 0.15s; }
.grid-item:hover { border-color: var(--accent); background: var(--bg-card-hover); }
.grid-item.selected { border-color: var(--accent); background: var(--bg-card-hover); }
.grid-icon { font-size: 32px; margin-bottom: 6px; }
.grid-name { font-size: 12px; color: var(--text-secondary); word-break: break-all; line-height: 1.3; }
.grid-meta { display: flex; gap: 4px; justify-content: center; align-items: center; margin-top: 4px; }
.source-tag { font-size: 10px; padding: 1px 6px; border-radius: 3px; }
.source-tag { background: var(--tag-local-bg); color: var(--tag-local-text); }
.source-tag.source-tg, .source-tag.source-notion { background: var(--tag-web-bg); color: var(--tag-web-text); }
.file-size { font-size: 10px; color: var(--text-dim); }
.list-header { display: flex; padding: 8px 12px; font-size: 12px; color: var(--text-muted); border-bottom: 1px solid var(--border-light); position: sticky; top: 0; background: var(--bg-primary); z-index: 1; }
.list-item { display: flex; align-items: center; padding: 8px 12px; border-radius: 6px; cursor: pointer; font-size: 13px; }
.list-item:hover { background: var(--bg-card); }
.list-item.selected { background: var(--bg-card-hover); }
.list-col-name { flex: 2; display: flex; align-items: center; gap: 8px; color: var(--text-secondary); }
.list-col-source { flex: 1; }
.list-col-size { width: 80px; text-align: right; color: var(--text-muted); font-size: 12px; }
.list-col-time { width: 140px; text-align: right; color: var(--text-dim); font-size: 12px; }
.list-icon { font-size: 18px; }
.preview-panel { width: 45%; background: var(--preview-bg); border-left: 1px solid var(--border); display: flex; flex-direction: column; flex-shrink: 0; }
.preview-header { padding: 8px 14px; font-size: 13px; border-bottom: 1px solid var(--border-light); display: flex; align-items: center; gap: 8px; }
.close-btn { margin-left: auto; background: none; border: none; color: var(--text-muted); cursor: pointer; font-size: 14px; }
.preview-body { flex: 1; display: flex; align-items: center; justify-content: center; overflow: auto; }
.preview-body img { max-width: 100%; max-height: 100%; object-fit: contain; }
.preview-placeholder { text-align: center; color: var(--text-dim); }
.statusbar { padding: 6px 16px; font-size: 11px; color: var(--text-dim); border-top: 1px solid var(--border-light); flex-shrink: 0; }
.status-text { color: var(--text-muted); text-align: center; padding: 60px; }
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.6); display: flex; align-items: center; justify-content: center; z-index: 100; }
.modal { background: var(--bg-card); border: 1px solid var(--border); border-radius: 10px; width: 440px; max-height: 80vh; overflow-y: auto; }
.modal.wider { width: 580px; }
.modal-header { padding: 14px 18px; border-bottom: 1px solid var(--border-light); font-size: 15px; }
.modal-body { padding: 18px; }
.modal-body label { display: block; font-size: 12px; color: var(--text-muted); margin-bottom: 4px; margin-top: 12px; }
.modal-body label:first-of-type { margin-top: 0; }
.modal-footer { padding: 12px 18px; border-top: 1px solid var(--border-light); display: flex; justify-content: flex-end; gap: 8px; }
.input { width: 100%; background: var(--bg-input); border: 1px solid var(--border); border-radius: 6px; padding: 8px 12px; color: var(--text-primary); font-size: 13px; outline: none; }
.input:focus { border-color: var(--accent); }
.input.flex-1 { flex: 1; }
.btn-primary { background: var(--accent-bg); color: #eef; border-color: var(--accent-bg); }
.btn-primary:hover { background: var(--accent-hover); }
.btn-sm { padding: 3px 10px; font-size: 12px; }
.btn-danger { color: var(--danger); border-color: var(--danger); }
.btn-danger:hover { background: var(--danger); color: #fff; }
.path-list { display: flex; flex-direction: column; gap: 8px; margin-bottom: 16px; }
.path-item { display: flex; align-items: center; gap: 10px; padding: 10px; background: var(--bg-secondary); border-radius: 6px; }
.path-text { flex: 1; font-size: 13px; color: var(--accent); font-family: 'Consolas', monospace; word-break: break-all; background: none; }
.path-badge { font-size: 11px; padding: 2px 8px; border-radius: 10px; white-space: nowrap; }
.badge-auto { background: #1a3a1a; color: #6c6; }
.badge-manual { background: #3a3a1a; color: #cc6; }
.path-actions { display: flex; gap: 4px; flex-shrink: 0; }
.no-paths { text-align: center; color: var(--text-dim); padding: 20px; font-size: 13px; margin-bottom: 12px; }
.add-path-row { display: flex; align-items: center; gap: 8px; }
.checkbox-label { display: flex; align-items: center; gap: 4px; font-size: 12px; color: var(--text-secondary); white-space: nowrap; cursor: pointer; }
</style>