<template>
  <div :class="$style.container">
    <div :class="$style.main">
      <div :class="$style.statistics">
        <TextMark text="score" bgColor="#ffdd44" />
        <div :class="$style.score">{{ gameScore }}</div>
        <TextMark text="time" bgColor="#ff4400" />
        <div :class="$style.time">{{ formatTime(gameTime) }}</div>
      </div>
      <div :class="$style.right">
        <TextMark text="players" bgColor="#bb0000" :class="$style.playerMark" />
        <OpponentInfo v-for="player in players" :key="player.name" :id="player.name" :score="player.score" />
      </div>
    </div>
    <div :class="$style.board">
      <MainGameBoard v-model:board="board" :highlighted-numbers="highlightedNumbers" />
    </div>

    <div :class="$style.inputbox">
      <MathInput
        v-model:board="board"
        v-model:current-expression="currentExpression"
        v-model:currentRoom="currentRoom"
      />
    </div>

    <StartModal v-model:showStartModal="showStartModal" />
    <ResultModal v-model:showResultModal="showResultModal" v-model:showStartModal="showStartModal" />

    <TopBar :room="currentRoom" />
    <!-- Debug button to simulate countdown (replace with WebSocket callback in production) -->
    <button @click="debugStartCountdown(3)">Debug Countdown</button>
    <CountDown v-if="countdown >= 0" :time="countdown" />
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from "vue";
import { useWebSocketStore, useGameResultStore, useRoomPlayersStore, useCurrentRoomStore } from "@/store";
import TopBar from "@/components/playgame/TopBar.vue";
import type { Room } from "@/lib/types";

import OpponentInfo from "@/components/playgame/OpponentInfo.vue";
import MainGameBoard from "@/components/playgame/MainGameBoard.vue";
import MathInput from "@/components/playgame/MathInput.vue";

import StartModal from "@/components/playgame/start/StartModal.vue";
import ResultModal from "@/components/playgame/result/ResultModal.vue";
import CountDown from "@/components/playgame/CountDown.vue";
import TextMark from "@/components/TextMark.vue";

import { onMounted, onBeforeUnmount } from "vue";
import { useRoute } from "vue-router";

// ルーターから情報を取得
const route = useRoute();

// WebSocketストアを取得
const webSocketStore = useWebSocketStore();
const gameResultStore = useGameResultStore();
const roomPlayersStore = useRoomPlayersStore();
const currentRoomStore = useCurrentRoomStore();

// ゲーム状態
const gameScore = ref(0);
const gameTime = ref(60); // 初期時間60秒
const gameStarted = ref(false);
const countdown = ref(-1); // -1 means hide the countdown screen

// 現在のユーザーIDを取得（仮実装）
function getCurrentUserId(): number {
  // 実際の実装では認証情報から取得
  return 1;
}

// ゲームタイマー管理
let gameTimer: number | null = null;

function startGameTimer() {
  if (gameTimer) return;

  gameTimer = setInterval(() => {
    if (gameTime.value > 0) {
      gameTime.value--;
    } else {
      stopGameTimer();
    }
  }, 1000);
}

function stopGameTimer() {
  if (gameTimer) {
    clearInterval(gameTimer);
    gameTimer = null;
  }
}

// 盤面更新をサーバーに送信
function sendBoardUpdate(newBoard: number[], gainScore: number) {
  if (!currentRoom.value) return;

  const boardUpdateEvent = {
    event: "board_update",
    content: {
      room_id: currentRoom.value.roomId,
      user_id: getCurrentUserId(),
      board: {
        content: newBoard,
        version: Date.now(), // 簡易的なバージョン管理
        size: 16,
      },
      gain_score: gainScore,
    },
  };

  // グローバルWebSocketストア経由で送信
  webSocketStore.sendMessage(boardUpdateEvent);
}

// ルーターから渡された情報を取得
onMounted(() => {
  // まずストアからルーム情報を取得
  let room = currentRoomStore.getCurrentRoom();

  if (!room) {
    // ストアにない場合はクエリパラメータから復元
    const roomId = route.params.roomId as string;
    const roomName = route.query.roomName as string;
    const isOpened = route.query.isOpened === "true";

    if (roomId && roomName) {
      room = {
        roomId: parseInt(roomId),
        roomName,
        isOpened,
        users: [],
      };
      // ストアに保存
      currentRoomStore.setCurrentRoom(room);
    }
  }

  currentRoom.value = room;
  console.log("Current room:", currentRoom.value);

  // 各ストアの初期化
  // ゲーム結果ストアの初期化
  gameResultStore.clearPlayers();

  // ルームプレイヤーストアの初期化
  roomPlayersStore.clearPlayers();

  // 現在のルーム情報があれば、プレイヤー情報も初期化
  if (room?.users && room.users.length > 0) {
    // Room型のusersをroomPlayersStoreが期待する形式に変換
    const roomPlayers = room.users.map((user) => ({
      user_id: 0,
      user_name: user.username, // Room型ではnameがないため、idをnameとして使用
      is_ready: user.isReady,
    }));
    roomPlayersStore.updatePlayers(roomPlayers);
  }

  // グローバルWebSocketストアに現在のコンポーネントのイベントハンドラーを設定
  // 既存のWebSocket接続がない場合は、ローカルストレージからユーザー名を取得して接続
});

onBeforeUnmount(() => {
  stopGameTimer();
  // グローバルWebSocketは他のコンポーネントでも使用される可能性があるため、
  // ここでは接続を切断しない
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

const showStartModal = ref(true);
const showResultModal = ref(false);
const currentRoom = ref<Room | null>(null);

// 盤面の変更を監視してWebSocketで送信
watch(board, (newBoard: number[], oldBoard: number[]) => {
  console.log("Board updated:", newBoard);

  // ゲームが開始されていて、実際に盤面が変更された場合のみ送信
  if (gameStarted.value && oldBoard && JSON.stringify(newBoard) !== JSON.stringify(oldBoard)) {
    // スコア計算（簡易実装）
    const gainScore = calculateScoreGain(newBoard, oldBoard);
    sendBoardUpdate(newBoard, gainScore);
  }
});

// スコア計算（簡易実装）
function calculateScoreGain(newBoard: number[], oldBoard: number[]): number {
  // 実際のゲームロジックに応じてスコアを計算
  // ここでは例として、変更された数の数をスコアとする
  let changes = 0;
  for (let i = 0; i < newBoard.length; i++) {
    if (newBoard[i] !== oldBoard[i]) {
      changes++;
    }
  }
  return changes * 10; // 変更1つあたり10点
}

// 時間をMM:SS形式でフォーマット
function formatTime(seconds: number): string {
  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = seconds % 60;
  return `${minutes.toString().padStart(2, "0")}:${remainingSeconds.toString().padStart(2, "0")}`;
}

// デバッグ関数
function debugStartGame() {
  gameStarted.value = true;
  gameTime.value = 60;
  startGameTimer();
  console.log("Debug: Game started");
}

// デバッグ用カウントダウン
async function debugStartCountdown(startNum: number) {
  for (let i = startNum; i > 0; i--) {
    countdown.value = i;
    await new Promise((resolve) => setTimeout(resolve, 1000));
  }
  countdown.value = -1;
  debugStartGame(); // カウントダウン後にゲーム開始
}

const currentExpression = ref("");

// 現在の数式に含まれる数字を抽出してハイライト対象を決定
const highlightedNumbers = computed(() => {
  if (!currentExpression.value) return [];

  // 数式から数字のみを抽出（演算子や括弧を除外）
  const numbersInExpression = currentExpression.value.match(/[1-9]/g) || [];

  // 重複を除去して数値に変換
  return [...new Set(numbersInExpression.map(Number))];
});

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
