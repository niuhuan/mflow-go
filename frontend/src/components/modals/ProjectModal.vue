<script setup>
import PopupMenu from '../PopupMenu.vue';
import { backend } from '../../api/backend';
import { dialogs } from '../../composables/useDialogs';
import { templates } from '../../blockly/templates';

defineProps({
  show: { type: Boolean, default: false },
  anchor: { type: Object, default: null },
});

const emit = defineEmits(['close', 'apply', 'save-as']);

async function openProject() {
  emit('close');
  try {
    const path = await backend.openProjectDialog();
    if (!path) return;
    const text = await backend.readTextFile(path);
    emit('apply', { path, text });
  } catch (e) {
    await dialogs.alert('打开失败: ' + e);
  }
}

async function newFromTemplate(tpl) {
  emit('close');
  const ok = await dialogs.confirm(`使用「${tpl.name}」新建工程？未保存的更改将丢失。`);
  if (!ok) return;
  emit('apply', { path: '', text: tpl.xml });
}

function saveAs() {
  emit('close');
  emit('save-as');
}
</script>

<template>
  <PopupMenu :show="show" :anchor="anchor" :min-width="220" @close="emit('close')">
    <button class="menu-item" @click="openProject">打开工程...</button>
    <button class="menu-item" @click="saveAs">另存为...</button>

    <div class="menu-divider" />
    <div class="menu-label">新建（模板）</div>
    <button v-for="tpl in templates" :key="tpl.key" class="menu-item" @click="newFromTemplate(tpl)">
      {{ tpl.name }}
    </button>
  </PopupMenu>
</template>

<style scoped>
.menu-label {
  padding: 6px 14px 4px;
  font-size: 11px;
  color: #a0aec0;
  font-weight: 600;
  letter-spacing: 0.03em;
}
.menu-item {
  display: block;
  width: 100%;
  text-align: left;
  padding: 8px 14px;
  border: none;
  background: transparent;
  cursor: pointer;
  font-size: 13px;
  color: #2d3748;
  white-space: nowrap;
}
.menu-item:hover {
  background: #ebf8ff;
  color: #2b6cb0;
}
.menu-divider {
  height: 1px;
  background: #edf2f7;
  margin: 4px 0;
}
</style>
