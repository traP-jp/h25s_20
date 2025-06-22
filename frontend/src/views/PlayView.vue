<template>
  <div :class="$style.container">
    <TopBar v-model:room="roomData[0]" />
    <div :class="$style.userinfo">
      <OpponentInfo
        v-for="player in players"
        :key="player.name"
        :icon="player.icon"
        :name="player.name"
        :score="player.score"
      />
    </div>

    <ScoreInfo icon="/images/player-self.png" name="Me" :score="50" :time="'10:00'" :class="$style.myinfo" />

    <div :class="$style.board">
      <MainGameBoard v-model:board="board" />
    </div>

    <div :class="$style.inputbox">
      <MathInput v-model:board="board" />
    </div>

    <StartModal />
    <ResultModal />
  </div>
</template>

<script setup lang="ts">
import { ref, watch, provide } from "vue";
import { roomData } from "@/lib/sample-data";
import TopBar from "@/components/playgame/TopBar.vue";

import OpponentInfo from "@/components/playgame/OpponentInfo.vue";
import MainGameBoard from "@/components/playgame/MainGameBoard.vue";
import MathInput from "@/components/playgame/MathInput.vue";
import MyInfo from "@/components/playgame/MyInfo.vue";
import StartModal from "@/components/playgame/start/StartModal.vue";
import ResultModal from "@/components/playgame/result/ResultModal.vue";

const players = [
  { icon: "/images/player1.png", name: "Player 01", score: 30 },
  { icon: "/images/player2.png", name: "Player 02", score: 50 },
];

const showStartModal = ref(false);
const showResultModal = ref(false);

provide("showStartModal", showStartModal);
provide("showResultModal", showResultModal);

const board = ref([1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8]);

watch(board, (newBoard) => {
  console.log("Board updated:", newBoard);
});
</script>

<style module>
.container {
  width: 360px;
  height: 100vh;
  margin: 0 auto;
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
  z-index: 1010;
}

.userinfo {
  display: flex;
  flex-direction: column;
  margin-bottom: 20px;
  border: 1px solid var(--border-color, #ccc);
}
</style>
