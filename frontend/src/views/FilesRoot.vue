<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { collections, type Collection } from '../api'

const router = useRouter()
const items = ref<Collection[]>([])
const loading = ref(true)
const viewMode = ref<'grid' | 'list'>('grid')
const showCreate = ref(false)
const editItem = ref<Collection | null>(null)
const searchQuery = ref('')
const tagCache = ref<Record<string, any[]>>({})

onMounted(() => loadCollections())

async function loadCollections() {
  loading.value = true
  try {
    items.value = await collections.listRoot()
    for (const c of items.value) {
      try { const t = await fetch(`/api/collections/${c.id}/tags`).then(r => r.json()); tagCache.value[c.id] = t } catch {}
    }
  } catch (e) { console.error(e) }
  finally { loading.value = false }
}

const form = ref({ name: '', icon: '📁' })

function openCreate() { form.value = { name: '', icon: '📁' }; editItem.value = null; showCreate.value = true }
function openEdit(item: Collection) { form.value = { name: item.name, icon: item.icon }; editItem.value = item; showCreate.value = true }

async function save() {
  if (!form.value.name.trim()) return
  try {
    if (editItem.value) await collections.update(editItem.value.id, form.value)
    else await collections.create(form.value)
    showCreate.value = false; await loadCollections()
  } catch (e) { console.error(e) }
}

async function remove(item: Collection) {
  if (!confirm(`确定删除「${item.name}」？此操作不可撤销。`)) return
  try { await collections.delete(item.id); await loadCollections() } catch (e) { console.error(e) }
}

function enter(id: string) { router.push(`/browse/${id}`) }

const filteredItems = computed(() => {
  if (!searchQuery.value) return items.value
  return items.value.filter(c => c.name.toLowerCase().includes(searchQuery.value.toLowerCase()))
})
</script>

<template>
  <div class="files-page">
    <div class="toolbar">
      <div class="breadcrumb">📁 全部文件夹</div>
      <div class="toolbar-actions">
        <input v-model="searchQuery" class="search-input" placeholder="搜索文件夹…" />
        <button class="btn btn-primary" @click="openCreate">+ 新建文件夹</button>
        <div class="view-toggle">
          <button :class="{ on: viewMode === 'grid' }" @click="viewMode = 'grid'">▦</button>
          <button :class="{ on: viewMode === 'list' }" @click="viewMode = 'list'">☰</button>
        </div>
      </div>
    </div>

    <div v-if="viewMode === 'grid'" class="grid-view">
      <div v-if="loading" class="status-text">加载中…</div>
      <div v-else-if="filteredItems.length === 0" class="status-text">
        {{ searchQuery ? '无匹配文件夹' : '还没有文件夹' }}
      </div>
      <div v-else class="grid">
        <div v-for="item in filteredItems" :key="item.id" class="grid-item" @dblclick="enter(item.id)" @contextmenu.prevent>
          <div class="grid-icon">{{ item.icon }}</div>
          <div class="grid-name">{{ item.name }}</div>
          <div class="grid-tags" v-if="tagCache[item.id]?.length">
            <span v-for="t in tagCache[item.id]" :key="t.id" class="tag-pill" :style="{ background: t.color + '33', color: t.color }">{{ t.name }}</span>
          </div>
          <div class="grid-actions">
            <button class="action-btn" @click.stop="openEdit(item)">✏️</button>
            <button class="action-btn" @click.stop="remove(item)">🗑️</button>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="list-view">
      <div class="list-header">
        <span class="list-col-name">名称</span><span class="list-col-tags">标签</span><span class="list-col-actions">操作</span>
      </div>
      <div v-for="item in filteredItems" :key="item.id" class="list-item" @dblclick="enter(item.id)">
        <span class="list-col-name"><span class="list-icon">{{ item.icon }}</span>{{ item.name }}</span>
        <span class="list-col-tags">
          <span v-for="t in (tagCache[item.id] || [])" :key="t.id" class="tag-pill" :style="{ background: t.color + '33', color: t.color }">{{ t.name }}</span>
        </span>
        <span class="list-col-actions">
          <button class="action-btn" @click.stop="openEdit(item)">✏️</button>
          <button class="action-btn" @click.stop="remove(item)">🗑️</button>
        </span>
      </div>
    </div>

    <div class="statusbar">共 {{ items.length }} 个文件夹</div>

    <Teleport to="body">
      <div v-if="showCreate" class="modal-overlay" @click.self="showCreate = false">
        <div class="modal">
          <div class="modal-header">{{ editItem ? '编辑文件夹' : '新建文件夹' }}</div>
          <div class="modal-body">
            <label>图标</label><input v-model="form.icon" class="input" placeholder="📁" maxlength="4" />
            <label>名称</label><input v-model="form.name" class="input" placeholder="文件夹名称" @keyup.enter="save" />
          </div>
          <div class="modal-footer">
            <button class="btn" @click="showCreate = false">取消</button>
            <button class="btn btn-primary" @click="save">确定</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<style scoped>
.files-page { display: flex; flex-direction: column; height: 100%; }
.toolbar { display: flex; align-items: center; gap: 12px; padding: 10px 16px; border-bottom: 1px solid var(--border-light); flex-shrink: 0; }
.breadcrumb { font-size: 14px; color: var(--text-secondary); font-weight: 500; }
.toolbar-actions { display: flex; align-items: center; gap: 8px; margin-left: auto; }
.search-input { background: var(--bg-input); border: 1px solid var(--border); border-radius: 6px; padding: 6px 12px; color: var(--text-primary); font-size: 13px; width: 180px; outline: none; }
.search-input:focus { border-color: var(--accent); }
.btn { padding: 6px 14px; border-radius: 6px; border: 1px solid var(--border); background: var(--bg-card); color: var(--text-secondary); cursor: pointer; font-size: 13px; transition: background 0.15s; }
.btn:hover { background: var(--bg-card-hover); }
.btn-primary { background: var(--accent-bg); color: #eef; border-color: var(--accent-bg); }
.btn-primary:hover { background: var(--accent-hover); }
.view-toggle { display: flex; gap: 2px; }
.view-toggle button { padding: 6px 10px; border: 1px solid var(--border); background: var(--bg-card); color: var(--text-secondary); cursor: pointer; font-size: 13px; border-radius: 4px; }
.view-toggle button.on { background: var(--accent-bg); color: #eef; }
.grid-view { flex: 1; overflow-y: auto; padding: 16px; }
.grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(140px, 1fr)); gap: 12px; }
.grid-item { background: var(--bg-card); border-radius: 8px; padding: 16px; text-align: center; cursor: pointer; border: 1px solid transparent; position: relative; transition: all 0.15s; }
.grid-item:hover { border-color: var(--accent); background: var(--bg-card-hover); }
.grid-icon { font-size: 36px; margin-bottom: 6px; }
.grid-name { font-size: 13px; color: var(--text-secondary); word-break: break-all; }
.grid-tags { margin-top: 6px; display: flex; gap: 4px; flex-wrap: wrap; justify-content: center; }
.tag-pill { font-size: 10px; padding: 1px 8px; border-radius: 10px; white-space: nowrap; }
.grid-actions { display: none; position: absolute; top: 6px; right: 6px; gap: 2px; }
.grid-item:hover .grid-actions { display: flex; }
.action-btn { width: 22px; height: 22px; border-radius: 4px; border: none; background: var(--bg-secondary); cursor: pointer; font-size: 11px; display: flex; align-items: center; justify-content: center; }
.action-btn:hover { background: var(--accent-bg); }
.list-view { flex: 1; overflow-y: auto; padding: 0 16px; }
.list-header { display: flex; padding: 8px 12px; font-size: 12px; color: var(--text-muted); border-bottom: 1px solid var(--border-light); position: sticky; top: 0; background: var(--bg-primary); }
.list-item { display: flex; align-items: center; padding: 8px 12px; border-radius: 6px; cursor: pointer; font-size: 13px; gap: 10px; }
.list-item:hover { background: var(--bg-card); }
.list-col-name { flex: 1; display: flex; align-items: center; gap: 8px; color: var(--text-secondary); }
.list-col-tags { flex: 1; display: flex; gap: 4px; flex-wrap: wrap; }
.list-col-actions { width: 60px; display: flex; gap: 4px; }
.list-icon { font-size: 18px; }
.statusbar { padding: 6px 16px; font-size: 11px; color: var(--text-dim); border-top: 1px solid var(--border-light); flex-shrink: 0; }
.status-text { color: var(--text-muted); text-align: center; padding: 60px; }
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.6); display: flex; align-items: center; justify-content: center; z-index: 100; }
.modal { background: var(--bg-card); border: 1px solid var(--border); border-radius: 10px; width: 400px; max-height: 80vh; overflow-y: auto; }
.modal-header { padding: 14px 18px; border-bottom: 1px solid var(--border-light); font-size: 15px; }
.modal-body { padding: 18px; }
.modal-body label { display: block; font-size: 12px; color: var(--text-muted); margin-bottom: 4px; margin-top: 12px; }
.modal-body label:first-of-type { margin-top: 0; }
.input { width: 100%; background: var(--bg-input); border: 1px solid var(--border); border-radius: 6px; padding: 8px 12px; color: var(--text-primary); font-size: 13px; outline: none; }
.input:focus { border-color: var(--accent); }
.modal-footer { padding: 12px 18px; border-top: 1px solid var(--border-light); display: flex; justify-content: flex-end; gap: 8px; }
</style>
