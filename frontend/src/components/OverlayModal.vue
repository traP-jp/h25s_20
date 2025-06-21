<template>
  <div class="overlay" @click.self="close">
    <div class="modal">
      <slot></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
const emit = defineEmits<{ (e: 'close'): void }>();
function close() {
  emit('close');
}
</script>

<style scoped>
.overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0,  0,  0,  0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}
.modal {
  background: #fff;
  padding: 1.5rem;
  border-radius: 8px;
  max-width: 90%;
  max-height: 90%;
  overflow: auto;
}
</style>