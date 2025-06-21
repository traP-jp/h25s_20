<template>
  <div :class="$style.container">
    <div>
      <div>{{ room.name }}</div>
      <div :class="$style.details">
        <div>ID: {{ room.id }}</div>
        <div>Status: {{ room.status }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps } from "vue";

const props = defineProps<{
  room: {
    id: number;
    name: string;
    status: string;
  };
}>();
</script>

<style module>
.container {
  width: 100%;
  border: 1px solid #ccc;
  text-align: left;
}
.details {
  padding: 10px;
  border-top: 1px solid #ddd;
}
</style>
