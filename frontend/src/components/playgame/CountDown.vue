<template>
  <div v-if="time !== -1" class="overlay">
    <div class="count">{{ time }}</div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, toRef } from "vue";
const props = defineProps<{ time: number }>();
const time = toRef(props, "time");
</script>

<style scoped>
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
  z-index: 1005; /* ゲーム画面より上で、タイトルより下 */
}
.count {
  color: #fff;
  font-size: 15rem;
  font-weight: bold;
}
</style>
