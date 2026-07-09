<script setup>
import { ref, watch, nextTick } from 'vue';

const props = defineProps({
  messages: { type: Array, default: () => [] },
  position: { type: String, default: 'right' }, // 'right' | 'bottom'
});

const emit = defineEmits(['toggle-position']);

const container = ref(null);

watch(
  () => props.messages.length,
  async () => {
    await nextTick();
    if (container.value) {
      container.value.scrollTop = container.value.scrollHeight;
    }
  }
);
</script>

<template>
  <div class="console-panel">
    <div class="console-header">
      <span class="title">控制台</span>
      <button
        class="toggle"
        :title="position === 'right' ? '切换到窗口下方' : '切换到窗口右侧'"
        @click="emit('toggle-position')"
      >
        <!-- 右侧时：方框内左下箭头（点击移到下方） -->
        <svg v-if="position === 'right'" viewBox="0 0 24 24" width="16" height="16"
          fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <rect x="3" y="3" width="18" height="18" rx="2" />
          <line x1="15" y1="9" x2="9" y2="15" />
          <polyline points="13 15 9 15 9 11" />
        </svg>
        <!-- 下方时：方框内右上箭头（点击移到右侧） -->
        <svg v-else viewBox="0 0 24 24" width="16" height="16"
          fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <rect x="3" y="3" width="18" height="18" rx="2" />
          <line x1="9" y1="15" x2="15" y2="9" />
          <polyline points="11 9 15 9 15 13" />
        </svg>
      </button>
    </div>
    <div ref="container" class="console">
      <p v-for="(msg, i) in messages" :key="i" class="console-line">{{ msg }}</p>
    </div>
  </div>
</template>

<style scoped>
.console-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
  width: 100%;
  background: #1e1e1e;
}
.console-header {
  flex: 0 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 4px 8px;
  background: #252526;
  border-bottom: 1px solid #333;
}
.console-header .title {
  color: #cccccc;
  font-size: 12px;
}
.toggle {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: transparent;
  color: #9da5b4;
  cursor: pointer;
  padding: 2px;
  border-radius: 4px;
}
.toggle:hover {
  color: #fff;
  background: #3a3d41;
}
.console {
  flex: 1 1 auto;
  overflow-y: auto;
  color: #d4d4d4;
  font-family: 'Cascadia Code', 'Consolas', monospace;
  font-size: 12px;
  padding: 8px 12px;
  box-sizing: border-box;
  text-align: left;
}
.console-line {
  margin: 0 0 2px;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
