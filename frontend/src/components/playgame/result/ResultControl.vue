<template>
  <div :class="$style.container">
    <button :class="$style.mainBtn" @click="onClickMain">Got it!</button>
    <button :class="$style.quitBtn" @click="onClickQuit">部屋から抜ける</button>
  </div>
</template>
<script setup lang="ts">
import { defineModel } from "vue";
import { useRouter } from "vue-router";
import { useCurrentRoomStore } from "@/store";
import { apiClient } from "@/api";

const showStartModal = defineModel<boolean>("showStartModal");
const showResultModal = defineModel<boolean>("showResultModal");
const currentRoomStore = useCurrentRoomStore();
const router = useRouter();

const onClickMain = async () => {
  const room = currentRoomStore.getCurrentRoom();
  if (!room) {
    console.error("No current room found");
    return;
  }

  try {
    // CLOSE_RESULTアクションをバックエンドに送信（結果画面を閉じて次のゲームの準備）
    const response = await apiClient.performRoomAction(room.roomId, { action: "CLOSE_RESULT" });

    if (response.success) {
      console.log("Successfully closed result and ready for next game");
      // 結果モーダルを閉じて開始モーダルを表示
      showResultModal.value = false;
      showStartModal.value = true;
    } else {
      console.error("Failed to close result:", response.data);
      alert("結果画面の終了に失敗しました");
    }
  } catch (error) {
    console.error("Error closing result:", error);
    alert("結果画面の終了中にエラーが発生しました");
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
      console.log("Successfully left the room from result screen");
      showResultModal.value = false;
      showStartModal.value = false;
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
  padding: 12px 24px;
  font-size: 16px;
  font-weight: bold;
  border-radius: 8px;
  margin: 20px 0;
  transform: scale(1.5);
  min-width: 150px;
  box-sizing: border-box;
  text-align: center;
}

.quitBtn {
  align-self: flex-start;
}
</style>
