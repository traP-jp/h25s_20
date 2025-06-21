<template>
  <div :class="$style.container">
    <div>
      <button
        :disabled="room.status === 'closed'"
        @click="emit('click')"
        :class="[
          $style.button,
          room.status === 'closed' ? $style.closed : $style.open,
        ]"
      >
        <div>{{ room.name }}</div>
        <div>ID: {{ room.id }}</div>
        <div>Status: {{ room.status }}</div>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from "vue";
import type { Room } from "@/lib/sample-data";

defineProps<{
  room: Room;
}>();

const emit = defineEmits<{
  click: [];
}>();
</script>

<style module>
.container {
  width: 100%;
  border: 1px solid #ccc;
  border-radius: 4px;
  text-align: left;
  cursor: pointer;
}
.button {
  width: 100%;
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  font-size: 16px;
  transition: background-color 0.3s ease;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  text-align: left;
}
.open {
  background-color: #4caf50;
  cursor: pointer;
}
.closed {
  background-color: #ccc;
  opacity: 0.4;
  color: #888;
  cursor: not-allowed;
}
</style>
