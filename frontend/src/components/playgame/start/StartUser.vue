<template>
  <div :class="$style.container">
    <ReadyIcon :id="props.player.name" :size="40" :isReady="props.player.isReady" />
    <div :class="$style.userName">{{ props.player.name }}</div>
    <StartBadge :gold="props.player.gold" :silver="props.player.silver" :bronze="props.player.bronze" />
  </div>
</template>
<script setup lang="ts">
import StartBadge from "@/components/playgame/start/StartBadge.vue";
import ReadyIcon from "@/components/ReadyIcon.vue";
import type { StartPlayer } from "@/lib/types";

const props = defineProps<{
  player: StartPlayer;
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
