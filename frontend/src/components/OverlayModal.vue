<template>
  <div :class="$style.overlay" @click.self="close">
    <div :class="$style.modal">
      <slot></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
const emit = defineEmits<{ (e: "close"): void }>();
function close() {
  emit("close");
}
</script>

<style module>
.overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 3;
}
.modal {
  margin-top: 40px;
  background: #fff;
  padding: 1.5rem;
  width: 100%;
  height: 100%;
  overflow: auto;
}
</style>
