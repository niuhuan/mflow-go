<script setup>
import { onMounted, onBeforeUnmount, ref, markRaw } from 'vue';
import * as Blockly from 'blockly/core';
import * as zhHans from 'blockly/msg/zh-hans';
import 'blockly/blocks';
import 'blockly/javascript';
import '../blockly/customBlocks';
import '../blockly/generators';
import { toolboxXml } from '../blockly/toolbox';
import { emptyXml } from '../blockly/templates';
import { mflowFlatTheme } from '../blockly/theme';

Blockly.setLocale(zhHans);

const emit = defineEmits(['change', 'ready']);

const blocklyDiv = ref(null);
// 保存原始 Blockly 实例，避免被 Vue 响应式代理包裹
let workspace = null;
let suppressChange = false;

function ensureStartBlock() {
  const starts = workspace.getBlocksByType('start_flow', false);
  if (starts.length === 0) {
    // 恢复一个开始积木
    Blockly.Xml.domToWorkspace(Blockly.utils.xml.textToDom(emptyXml), workspace);
  }
  workspace.getBlocksByType('start_flow', false).forEach((b) => {
    b.setDeletable(false);
  });
}

function onChangeListener(event) {
  if (!workspace) return;
  // 阻止删除开始积木：删除后自动恢复
  if (event.type === Blockly.Events.BLOCK_DELETE) {
    ensureStartBlock();
  }
  if (suppressChange) return;
  emit('change');
}

onMounted(() => {
  workspace = markRaw(
    Blockly.inject(blocklyDiv.value, {
      toolbox: toolboxXml,
      theme: mflowFlatTheme,
      renderer: 'thrasos',
      trashcan: true,
      zoom: { controls: true, wheel: true, startScale: 1.0, maxScale: 3, minScale: 0.3, scaleSpeed: 1.2 },
      move: { scrollbars: true, drag: true, wheel: true },
      grid: { spacing: 22, length: 3, colour: '#e3e8ef', snap: true },
    })
  );

  loadXml(emptyXml);
  workspace.addChangeListener(onChangeListener);
  emit('ready');
});

onBeforeUnmount(() => {
  if (workspace) {
    workspace.dispose();
    workspace = null;
  }
});

function loadXml(xmlText) {
  suppressChange = true;
  try {
    workspace.clear();
    Blockly.Xml.domToWorkspace(Blockly.utils.xml.textToDom(xmlText), workspace);
    ensureStartBlock();
  } finally {
    suppressChange = false;
  }
}

// 加载 JSON（新 .m7p 格式）
function loadJson(state) {
  suppressChange = true;
  try {
    Blockly.serialization.workspaces.load(state, workspace);
    ensureStartBlock();
  } finally {
    suppressChange = false;
  }
}

// 从文本加载：自动识别 JSON / XML
function loadFromText(text) {
  const trimmed = (text || '').trim();
  if (!trimmed) {
    loadXml(emptyXml);
    return;
  }
  if (trimmed.startsWith('{')) {
    loadJson(JSON.parse(trimmed));
  } else {
    loadXml(trimmed);
  }
}

// 导出为 JSON 文本（.m7p）
function toJsonText() {
  return JSON.stringify(Blockly.serialization.workspaces.save(workspace), null, 2);
}

function getWorkspace() {
  return workspace;
}

defineExpose({ loadXml, loadJson, loadFromText, toJsonText, getWorkspace });
</script>

<template>
  <div ref="blocklyDiv" class="blockly-div"></div>
</template>

<style scoped>
.blockly-div {
  width: 100%;
  height: 100%;
}
</style>
