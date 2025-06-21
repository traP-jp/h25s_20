<template>
  <div :class="$style.container">
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
</script>

<style module>
.container {
  width: 60px;
  height: 60px;
  border-radius: 5px;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #446688;
}

.text {
  color: white;
  font-size: 28px;
  font-weight: bold;
}

.icon {
  color: white;
  font-size: 32px;
}
</style>
