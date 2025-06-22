<template>
  <div :class="$style.content" :style="{ backgroundColor: bgColor }">{{ text }}</div>
</template>

<script setup lang="ts">
import { defineProps } from "vue";
defineProps<{
  text: string;
  bgColor: string;
}>();
</script>

<style module>
.content {
  color: white;
  padding: 0 10px;
  border-radius: 1em;
  font-weight: bold;
  font-size: 12px;
}
</style>
