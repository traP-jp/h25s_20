<template>
  <div :class="$style.container">
    <button
      :class="$style.mainBtn"
      @click="onClickMain"
      :disabled="isBtnDisabled"
      v-show="!isLeader || !allPlayersReady"
    >
      {{ BtnMsg }}
    </button>
    <!-- リーダーかつ全プレイヤーが準備完了の場合にゲーム開始ボタンを表示 -->
    <button v-show="isLeader && allPlayersReady" :class="$style.startBtn" @click="onClickStart">ゲーム開始</button>
    <button v-show="!isBtnDisabled" :class="$style.quitBtn" @click="onClickQuit">部屋から抜ける</button>
    <button v-show="isBtnDisabled" :class="$style.cancelBtn" @click="onClickCancel">キャンセル</button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, defineModel, onMounted, onUnmounted } from "vue";
import { useRouter } from "vue-router";
import { useCurrentRoomStore, useRoomPlayersStore, useWebSocketStore } from "@/store";
import { apiClient } from "@/api";
import { WS_EVENTS, type WebSocketEvent } from "@/lib/websocket";

const showResultModal = defineModel<boolean>("showResultModal");
const currentRoomStore = useCurrentRoomStore();
const roomPlayersStore = useRoomPlayersStore();
const webSocketStore = useWebSocketStore();
const router = useRouter();

const BtnMsg = ref("準備OK!");
const isBtnDisabled = ref(false);

// 現在のユーザーがリーダー（プレイヤーリストの先頭）かどうかを判定
const isLeader = computed(() => {
  const players = roomPlayersStore.players;
  const currentUsername = webSocketStore.currentUsername;
  return players.length > 0 && players[0].name === currentUsername;
});

// 全プレイヤーが準備完了かどうかを判定
const allPlayersReady = computed(() => {
  const players = roomPlayersStore.players;
  return players.length > 0 && players.every((player) => player.isReady);
});

// WebSocketイベントハンドラー
const handleWebSocketEvent = (event: WebSocketEvent) => {
  if (event.event === WS_EVENTS.PLAYER_ALL_READY) {
    // 全プレイヤーのisReadyをtrueに設定
    roomPlayersStore.players.forEach((player) => {
      player.isReady = true;
    });
    console.log("All players are ready - updated store");
  }
};

// WebSocketイベントリスナーの設定
onMounted(() => {
  // 既存のWebSocket接続にイベントハンドラーを追加
  if (webSocketStore.wsManager) {
    webSocketStore.wsManager.addMessageHandler(handleWebSocketEvent);
    console.log("Added PLAYER_ALL_READY event handler to existing WebSocket connection");
  }
});

onUnmounted(() => {
  // コンポーネントが破棄される際の処理
  // 他のコンポーネントが使用している可能性があるため、WebSocket自体は切断しない
});

const onClickMain = async () => {
  const room = currentRoomStore.getCurrentRoom();
  if (!room) {
    console.error("No current room found");
    return;
  }

  try {
    // READYアクションをバックエンドに送信
    const response = await apiClient.performRoomAction(room.roomId, { action: "READY" });

    if (response.success) {
      console.log("Successfully sent READY action");
      // メッセージとボタン状態を変更
      BtnMsg.value = "他のプレイヤーを待っています...";
      isBtnDisabled.value = true;
    } else {
      console.error("Failed to send READY action:", response.data);
      alert("準備完了の送信に失敗しました");
    }
  } catch (error) {
    console.error("Error sending READY action:", error);
    alert("準備完了の送信中にエラーが発生しました");
  }
};

const onClickQuit = async () => {
  const room = currentRoomStore.getCurrentRoom();
  if (!room) {
    console.error("No current room found");
    return;
  }

  try {
    // ABORTアクションをバックエンドに送信（部屋から抜ける）
    const response = await apiClient.performRoomAction(room.roomId, { action: "ABORT" });

    if (response.success) {
      console.log("Successfully left the room");
      showResultModal.value = false;
      // ルーム情報をクリアしてルーム選択画面に戻る
      currentRoomStore.clearCurrentRoom();
      router.push("/rooms");
    } else {
      console.error("Failed to leave room:", response.data);
      alert("部屋から抜ける処理に失敗しました");
    }
  } catch (error) {
    console.error("Error leaving room:", error);
    alert("部屋から抜ける処理中にエラーが発生しました");
  }
};

const onClickCancel = async () => {
  const room = currentRoomStore.getCurrentRoom();
  if (!room) {
    console.error("No current room found");
    return;
  }

  try {
    // CANCELアクションをバックエンドに送信
    const response = await apiClient.performRoomAction(room.roomId, { action: "CANCEL" });

    if (response.success) {
      console.log("Successfully canceled READY state");
      // メッセージとボタン状態を元に戻す
      BtnMsg.value = "準備OK!";
      isBtnDisabled.value = false;
    } else {
      console.error("Failed to cancel READY state:", response.data);
      alert("キャンセルの送信に失敗しました");
    }
  } catch (error) {
    console.error("Error sending CANCEL action:", error);
    alert("キャンセルの送信中にエラーが発生しました");
  }
};

const onClickStart = async () => {
  const room = currentRoomStore.getCurrentRoom();
  if (!room) {
    console.error("No current room found");
    return;
  }

  try {
    // STARTアクションをバックエンドに送信
    const response = await apiClient.performRoomAction(room.roomId, { action: "START" });

    if (response.success) {
      console.log("Successfully started the game");
      // ゲーム開始後の処理は WebSocket イベントで処理される想定
      // モーダルの制御はWebSocketイベントで行われるため、ここでは何もしない
    } else {
      console.error("Failed to start game:", response.data);
      alert("ゲーム開始の送信に失敗しました");
    }
  } catch (error) {
    console.error("Error starting game:", error);
    alert("ゲーム開始の送信中にエラーが発生しました");
  }
};
</script>

<style module>
.container {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
  background-color: white;
  padding: 20px;
  z-index: 1000;
}

.mainBtn {
  transform: scale(1.5);
  margin: 20px 0;
}

.startBtn {
  border: 1px solid red;
  /* background-color: #4caf50; */
  color: white;
  border: none;
  padding: 15px 30px;
  font-size: 18px;
  border-radius: 8px;
  cursor: pointer;
  margin: 10px 0;
  transform: scale(1.2);
}

.startBtn:hover {
  opacity: 0.8;
}

.quitBtn {
  align-self: flex-start;
}

.cancelBtn {
  align-self: flex-start;
}
</style>
