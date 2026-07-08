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
import { loadLastFile, saveLastFile } from './composables/useConfig';

const workspaceRef = ref(null);
const messages = ref([]);
const version = ref('');
const newVersion = ref('');
const filePath = ref('');
const dirty = ref(false);

// 最近一次已保存/已加载的工程内容快照，用于判断是否"脏"。
let lastSavedText = '';

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

// 记录当前内容为"已保存"基线。
function markSaved() {
  try {
    lastSavedText = workspaceRef.value.toJsonText();
  } catch (_) {
    lastSavedText = '';
  }
  dirty.value = false;
}

// 内容变化时通过与基线比较来判断是否真的被编辑过。
function onChange() {
  try {
    dirty.value = workspaceRef.value.toJsonText() !== lastSavedText;
  } catch (_) {
    /* ignore */
  }
}

// 按扩展名选择序列化格式：.m7f 写 XML，其余（.m7p）写 JSON。
function serializeFor(path) {
  if (/\.m7f$/i.test(path || '')) {
    return workspaceRef.value.toXmlText();
  }
  return workspaceRef.value.toJsonText();
}

async function save() {
  try {
    if (!filePath.value) {
      // 首次保存（新建/模板工程）走"另存为"，默认 .m7p。
      const picked = await backend.saveProjectDialog();
      if (!picked) return false;
      filePath.value = picked;
    }
    await backend.writeTextFile(filePath.value, serializeFor(filePath.value));
    // 脏检测基线始终使用 JSON，与 onChange 保持一致。
    lastSavedText = workspaceRef.value.toJsonText();
    dirty.value = false;
    await saveLastFile(filePath.value);
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
    await backend.writeTextFile(picked, serializeFor(picked));
    lastSavedText = workspaceRef.value.toJsonText();
    dirty.value = false;
    await saveLastFile(picked);
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
  if (path) {
    // 打开已存在工程（.m7p 或 .m7f）：关联路径，保存时按扩展名写回对应格式
    filePath.value = path;
    markSaved();
    await saveLastFile(path);
    log('已打开工程: ' + path);
  } else {
    // 新建/模板工程：不关联任何路径，首次保存走另存为（默认 .m7p）
    filePath.value = '';
    markSaved();
    await saveLastFile('');
    log('已新建工程（首次保存将提示选择保存位置）');
  }
}

onMounted(async () => {
  try {
    version.value = await backend.getVersion();
  } catch (_) { /* ignore */ }
  backend.getNewVersion().then((v) => { newVersion.value = v || ''; }).catch(() => {});

  try {
    const autoRun = await backend.getAutoRunFile();
    if (autoRun && (await backend.exists(autoRun))) {
      const text = await backend.readTextFile(autoRun);
      workspaceRef.value.loadFromText(text);
      filePath.value = autoRun;
      markSaved();
      await saveLastFile(autoRun);
      log('已打开工程: ' + autoRun);
      setTimeout(() => runAndSave(), 300);
      return;
    }

    // 非首次：恢复上次打开的工程（若仍存在）。首次打开则保持空白工作区。
    const last = await loadLastFile();
    if (last && (await backend.exists(last))) {
      const text = await backend.readTextFile(last);
      workspaceRef.value.loadFromText(text);
      filePath.value = last;
      markSaved();
      log('已恢复上次工程: ' + last);
    } else {
      // 空白工作区，未关联文件，且不标记为已编辑。
      markSaved();
    }
  } catch (e) {
    log('初始化失败: ' + e);
    markSaved();
  }
});

function onReady() {
  // 工作区就绪时以当前（空白）内容作为基线，避免误判为已编辑。
  markSaved();
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
