<script setup>
import PopupMenu from '../PopupMenu.vue';
import { backend } from '../../api/backend';
import { dialogs } from '../../composables/useDialogs';

defineProps({
  show: { type: Boolean, default: false },
  anchor: { type: Object, default: null },
});

const emit = defineEmits(['close']);

const launchers = [
  { label: '启动三月七小助手 (GUI)', fn: () => backend.runM7Launcher() },
  { label: '启动更好的原神 (GUI)', fn: () => backend.runBetterGiGui() },
  { label: '启动绝区零一条龙 (GUI)', fn: () => backend.runZzzodGui() },
];

async function runLauncher(item) {
  emit('close');
  try {
    await item.fn();
  } catch (e) {
    await dialogs.alert(String(e));
  }
}

async function exportAccount(kind) {
  emit('close');
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
  emit('close');
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
  emit('close');
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
  <PopupMenu :show="show" :anchor="anchor" :min-width="240" @close="emit('close')">
    <div class="menu-label">启动程序</div>
    <button v-for="l in launchers" :key="l.label" class="menu-item" @click="runLauncher(l)">
      {{ l.label }}
    </button>

    <div class="menu-divider" />
    <div class="menu-label">账号导出</div>
    <button class="menu-item" @click="exportAccount('hsr')">导出星铁账号</button>
    <button class="menu-item" @click="exportAccount('gi')">导出原神账号</button>
    <button class="menu-item" @click="exportAccount('zzz')">导出绝区零账号</button>

    <div class="menu-divider" />
    <div class="menu-label">账号导入</div>
    <button class="menu-item" @click="importAccount('hsr')">导入星铁账号</button>
    <button class="menu-item" @click="importAccount('gi')">导入原神账号</button>
    <button class="menu-item" @click="importAccount('zzz')">导入绝区零账号</button>

    <div class="menu-divider" />
    <div class="menu-label">清除注册表</div>
    <button
      v-for="item in clearRegItems"
      :key="item.label"
      class="menu-item danger"
      @click="clearReg(item)"
    >
      {{ item.label }}
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
  background: #fffaf0;
  color: #c05621;
}
.menu-item.danger:hover {
  background: #fff5f5;
  color: #c53030;
}
.menu-divider {
  height: 1px;
  background: #edf2f7;
  margin: 4px 0;
}
</style>
