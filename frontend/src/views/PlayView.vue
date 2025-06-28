<template>
  <div :class="$style.container">
    <div :class="$style.main">
      <div :class="$style.statistics">
        <TextMark text="score" bgColor="#ffdd44" />
        <div :class="$style.score">{{ gameScore }}</div>
        <TextMark text="time" bgColor="#ff4400" />
        <div
          :class="[
            $style.time,
            getTimeClass(gameTime) && $style[getTimeClass(gameTime)],
          ]"
        >
          {{ isTimePlaceholder(gameTime) ? "xx:xx" : formatTime(gameTime) }}
        </div>
      </div>
      <div :class="$style.right">
        <TextMark text="players" bgColor="#bb0000" :class="$style.playerMark" />
        <OpponentInfo
          v-for="player in Array.from(playerScores.values()).sort(
            (a, b) => b.score - a.score
          )"
          :key="player.name"
          :id="player.name"
          :score="player.score"
        />
      </div>
    </div>
    <div :class="$style.board">
      <MainGameBoard
        v-model:board="board"
        :highlighted-numbers="highlightedNumbers"
        :is-disabled="countdown > 0"
      />
    </div>

    <div :class="$style.inputbox">
      <MathInput
        v-model:version="version"
        v-model:board="board"
        v-model:current-expression="currentExpression"
        v-model:currentRoom="currentRoom"
        v-model:expression="expression"
      />
    </div>

    <StartModal v-model:showStartModal="showStartModal" />
    <ResultModal
      v-model:showResultModal="showResultModal"
      v-model:showStartModal="showStartModal"
    />

    <TopBar :room="currentRoom" />
    <!-- Debug button to simulate countdown (replace with WebSocket callback in production) -->
    <!-- <button @click="debugStartCountdown(3)">Debug Countdown</button> -->
    <CountDown
      v-if="countdown > 0 && !gameStarted"
      :time="countdown"
      :isStartCountdown="true"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from "vue";
import {
  useWebSocketStore,
  useGameResultStore,
  useRoomPlayersStore,
  useCurrentRoomStore,
} from "@/store";
import TopBar from "@/components/playgame/TopBar.vue";
import type { Room } from "@/lib/types";
import { WS_EVENTS, type WebSocketEvent } from "@/lib/websocket";

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
const gameTime = ref(-1); // 初期状態は-1（プレースホルダー表示）
const gameStarted = ref(false);
const countdown = ref(-1); // -1 means hide the countdown screen
const expression = ref("");

// 各プレイヤーのリアルタイムスコアを追跡
const playerScores = ref<Map<string, { name: string; score: number }>>(
  new Map()
);

const version = ref(0); // フォーミュラのバージョン管理

// ゲームタイマー管理
let gameTimer: number | null = null;

function startGameTimer() {
  if (gameTimer) return;

  gameTimer = setInterval(() => {
    if (gameTime.value > 0) {
      gameTime.value--;
      // 時間が0になってもまだゲームが終了していない場合は継続
      // サーバーからのGAME_ENDEDイベントを待つ
    } else if (gameTime.value === 0) {
      // 0秒になったら一旦タイマーを停止するが、ゲーム状態は維持
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
      user_name: webSocketStore.currentUsername,
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
      user_name: user.username,
      is_ready: user.isReady,
    }));
    roomPlayersStore.updatePlayers(roomPlayers);
  }

  // WebSocketイベントハンドラーの設定
  const handleWebSocketEvent = (event: WebSocketEvent) => {
    console.log("PlayView received WebSocket event:", event);

    switch (event.event) {
      case WS_EVENTS.PLAYER_JOINED:
        console.log("Player joined event received:", event.content);
        if (
          event.content &&
          typeof event.content === "object" &&
          "room" in event.content
        ) {
          const roomContent = event.content as any;
          if (roomContent.room && roomContent.room.players) {
            // ルームプレイヤーストアを更新
            const roomPlayers = roomContent.room.players.map((player: any) => ({
              user_name: player.user_name,
              is_ready: player.is_ready,
            }));
            roomPlayersStore.updatePlayers(roomPlayers);
            console.log(
              "Updated room players after PLAYER_JOINED:",
              roomPlayers
            );
          }
        }
        break;

      case WS_EVENTS.PLAYER_READY:
        console.log("Player ready event received:", event.content);
        if (
          event.content &&
          typeof event.content === "object" &&
          "user_name" in event.content
        ) {
          const playerContent = event.content as any;
          roomPlayersStore.setPlayerReady(playerContent.user_name, true);
          console.log(`Player ${playerContent.user_name} is now ready`);
        }
        break;

      case WS_EVENTS.PLAYER_CANCELED:
        console.log("Player canceled event received:", event.content);
        if (
          event.content &&
          typeof event.content === "object" &&
          "user_name" in event.content
        ) {
          const playerContent = event.content as any;
          roomPlayersStore.setPlayerReady(playerContent.user_name, false);
          console.log(`Player ${playerContent.user_name} canceled ready state`);
        }
        break;

      case WS_EVENTS.PLAYER_LEFT:
        console.log("Player left event received:", event.content);
        if (
          event.content &&
          typeof event.content === "object" &&
          "room" in event.content
        ) {
          const roomContent = event.content as any;
          if (roomContent.room && roomContent.room.players) {
            // ルームプレイヤーストアを更新
            const roomPlayers = roomContent.room.players.map((player: any) => ({
              user_name: player.user_name,
              is_ready: player.is_ready,
            }));
            roomPlayersStore.updatePlayers(roomPlayers);
            console.log("Updated room players after PLAYER_LEFT:", roomPlayers);
          }
        }
        break;

      case WS_EVENTS.COUNTDOWN_START:
        console.log("Countdown start event received");
        showStartModal.value = false; // スタートモーダルを閉じてゲーム画面に移行
        // ゲーム開始準備中は時間をプレースホルダー状態に
        gameTime.value = -1;
        gameStarted.value = false; // まだゲームは開始していない
        countdown.value = -1; // カウントダウンをリセット
        console.log("Countdown start - gameStarted:", gameStarted.value, "gameTime:", gameTime.value);
        break;

      case WS_EVENTS.COUNTDOWN:
        console.log("Countdown event received:", event.content, "gameStarted:", gameStarted.value);
        if (
          event.content &&
          typeof event.content === "object" &&
          "count" in event.content
        ) {
          const countValue = (event.content as any).count;
          // ゲーム開始前のカウントダウンのみ処理
          // ゲームが開始されている場合は完全に無視
          if (!gameStarted.value) {
            countdown.value = countValue;
            console.log("Setting pre-game countdown:", countValue);
          } else {
            console.log("Ignoring countdown during game:", countValue);
          }
        }
        break;

      case WS_EVENTS.GAME_START:
        console.log("Game start event received");
        countdown.value = -1; // カウントダウンを非表示
        showStartModal.value = false; // スタートモーダルを閉じる
        gameStarted.value = true; // ゲーム開始フラグを設定
        gameTime.value = 120; // 120秒ゲーム
        version.value = 0; // バージョンをリセット
        expression.value = ""; // 数式をリセット

        console.log("Game started - gameStarted:", gameStarted.value, "gameTime:", gameTime.value);

        startGameTimer();

        console.log("players:", roomPlayersStore.players);

        // プレイヤースコアを初期化
        playerScores.value.clear();
        // 現在のルームのプレイヤー情報からスコアを初期化
        for (const player of roomPlayersStore.players) {
          playerScores.value.set(player.name, {
            name: player.name,
            score: 0,
          });
        }

        console.log("playerScores initialized:", playerScores.value);

        // 現在のユーザーも確実に追加
        // const currentUserId = getCurrentUserId();
        // if (currentUserId && !playerScores.value.has(currentUserId)) {
        //   const currentPlayer = roomPlayersStore.players.find((p) => parseInt(p.id) === currentUserId);
        //   playerScores.value.set(currentUserId, {
        //     name: currentPlayer?.name || `Player ${currentUserId}`,
        //     score: 0,
        //   });
        // }

        console.log("Initialized player scores:", playerScores.value);

        // ボード情報があれば更新
        if (
          event.content &&
          typeof event.content === "object" &&
          "board" in event.content
        ) {
          const boardContent = event.content as any;
          if (boardContent.board && boardContent.board.content) {
            board.value = boardContent.board.content;
          }
        }
        break;

      case WS_EVENTS.BOARD_UPDATED:
        console.log("Board updated event received:", event.content);
        if (
          event.content &&
          typeof event.content === "object" &&
          "board" in event.content
        ) {
          const boardContent = event.content as any;
          if (boardContent.board && boardContent.board.content) {
            console.log(
              "Updating board from WebSocket:",
              boardContent.board.content
            );
            board.value = boardContent.board.content;
          }

          // プレイヤーのスコアを更新
          if (boardContent.gain_score && boardContent.user_name) {
            const userName = boardContent.user_name;
            const gainScore = boardContent.gain_score;

            // 全プレイヤーのスコアマップを更新
            if (playerScores.value.has(userName)) {
              const playerData = playerScores.value.get(userName)!;
              playerData.score += gainScore;
              console.log(
                `Updated score for ${playerData.name}: ${playerData.score} (+${gainScore})`
              );
            } else {
              // 新しいプレイヤーの場合、追加
              playerScores.value.set(userName, {
                name: userName,
                score: gainScore,
              });
              console.log(
                `Added new player ${userName} with score: ${gainScore}`
              );
              console.log("Current player scores:", playerScores.value);
            }

            // 自分のスコア更新ログ
            const currentUsername = webSocketStore.currentUsername;
            if (userName === currentUsername) {
              console.log(
                "Score updated (own submission):",
                gameScore.value,
                "gained:",
                gainScore
              );
            }
          }

          version.value = boardContent.board.version;
        }
        break;

      case WS_EVENTS.GAME_ENDED:
        console.log("Game ended event received");
        
        // ゲーム状態を即座にリセット（これにより以降のCOUNTDOWNイベントは無視される）
        gameStarted.value = false;
        countdown.value = -1; // カウントダウンを非表示
        
        // ゲーム終了時は時間を0で固定（プレースホルダーにリセットしない）
        if (gameTime.value > 0) {
          gameTime.value = 0; // 残り時間がある場合は0に設定
        }
        stopGameTimer();
        
        console.log("Game state reset - gameStarted:", gameStarted.value, "gameTime:", gameTime.value);
        
        showResultModal.value = true;

        // 蓄積されたプレイヤースコア情報をgameResultStoreに反映
        const finalScores = Array.from(playerScores.value.entries()).map(
          ([userName, playerData]) => ({
            user_name: userName,
            score: playerData.score,
          })
        );

        if (finalScores.length > 0) {
          gameResultStore.updatePlayers(finalScores);
          console.log(
            "Updated game result store with tracked scores:",
            finalScores
          );
        } else {
          console.warn("No player scores tracked during the game");
        }

        // 全てのプレイヤーの isReady を false に設定
        for (const player of roomPlayersStore.players) {
          player.isReady = false;
        }
        break;
    }
  };

  // グローバルWebSocketストアにイベントハンドラーを追加
  if (webSocketStore.wsManager) {
    webSocketStore.wsManager.addMessageHandler(handleWebSocketEvent);
    console.log("Added PlayView WebSocket event handler");
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
  // WebSocketで盤面が送られてくるまでは空（0）で埋める
  return new Array(16).fill(0);
}

const board = ref(generateInitialBoard());

const showStartModal = ref(true);
const showResultModal = ref(false);
const currentRoom = ref<Room | null>(null);

// 盤面の変更を監視してWebSocketで送信
watch(board, (newBoard: number[], oldBoard: number[]) => {
  console.log("Board updated:", newBoard);

  // ゲームが開始されていて、実際に盤面が変更された場合のみ送信
  if (
    gameStarted.value &&
    oldBoard &&
    JSON.stringify(newBoard) !== JSON.stringify(oldBoard)
  ) {
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

// 現在のユーザーのスコアを playerScores から取得
const gameScore = computed(() => {
  const currentUsername = webSocketStore.currentUsername;
  console.log("Current username:", currentUsername);
  const playerData = playerScores.value.get(currentUsername);
  return playerData ? playerData.score : 0;
});

// 時間をMM:SS形式でフォーマット
function formatTime(seconds: number): string {
  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = seconds % 60;
  return `${minutes.toString().padStart(2, "0")}:${remainingSeconds
    .toString()
    .padStart(2, "0")}`;
}

// 時間の状態に応じてCSSクラスを返す
function getTimeClass(seconds: number): string {
  if (isTimePlaceholder(seconds)) return "";

  if (seconds <= 10 && seconds > 0) {
    return "urgent-pulse";
  } else if (seconds <= 30 && seconds > 0) {
    return "warning";
  }
  return "";
}

// 時間がプレースホルダー状態かどうかを判定
function isTimePlaceholder(seconds: number): boolean {
  // ゲーム開始前（-1）またはundefined/nullの場合のみプレースホルダー
  // ゲーム終了時の0秒は有効な時間として扱う
  return seconds === undefined || seconds === null || seconds < 0;
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
  transition: color 0.3s ease;
}

/* 30秒以下で赤色に */
.warning {
  color: #ff4757 !important;
}

/* 10秒以下で強調アニメーション */
.urgent-pulse {
  color: #ff3838 !important;
  animation: time-urgent 0.8s ease-in-out infinite !important;
  font-weight: bold !important;
}

@keyframes time-urgent {
  0% {
    transform: scale(1);
    opacity: 1;
    text-shadow: 0 0 0px rgba(255, 56, 56, 0);
  }
  25% {
    transform: scale(1.3);
    opacity: 0.9;
    text-shadow: 0 0 15px rgba(255, 56, 56, 0.8);
  }
  50% {
    transform: scale(1.2);
    opacity: 0.7;
    text-shadow: 0 0 25px rgba(255, 56, 56, 1);
  }
  75% {
    transform: scale(1.25);
    opacity: 0.9;
    text-shadow: 0 0 10px rgba(255, 56, 56, 0.6);
  }
  100% {
    transform: scale(1);
    opacity: 1;
    text-shadow: 0 0 0px rgba(255, 56, 56, 0);
  }
}

.board {
  padding: 14px;
}
</style>
