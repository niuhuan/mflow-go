<script setup>
defineProps({
  title: { type: String, default: '' },
  width: { type: String, default: '460px' },
});
const emit = defineEmits(['close']);
</script>

<template>
  <div class="overlay" @click.self="emit('close')">
    <div class="card" :style="{ width }">
      <div class="header">
        <span class="title">{{ title }}</span>
        <button class="close" @click="emit('close')">×</button>
      </div>
      <div class="content">
        <slot />
      </div>
      <div v-if="$slots.footer" class="footer">
        <slot name="footer" />
      </div>
    </div>
  </div>
</template>

<style scoped>
.overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}
.card {
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.25);
  max-width: 90vw;
  max-height: 85vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid #eef0f3;
}
.title {
  font-weight: 600;
  font-size: 15px;
  color: #1a202c;
}
.close {
  border: none;
  background: transparent;
  font-size: 20px;
  line-height: 1;
  cursor: pointer;
  color: #718096;
}
.close:hover {
  color: #1a202c;
}
.content {
  padding: 16px;
  overflow-y: auto;
  text-align: left;
}
.footer {
  padding: 12px 16px;
  border-top: 1px solid #eef0f3;
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}
</style>
