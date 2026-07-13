<script setup>
import { ref, watch, onMounted, onBeforeUnmount, nextTick } from 'vue';

const props = defineProps({
  show: { type: Boolean, default: false },
  anchor: { type: Object, default: null },
  minWidth: { type: Number, default: 200 },
});

const emit = defineEmits(['close']);

const menuRef = ref(null);
const menuStyle = ref({});

function getAnchorEl() {
  const a = props.anchor;
  if (!a) return null;
  if (a instanceof HTMLElement) return a;
  return a.value ?? null;
}

function updatePosition() {
  const el = getAnchorEl();
  if (!el || !props.show) return;
  const rect = el.getBoundingClientRect();
  menuStyle.value = {
    top: `${rect.bottom + 6}px`,
    right: `${window.innerWidth - rect.right}px`,
    minWidth: `${Math.max(props.minWidth, rect.width)}px`,
  };
}

function onKeydown(e) {
  if (e.key === 'Escape' && props.show) emit('close');
}

watch(
  () => props.show,
  async (v) => {
    if (v) {
      await nextTick();
      updatePosition();
    }
  }
);

watch(() => props.anchor, updatePosition);

onMounted(() => {
  window.addEventListener('resize', updatePosition);
  window.addEventListener('scroll', updatePosition, true);
  window.addEventListener('keydown', onKeydown);
});

onBeforeUnmount(() => {
  window.removeEventListener('resize', updatePosition);
  window.removeEventListener('scroll', updatePosition, true);
  window.removeEventListener('keydown', onKeydown);
});
</script>

<template>
  <Teleport to="body">
    <Transition name="popup">
      <div v-if="show" class="popup-backdrop" @click="emit('close')">
        <div ref="menuRef" class="popup-menu" :style="menuStyle" @click.stop>
          <slot />
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.popup-backdrop {
  position: fixed;
  inset: 0;
  z-index: 1000;
  background: transparent;
}

.popup-menu {
  position: fixed;
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  box-shadow:
    0 4px 6px -1px rgba(0, 0, 0, 0.08),
    0 10px 24px -4px rgba(0, 0, 0, 0.14);
  padding: 6px 0;
  max-height: min(70vh, 520px);
  overflow-y: auto;
  transform-origin: top right;
}

.popup-enter-active {
  transition: opacity 0.18s ease;
}
.popup-enter-active .popup-menu {
  transition:
    opacity 0.18s ease,
    transform 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}
.popup-leave-active {
  transition: opacity 0.12s ease;
}
.popup-leave-active .popup-menu {
  transition:
    opacity 0.12s ease,
    transform 0.12s ease;
}
.popup-enter-from,
.popup-leave-to {
  opacity: 0;
}
.popup-enter-from .popup-menu,
.popup-leave-to .popup-menu {
  opacity: 0;
  transform: translateY(-8px) scale(0.96);
}
</style>
