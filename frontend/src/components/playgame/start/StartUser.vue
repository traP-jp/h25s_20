<template>
  <div :class="$style.container">
    <ReadyIcon :id="props.name" :size="40" :isReady="props.isReady" />
    <div :class="$style.userName">{{ props.name }}</div>
    <StartBadge :gold="props.gold" :silver="props.silver" :bronze="props.bronze" />
  </div>
</template>
<script setup lang="ts">
import StartBadge from "@/components/playgame/start/StartBadge.vue";
import ReadyIcon from "@/components/ReadyIcon.vue";

const props = defineProps<{
  name: string;
  gold: number;
  silver: number;
  bronze: number;
  isReady: boolean;
}>();
</script>

<style module>
.userName {
  font-size: 16px;
  font-weight: bold;
  margin-right: 10px;
}

.container {
  width: 100%;
  display: flex;
  align-items: center;
}
</style>
