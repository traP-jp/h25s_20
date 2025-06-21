<template>
  <div :class="$style.container">
    <div :class="$style.header">Play View: Room{{ roomData[0].name }}</div>
    <div :class="$style.userinfo">
      <OpponentInfo
        v-for="player in players"
        :key="player.name"
        :icon="player.icon"
        :name="player.name"
        :score="player.score"
      />
    </div>

    <MyInfo icon="/images/player-self.png" name="Me" :score="50" :time="'10:00'" />

    <div :class="$style.board">
      <MainGameBoard v-model:board="board" />
    </div>

    <div :class="$style.inputbox">
      <MathInput v-model:board="board" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { roomData } from "@/lib/sample-data";

import type { Room } from "@/lib/types.ts";

import OpponentInfo from "@/components/playgame/OpponentInfo.vue";
import MainGameBoard from "@/components/playgame/MainGameBoard.vue";
import MathInput from "@/components/playgame/MathInput.vue";
import MyInfo from "@/components/playgame/MyInfo.vue";

import { ref, watch } from "vue";

const players = [
  { icon: "/images/player1.png", name: "Player 01", score: 30 },
  { icon: "/images/player2.png", name: "Player 02", score: 50 },
];

function handleRoomClick(room: Room) {
  console.log("Room clicked:", room);
}

const board = ref([1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8]);

watch(board, (newBoard) => {
  console.log("Board updated:", newBoard);
});
</script>

<style module>
.container {
  width: 500px;
  height: 100vh;
  margin: 0 auto;
  padding: 20px;
  border: 1px solid var(--border-color, #ccc);
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
}

.header {
  font-size: 24px;
  font-weight: bold;
  text-align: center;
  margin-bottom: 30px;
}

.rooms {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.myinfo {
  display: flex;
  flex-direction: column;
  border: 1px solid var(--border-color, #ccc);
  text-align: left;
}

.userinfo {
  display: flex;
  flex-direction: column;
  margin-bottom: 20px;
  border: 1px solid var(--border-color, #ccc);
}
</style>
