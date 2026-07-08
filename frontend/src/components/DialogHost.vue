<script setup>
import { watch, nextTick, ref } from 'vue';
import BaseModal from './modals/BaseModal.vue';
import { dialogState, dialogs } from '../composables/useDialogs';

const inputRef = ref(null);

watch(
  () => dialogState.visible,
  async (v) => {
    if (v && dialogState.type === 'prompt') {
      await nextTick();
      inputRef.value?.focus();
    }
  }
);

function cancel() {
  dialogs._finish(dialogState.type === 'prompt' ? null : dialogState.type === 'confirm' ? false : null);
}
function ok() {
  if (dialogState.type === 'prompt') dialogs._finish(dialogState.value);
  else if (dialogState.type === 'confirm') dialogs._finish(true);
  else dialogs._finish(true);
}
function choose(item) {
  dialogs._finish(item);
}
</script>

<template>
  <BaseModal v-if="dialogState.visible" :title="dialogState.title" width="380px" @close="cancel">
    <p v-if="dialogState.message" class="msg">{{ dialogState.message }}</p>

    <input
      v-if="dialogState.type === 'prompt'"
      ref="inputRef"
      v-model="dialogState.value"
      class="input"
      @keydown.enter="ok"
    />

    <div v-if="dialogState.type === 'select'" class="list">
      <div v-if="dialogState.items.length === 0" class="empty">（空）</div>
      <button
        v-for="(item, i) in dialogState.items"
        :key="i"
        class="list-item"
        @click="choose(item)"
      >
        {{ item }}
      </button>
    </div>

    <template #footer>
      <button v-if="dialogState.type !== 'select'" class="btn" @click="cancel">
        {{ dialogState.type === 'alert' ? '关闭' : '取消' }}
      </button>
      <button v-if="dialogState.type === 'prompt' || dialogState.type === 'confirm'" class="btn primary" @click="ok">
        确定
      </button>
      <button v-if="dialogState.type === 'select'" class="btn" @click="cancel">取消</button>
    </template>
  </BaseModal>
</template>

<style scoped>
.msg {
  margin: 0 0 12px;
  color: #2d3748;
  white-space: pre-wrap;
}
.input {
  width: 100%;
  box-sizing: border-box;
  padding: 8px 10px;
  border: 1px solid #cbd5e0;
  border-radius: 6px;
  font-size: 14px;
}
.list {
  display: flex;
  flex-direction: column;
  gap: 6px;
  max-height: 320px;
  overflow-y: auto;
}
.list-item {
  text-align: left;
  padding: 10px 12px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  background: #fff;
  cursor: pointer;
  font-size: 14px;
}
.list-item:hover {
  background: #f7fafc;
}
.empty {
  color: #a0aec0;
  text-align: center;
  padding: 16px;
}
.btn {
  padding: 6px 14px;
  border: 1px solid #cbd5e0;
  border-radius: 6px;
  background: #fff;
  cursor: pointer;
  font-size: 13px;
}
.btn:hover {
  background: #f7fafc;
}
.btn.primary {
  background: #2f855a;
  color: #fff;
  border-color: #2f855a;
}
.btn.primary:hover {
  background: #38a169;
}
</style>
