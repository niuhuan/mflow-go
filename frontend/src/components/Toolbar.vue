<script setup>
import { ref } from 'vue';
import ProjectModal from './modals/ProjectModal.vue';
import ToolsModal from './modals/ToolsModal.vue';

defineProps({
  filePath: { type: String, default: '' },
  dirty: { type: Boolean, default: false },
  running: { type: Boolean, default: false },
  version: { type: String, default: '' },
  newVersion: { type: String, default: '' },
});

const emit = defineEmits([
  'run',
  'interrupt',
  'save',
  'open-settings',
  'open-release',
  'apply-project',
  'save-as',
]);

const projectBtn = ref(null);
const toolsBtn = ref(null);
const showProject = ref(false);
const showTools = ref(false);

function toggleProject() {
  showTools.value = false;
  showProject.value = !showProject.value;
}

function toggleTools() {
  showProject.value = false;
  showTools.value = !showTools.value;
}
</script>

<template>
  <div class="toolbar">
    <div class="left">
      <button class="btn primary" :disabled="running" @click="emit('run')">
        {{ running ? '运行中...' : '保存并运行' }}
      </button>
      <button class="btn danger" :disabled="!running" @click="emit('interrupt')">中断</button>
      <button class="btn" @click="emit('save')">保存（重新运行生效）</button>
      <span v-if="dirty" class="dirty">*</span>
    </div>

    <div class="center">
      <span class="file-path" :title="filePath">{{ filePath || '未打开工程' }}</span>
    </div>

    <div class="right">
      <button
        ref="projectBtn"
        class="btn"
        :class="{ active: showProject }"
        @click="toggleProject"
      >
        工程
      </button>
      <button ref="toolsBtn" class="btn" :class="{ active: showTools }" @click="toggleTools">
        工具
      </button>
      <button class="btn" @click="emit('open-settings')">设置</button>
      <span class="version" @click="emit('open-release')">
        <span v-if="newVersion" class="update-badge">有新版本</span>
        v{{ version }}
      </span>
    </div>

    <ProjectModal
      :show="showProject"
      :anchor="projectBtn"
      @close="showProject = false"
      @apply="emit('apply-project', $event)"
      @save-as="emit('save-as')"
    />
    <ToolsModal :show="showTools" :anchor="toolsBtn" @close="showTools = false" />
  </div>
</template>

<style scoped>
.toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 10px;
  background: #2b3648;
  color: #fff;
  flex: 0 0 auto;
}
.left,
.right {
  display: flex;
  align-items: center;
  gap: 6px;
}
.center {
  flex: 1 1 auto;
  overflow: hidden;
}
.file-path {
  display: block;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: #cdd6e4;
  font-size: 13px;
}
.btn {
  padding: 5px 12px;
  border: none;
  border-radius: 4px;
  background: #3f4c63;
  color: #fff;
  cursor: pointer;
  font-size: 13px;
}
.btn:hover:not(:disabled) {
  background: #4c5c78;
}
.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
.btn.active {
  background: #4c5c78;
  box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.25);
}
.btn.primary {
  background: #2f855a;
}
.btn.primary:hover:not(:disabled) {
  background: #38a169;
}
.btn.danger {
  background: #9b2c2c;
}
.btn.danger:hover:not(:disabled) {
  background: #c53030;
}
.dirty {
  color: #f6ad55;
  font-weight: bold;
}
.version {
  cursor: pointer;
  font-size: 12px;
  color: #cdd6e4;
  display: flex;
  align-items: center;
  gap: 4px;
}
.update-badge {
  background: #dd6b20;
  color: #fff;
  border-radius: 3px;
  padding: 1px 5px;
  font-size: 11px;
}
</style>
