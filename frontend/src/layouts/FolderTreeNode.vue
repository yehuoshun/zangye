<template>
  <div class="tree-node" :style="{ paddingLeft: depth * 16 + 'px' }">
    <div
      class="tree-node-label"
      :class="{ selected: node.id === selectedId }"
      @click="$emit('select', node.id)"
    >
      <span class="tree-toggle" @click.stop="toggle">
        {{ expanded ? '▼' : '▶' }}
      </span>
      <span class="tree-icon">📂</span>
      <span class="tree-name">{{ node.name }}</span>
      <span class="tree-count" v-if="node.file_count">{{ node.file_count }}</span>
    </div>
    <div v-if="expanded && node.children && node.children.length > 0">
      <FolderTreeNode
        v-for="child in node.children"
        :key="child.id"
        :node="child"
        :depth="depth + 1"
        :selected-id="selectedId"
        @select="$emit('select', $event)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { FolderItem } from '@/features/folders/types'

defineProps<{
  node: FolderItem
  depth: number
  selectedId: string | null
}>()

defineEmits<{
  select: [id: string]
}>()

const expanded = ref(true)

function toggle() {
  expanded.value = !expanded.value
}
</script>

<style scoped>
.tree-node-label {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  border-radius: var(--border-radius-sm);
  cursor: pointer;
  font-size: 13px;
  color: var(--text-secondary);
  transition: all var(--transition-fast);
}

.tree-node-label:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.tree-node-label.selected {
  background: var(--bg-active);
  color: var(--accent-primary);
}

.tree-toggle {
  font-size: 8px;
  width: 12px;
  text-align: center;
  flex-shrink: 0;
}

.tree-icon {
  font-size: 14px;
}

.tree-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tree-count {
  font-size: 11px;
  color: var(--text-muted);
  background: var(--bg-tertiary);
  padding: 0 6px;
  border-radius: 8px;
}
</style>
