<template>
  <div :class="$style.container">
    <ResultUser
      v-for="player in players"
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
// import { ref } from "vue";

const players = [
  { id: 0, name: "kitsne", rank: 1, score: 100 },
  { id: 1, name: "mikannkann", rank: 2, score: 80 },
  { id: 2, name: "rurun", rank: 3, score: 60 },
];
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
