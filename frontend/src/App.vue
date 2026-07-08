<script setup>
import { ref, onMounted } from 'vue';
import Toolbar from './components/Toolbar.vue';
import BlocklyWorkspace from './components/BlocklyWorkspace.vue';
import ConsolePanel from './components/ConsolePanel.vue';
import DialogHost from './components/DialogHost.vue';
import ToolsModal from './components/modals/ToolsModal.vue';
import SettingsModal from './components/modals/SettingsModal.vue';
import ProjectModal from './components/modals/ProjectModal.vue';
import { useRunner } from './composables/useRunner';
import { backend } from './api/backend';
import { dialogs } from './composables/useDialogs';

const workspaceRef = ref(null);
const messages = ref([]);
const version = ref('');
const newVersion = ref('');
const filePath = ref('');
const dirty = ref(false);

const showTools = ref(false);
const showSettings = ref(false);
const showProject = ref(false);

function log(msg) {
  const time = new Date().toLocaleTimeString();
  messages.value.push(`${time} ${msg}`);
  if (messages.value.length > 200) messages.value.shift();
}

const getWorkspace = () => workspaceRef.value?.getWorkspace();
const { running, run, interrupt } = useRunner(getWorkspace, log);

function onChange() {
  dirty.value = true;
}

async function defaultProjectPath() {
  const dir = await backend.appDataPath();
  return dir + '/default.m7p';
}

async function save() {
  try {
    if (!filePath.value) {
      const picked = await backend.saveProjectDialog();
      if (!picked) return false;
      filePath.value = picked;
    }
    const text = workspaceRef.value.toJsonText();
    await backend.writeTextFile(filePath.value, text);
    dirty.value = false;
    log('已保存到: ' + filePath.value);
    return true;
  } catch (e) {
    log('保存失败: ' + e);
    return false;
  }
}

async function saveAs() {
  try {
    const picked = await backend.saveProjectDialog();
    if (!picked) return;
    filePath.value = picked;
    const text = workspaceRef.value.toJsonText();
    await backend.writeTextFile(picked, text);
    dirty.value = false;
    log('已另存为: ' + picked);
  } catch (e) {
    log('另存失败: ' + e);
  }
}

async function runAndSave() {
  await save();
  await run();
}

async function applyProject({ path, text }) {
  if (dirty.value) {
    const ok = await dialogs.confirm('当前工程未保存，确定放弃更改？');
    if (!ok) return;
  }
  workspaceRef.value.loadFromText(text);
  filePath.value = path;
  dirty.value = false;
  log(path ? '已打开工程: ' + path : '已新建工程（请另存为以保存）');
}

onMounted(async () => {
  try {
    version.value = await backend.getVersion();
  } catch (_) { /* ignore */ }
  backend.getNewVersion().then((v) => { newVersion.value = v || ''; }).catch(() => {});

  try {
    const autoRun = await backend.getAutoRunFile();
    const target = autoRun || (await defaultProjectPath());
    if (await backend.exists(target)) {
      const text = await backend.readTextFile(target);
      workspaceRef.value.loadFromText(text);
      filePath.value = target;
      dirty.value = false;
      log('已打开工程: ' + target);
      if (autoRun) {
        setTimeout(() => runAndSave(), 300);
      }
    }
  } catch (e) {
    log('初始化失败: ' + e);
  }
});

function onReady() {
  log('工作空间就绪。');
}
</script>

<template>
  <div class="app">
    <Toolbar
      :file-path="filePath"
      :dirty="dirty"
      :running="running"
      :version="version"
      :new-version="newVersion"
      @run="runAndSave"
      @interrupt="interrupt"
      @save="save"
      @open-tools="showTools = true"
      @open-settings="showSettings = true"
      @open-project="showProject = true"
      @open-release="backend.openReleasePage()"
    />
    <div class="body">
      <div class="editor">
        <BlocklyWorkspace ref="workspaceRef" @change="onChange" @ready="onReady" />
      </div>
      <div class="console-wrap">
        <ConsolePanel :messages="messages" />
      </div>
    </div>

    <ToolsModal v-if="showTools" @close="showTools = false" />
    <SettingsModal v-if="showSettings" @close="showSettings = false" />
    <ProjectModal
      v-if="showProject"
      @close="showProject = false"
      @apply="applyProject"
      @save-as="saveAs"
    />
    <DialogHost />
  </div>
</template>

<style scoped>
.app {
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 100vw;
}
.body {
  flex: 1 1 auto;
  display: flex;
  flex-direction: column;
  min-height: 0;
}
.editor {
  flex: 1 1 auto;
  min-height: 0;
  position: relative;
}
.console-wrap {
  flex: 0 0 160px;
  border-top: 1px solid #333;
}
</style>
