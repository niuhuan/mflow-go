<script setup>
import BaseModal from './BaseModal.vue';
import { backend } from '../../api/backend';
import { dialogs } from '../../composables/useDialogs';

const emit = defineEmits(['close']);

const launchers = [
  { label: '启动三月七小助手 (GUI)', fn: () => backend.runM7Launcher() },
  { label: '启动更好的原神 (GUI)', fn: () => backend.runBetterGiGui() },
  { label: '启动绝区零一条龙 (GUI)', fn: () => backend.runZzzodGui() },
];

async function runLauncher(item) {
  try {
    await item.fn();
  } catch (e) {
    await dialogs.alert(String(e));
  }
}

async function exportAccount(kind) {
  const label = { hsr: '星铁', gi: '原神', zzz: '绝区零' }[kind];
  const name = await dialogs.prompt(`请输入${label}账号名称`);
  if (!name || !name.trim()) return;
  try {
    if (kind === 'hsr') await backend.exportAccount(name.trim(), '', '');
    else if (kind === 'gi') await backend.exportGiAccount(name.trim());
    else await backend.exportZzzAccount(name.trim());
    await dialogs.alert(`导出${label}账号成功`);
  } catch (e) {
    await dialogs.alert(String(e));
  }
}

async function importAccount(kind) {
  const label = { hsr: '星铁', gi: '原神', zzz: '绝区零' }[kind];
  try {
    let list = [];
    if (kind === 'hsr') list = await backend.listAccounts();
    else if (kind === 'gi') list = await backend.listGiAccounts();
    else list = await backend.listZzzAccounts();
    list = list || [];
    if (list.length === 0) {
      await dialogs.alert(`没有${label}账号`);
      return;
    }
    const name = await dialogs.select(list, `选择要导入的${label}账号`);
    if (!name) return;
    if (kind === 'hsr') await backend.loadAccount(name);
    else if (kind === 'gi') await backend.importGiAccount(name);
    else await backend.importZzzAccount(name);
    await dialogs.alert(`导入${label}账号成功`);
  } catch (e) {
    await dialogs.alert(String(e));
  }
}

const clearRegItems = [
  { label: '清除星铁注册表', name: '星铁', fn: () => backend.clearHsrReg() },
  { label: '清除原神注册表', name: '原神', fn: () => backend.clearGiReg() },
  { label: '清除绝区零注册表', name: '绝区零', fn: () => backend.clearZzzReg() },
  { label: '清除鸣潮注册表', name: '鸣潮', fn: () => backend.clearWwReg() },
];

async function clearReg(item) {
  const ok = await dialogs.confirm(`确定清除${item.name}游戏注册表？此操作不可撤销。`);
  if (!ok) return;
  try {
    await item.fn();
    await dialogs.alert(`清除${item.name}注册表成功`);
  } catch (e) {
    await dialogs.alert(String(e));
  }
}
</script>

<template>
  <BaseModal title="工具" width="480px" @close="emit('close')">
    <h4 class="section">启动程序</h4>
    <div class="row">
      <button v-for="l in launchers" :key="l.label" class="tool-btn" @click="runLauncher(l)">
        {{ l.label }}
      </button>
    </div>

    <h4 class="section">账号导出</h4>
    <div class="row">
      <button class="tool-btn" @click="exportAccount('hsr')">导出星铁账号</button>
      <button class="tool-btn" @click="exportAccount('gi')">导出原神账号</button>
      <button class="tool-btn" @click="exportAccount('zzz')">导出绝区零账号</button>
    </div>

    <h4 class="section">账号导入</h4>
    <div class="row">
      <button class="tool-btn" @click="importAccount('hsr')">导入星铁账号</button>
      <button class="tool-btn" @click="importAccount('gi')">导入原神账号</button>
      <button class="tool-btn" @click="importAccount('zzz')">导入绝区零账号</button>
    </div>

    <h4 class="section">清除注册表</h4>
    <div class="row">
      <button v-for="item in clearRegItems" :key="item.label" class="tool-btn danger" @click="clearReg(item)">
        {{ item.label }}
      </button>
    </div>
  </BaseModal>
</template>

<style scoped>
.section {
  margin: 4px 0 10px;
  font-size: 13px;
  color: #4a5568;
  border-left: 3px solid #dd6b20;
  padding-left: 8px;
}
.row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 14px;
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
  background: #fffaf0;
  border-color: #dd6b20;
}
.tool-btn.danger:hover {
  background: #fff5f5;
  border-color: #c53030;
  color: #c53030;
}
</style>
