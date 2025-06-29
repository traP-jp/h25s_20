<template>
  <div :class="$style.container">
    <StartUser v-for="player in roomPlayersStore.players" :player="player" />
  </div>
</template>

<script setup lang="ts">
import StartUser from "@/components/playgame/start/StartUser.vue";
import { useRoomPlayersStore } from "@/store";

const roomPlayersStore = useRoomPlayersStore();
</script>

<style module>
.container {
  margin-top: 40px;
  display: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}
</style>
