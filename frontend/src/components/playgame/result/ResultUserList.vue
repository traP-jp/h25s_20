<template>
  <div :class="$style.container">
    <ResultUser
      v-for="player in gameResultStore.players"
      :key="player.id"
      :id="player.id"
      :name="player.name"
      :rank="player.rank"
      :score="player.score"
    />
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
