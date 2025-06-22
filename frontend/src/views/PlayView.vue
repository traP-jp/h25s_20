<template>
  <div :class="$style.container">
    <TopBar :room="currentRoom || roomData[0]" />
    <div :class="$style.main">
      <div :class="$style.statistics">
        <TextMark text="score" bgColor="#ffdd44" />
        <div :class="$style.score">30</div>
        <TextMark text="time" bgColor="#ff4400" />
        <div :class="$style.time">01:00</div>
      </div>
      <div :class="$style.right">
        <TextMark text="players" bgColor="#bb0000" :class="$style.playerMark" />
        <OpponentInfo v-for="player in players" :key="player.name" :id="player.name" :score="player.score" />
      </div>
    </div>
    <div :class="$style.board">
      <MainGameBoard v-model:board="board" />
    </div>

    <div :class="$style.inputbox">
      <MathInput v-model:board="board" />
    </div>

    <StartModal />
    <ResultModal />
    <!-- Debug button to simulate countdown (replace with WebSocket callback in production) -->
    <button @click="debugStartCountdown(3)">Debug Countdown</button>
    <CountDown v-if="countdown >= 0" :time="countdown" />
  </div>
</template>

<script setup lang="ts">
import { ref, watch, provide, onMounted } from "vue";
import { useRoute } from "vue-router";
import { roomData } from "@/lib/sample-data";
import TopBar from "@/components/playgame/TopBar.vue";
import type { Room } from "@/lib/types";

import OpponentInfo from "@/components/playgame/OpponentInfo.vue";
import MainGameBoard from "@/components/playgame/MainGameBoard.vue";
import MathInput from "@/components/playgame/MathInput.vue";

import StartModal from "@/components/playgame/start/StartModal.vue";
import ResultModal from "@/components/playgame/result/ResultModal.vue";
import CountDown from "@/components/playgame/CountDown.vue";
import TextMark from "@/components/TextMark.vue";

// ルーターから情報を取得
const route = useRoute();
const currentRoom = ref<Room | null>(null);

// ルーターから渡された情報を取得
onMounted(() => {
  // ルートのstateからroom情報を取得
  if (history.state && history.state.room) {
    currentRoom.value = history.state.room;
  } else {
    // state経由でない場合はクエリパラメータから復元
    const roomId = route.params.roomId as string;
    const roomName = route.query.roomName as string;
    const isOpened = route.query.isOpened === "true";

    if (roomId && roomName) {
      currentRoom.value = {
        roomId: parseInt(roomId),
        roomName,
        isOpened,
        users: [],
      };
    }
  }

  console.log("Current room:", currentRoom.value);
});

// 初期盤面を生成する関数
function generateInitialBoard(): number[] {
  // 1-9の数字をランダムに4つずつ選んで16個の配列を作成
  const numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9];
  const board: number[] = [];

  // 各数字から4つまでランダムに選択
  for (let i = 0; i < 16; i++) {
    const randomIndex = Math.floor(Math.random() * numbers.length);
    board.push(numbers[randomIndex]);
  }

  return board;
}

const board = ref(generateInitialBoard());

// プレイヤー情報（後でWebSocketから取得する予定）
const players = [
  { icon: "/images/player1.png", name: "Player 01", score: 30 },
  { icon: "/images/player2.png", name: "Player 02", score: 50 },
];

const showStartModal = ref(false);
const showResultModal = ref(false);
const countdown = ref(-1); // -1 means hide the countdown screen

provide("showStartModal", showStartModal);
provide("showResultModal", showResultModal);

// デバック用 実際はwebsocketのコールバックからカウントダウンの更新をする
async function debugStartCountdown(startNum: number) {
  for (let i = startNum; i > 0; i--) {
    countdown.value = i;
    await new Promise((resolve) => setTimeout(resolve, 1000));
  }
  countdown.value = -1;
}

watch(board, (newBoard: number[]) => console.log("Board updated:", newBoard));
</script>

<style module>
.container {
  width: 360px;
  height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
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

.right {
  position: fixed;
  width: 60px;
  top: 80px;
  right: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.playerMark {
  position: absolute;
  right: 10px;
  top: -20px;
}

.main {
  margin-top: 40px;
  position: relative;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
}

.left {
  width: 260px;
}

.statistics {
  margin-top: 14px;
  height: 60px;
  gap: 10px;
  display: flex;
  align-items: center;
  justify-content: space-around;
}

.score {
  font-size: 40px;
  font-weight: bold;
}

.time {
  font-size: 20px;
  font-weight: bold;
}

.board {
  padding: 14px;
}
</style>
