<template>
  <div :class="$style.container">
    <ResultBadge :rank="props.player.rank" :class="$style.rank" />
    <UserIcon :id="props.player.name" :size="40" />
    <div :class="$style.userName">{{ props.player.name }}</div>
    <div :class="$style.score">{{ props.player.score }}点</div>
  </div>
</template>
<script setup lang="ts">
// import { ref } from "vue";
import ResultBadge from "@/components/playgame/result/ResultBadge.vue";
import UserIcon from "@/components/UserIcon.vue";
import type { ResultPlayer } from "@/lib/types";

const props = defineProps<{
  player: ResultPlayer;
}>();
</script>

<style module>
.container {
  width: 100%;
  display: flex;
  align-items: center;
}

.rank {
  margin-top: 5px;
  margin-right: 10px;
}

.userName {
  margin: 10px;
  font-weight: bold;
}
</style>
