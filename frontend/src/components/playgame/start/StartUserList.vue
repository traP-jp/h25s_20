<template>
  <div :class="$style.container">
    <StartUser
      v-for="player in roomPlayersStore.players"
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
