<template>
  <div :class="$style.container" @click="$emit('click')">
    <Icon v-if="icon" :icon="icon" :class="$style.icon" />
    <span v-else :class="$style.text">{{ text }}</span>
  </div>
</template>

<script setup lang="ts">
import { Icon } from "@iconify/vue";

defineProps<{
  text?: string;
  icon?: string;
}>();

defineEmits<{
  click: [];
}>();
</script>

<style module>
.container {
  width: 50px;
  height: 50px;
  border-radius: 5px;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f0f0f0;
  cursor: pointer;
  transition: all 0.2s ease;
}

.container:hover {
  background-color: #88aaff;
  transform: scale(1.05);
}

.container:active {
  transform: scale(0.95);
}

.text {
  color: black;
  font-size: 28px;
  font-weight: bold;
}

.icon {
  color: black;
  font-size: 32px;
}
</style>
