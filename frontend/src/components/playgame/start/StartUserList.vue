<template>
  <div :class="$style.container">
    <StartUser
      v-for="player in players"
      :key="player.id"
      :id="player.id"
      :name="player.name"
      :gold="player.gold"
      :silver="player.silver"
      :bronze="player.bronze"
      :isReady="player.isReady"
    />
  </div>


</template>
<script setup lang="ts">
import StartUser from "@/components/playgame/start/StartUser.vue";
import { ref } from "vue";

const players = [
  { id: 0, name: "kitsne", gold: 3, silver: 2, bronze: 1, isReady: true },
  { id: 1, name: "mikannkann", gold: 2, silver: 1, bronze: 0, isReady: false },
  { id: 2, name: "rurun", gold: 1, silver: 0, bronze: 0, isReady: true },
];
</script>

<style module>
.container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}
</style>
