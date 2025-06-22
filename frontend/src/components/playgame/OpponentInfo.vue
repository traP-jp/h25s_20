<template>
  <div :class="$style.user">
    <UserIcon :id="id" :size="20" :class="$style.icon" />
    <div :class="$style.score">{{ score }}</div>
  </div>
</template>

<script setup lang="ts">
defineProps<{ id: string; score: number }>();
import UserIcon from "@/components/UserIcon.vue";
</script>

<style module>
.icon {
  outline: black solid 2px;
  margin-right: 8px;
}

.score {
  font-weight: bold;
  font-size: 12px;
}

.user {
  width: 100%;
  display: flex;
  align-items: center;
  margin-top: 8px;
}
</style>
