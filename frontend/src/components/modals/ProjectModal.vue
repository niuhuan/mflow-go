<script setup>
import BaseModal from './BaseModal.vue';
import { backend } from '../../api/backend';
import { dialogs } from '../../composables/useDialogs';
import { templates } from '../../blockly/templates';

const emit = defineEmits(['close', 'apply', 'save-as']);

async function openProject() {
  try {
    const path = await backend.openProjectDialog();
    if (!path) return;
    const text = await backend.readTextFile(path);
    emit('apply', { path, text });
    emit('close');
  } catch (e) {
    await dialogs.alert('打开失败: ' + e);
  }
}

async function newFromTemplate(tpl) {
  const ok = await dialogs.confirm(`使用「${tpl.name}」新建工程？未保存的更改将丢失。`);
  if (!ok) return;
  emit('apply', { path: '', text: tpl.xml });
  emit('close');
}

function saveAs() {
  emit('save-as');
  emit('close');
}
</script>

<template>
  <BaseModal title="工程" width="440px" @close="emit('close')">
    <h4 class="section">打开 / 保存</h4>
    <div class="row">
      <button class="tool-btn" @click="openProject">打开工程...</button>
      <button class="tool-btn" @click="saveAs">另存为...</button>
    </div>

    <h4 class="section">新建（模板）</h4>
    <div class="row col">
      <button v-for="tpl in templates" :key="tpl.key" class="tool-btn" @click="newFromTemplate(tpl)">
        {{ tpl.name }}
      </button>
    </div>
  </BaseModal>
</template>

<style scoped>
.section {
  margin: 4px 0 10px;
  font-size: 13px;
  color: #4a5568;
  border-left: 3px solid #3182ce;
  padding-left: 8px;
}
.row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 14px;
}
.row.col {
  flex-direction: column;
}
.tool-btn {
  flex: 1 1 auto;
  min-width: 130px;
  padding: 9px 12px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  background: #fff;
  cursor: pointer;
  font-size: 13px;
}
.tool-btn:hover {
  background: #ebf8ff;
  border-color: #3182ce;
}
</style>
