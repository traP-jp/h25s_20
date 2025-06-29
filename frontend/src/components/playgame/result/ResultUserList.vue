<template>
  <div :class="$style.container">
    <ResultUser v-for="player in gameResultStore.players" :player="player" />
  </div>
</template>
<script setup lang="ts">
import ResultUser from "@/components/playgame/result/ResultUser.vue";
import { useGameResultStore } from "@/store";

const gameResultStore = useGameResultStore();
</script>

<style module>
.container {
  width: 320px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}
</style>
