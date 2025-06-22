<template>
  <div :class="$style.container">
    <button :class="$style.mainBtn" @click="onClickMain" :disabled="isBtnDisabled">
      {{ BtnMsg }}
    </button>
    <button v-show="!isBtnDisabled" :class="$style.quitBtn" @click="onClickQuit">部屋から抜ける</button>
    <button v-show="isBtnDisabled" :class="$style.cancelBtn" @click="onClickCancel">キャンセル</button>
  </div>
</template>

<script setup lang="ts">
import { ref, defineModel } from "vue";
import { useRouter } from "vue-router";
import { useCurrentRoomStore } from "@/store";
import { apiClient } from "@/api";

const showResultModal = defineModel<boolean>("showResultModal");
const currentRoomStore = useCurrentRoomStore();
const router = useRouter();

const BtnMsg = ref("準備OK!");
const isBtnDisabled = ref(false);

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

.quitBtn {
  align-self: flex-start;
}

.cancelBtn {
  align-self: flex-start;
}
</style>
