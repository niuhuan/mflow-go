<script setup>
import { ref, watch, nextTick } from 'vue';

const props = defineProps({
  messages: { type: Array, default: () => [] },
});

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
  <div ref="container" class="console">
    <p v-for="(msg, i) in messages" :key="i" class="console-line">{{ msg }}</p>
  </div>
</template>

<style scoped>
.console {
  height: 100%;
  overflow-y: auto;
  background: #1e1e1e;
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
