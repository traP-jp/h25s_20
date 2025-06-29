<template>
  <div :class="$style.container">
    <img :class="$style.logo" src="/logo.svg" alt="Logo" />
    <div :class="$style.header">部屋を選んで入室</div>
    
    <div v-if="isLoading" :class="$style.loading">
      部屋一覧を読み込み中...
    </div>
    
    <div v-else-if="roomData.length === 0 && retryCount < maxRetries" :class="$style.empty">
      <p>部屋が見つかりません</p>
      <button :class="$style.retryButton" @click="fetchRooms" :disabled="isLoading">
        再試行 ({{ retryCount }}/{{ maxRetries }})
      </button>
    </div>
    
    <div v-else-if="roomData.length === 0" :class="$style.empty">
      <p>部屋の取得に失敗しました</p>
      <button :class="$style.retryButton" @click="() => { retryCount = 0; fetchRooms(); }" :disabled="isLoading">
        最初からやり直す
      </button>
      <button :class="$style.backButton" @click="() => router.push('/')">
        ユーザー選択に戻る
      </button>
    </div>
    
    <div v-else :class="$style.rooms">
      <RoomButton v-for="room in roomData" :key="room.roomId" :room="room" @click="handleRoomClick(room)" />
    </div>

    <button :class="$style.button" @click="onClick">プレイ方法を確認</button>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import { apiClient } from "@/api";
import RoomButton from "@/components/RoomButton.vue";
import type { Room } from "@/lib/types";
import { useCurrentRoomStore } from "@/store";

const roomData = ref<Room[]>([]);
const router = useRouter();
const currentRoomStore = useCurrentRoomStore();
const isLoading = ref(false);
const retryCount = ref(0);
const maxRetries = 3;

async function fetchRooms() {
  if (isLoading.value) return;
  
  isLoading.value = true;
  try {
    const response = await apiClient.getRooms();
    if (response.success) {
      // ルームIDの昇順でソート
      roomData.value = response.data.sort((a: Room, b: Room) => a.roomId - b.roomId);
      console.log("Successfully fetched rooms:", response.data);
      retryCount.value = 0; // 成功したらリトライカウントをリセット
    } else {
      console.error("Failed to fetch rooms:", response.data);
      await handleFetchError(response.status);
    }
  } catch (error) {
    console.error("Error fetching rooms:", error);
    await handleFetchError(0);
  } finally {
    isLoading.value = false;
  }
}

async function handleFetchError(status: number) {
  // 認証エラー（401）またはユーザーが見つからない（404）の場合はユーザー選択画面に戻る
  if (status === 401 || status === 404) {
    console.warn("Authentication failed, redirecting to signup");
    alert("認証に失敗しました。ユーザー選択画面に戻ります。");
    router.push("/");
    return;
  }

  // その他のエラーでリトライ可能な場合は再試行
  if (retryCount.value < maxRetries) {
    retryCount.value++;
    console.log(`Retrying to fetch rooms (attempt ${retryCount.value}/${maxRetries})`);
    setTimeout(() => {
      fetchRooms();
    }, 1000 * retryCount.value); // 線形バックオフ的に待機時間を増加
  } else {
    console.error("Max retries reached, giving up");
    alert("部屋一覧の取得に失敗しました。ユーザー選択画面に戻ります。");
    router.push("/");
  }
}

function checkAuthToken() {
  // APIクライアントが認証トークンを持っているかチェック
  const hasToken = sessionStorage.getItem("authToken");
  if (!hasToken) {
    console.warn("No auth token found, redirecting to signup");
    alert("認証情報が見つかりません。ユーザー選択画面に戻ります。");
    router.push("/");
    return false;
  }
  return true;
}

async function handleRoomClick(room: Room) {
  try {
    console.log("Joining room:", room);
    const response = await apiClient.performRoomAction(room.roomId, { action: "JOIN" });

    if (response.success) {
      console.log("Successfully joined room:", room.roomId);

      // ルーム情報をストアに保存
      currentRoomStore.setCurrentRoom(room);

      // ルーム参加後、最新のルーム情報（ユーザー情報含む）を取得
      let finalRoom = room; // デフォルトは元のルーム情報
      try {
        const roomsResponse = await apiClient.getRooms();
        if (roomsResponse.success) {
          // 参加したルームの最新情報を取得
          const updatedRoom = roomsResponse.data.find((r: Room) => r.roomId === room.roomId);
          if (updatedRoom) {
            // 最新のルーム情報でストアを更新
            currentRoomStore.setCurrentRoom(updatedRoom);
            finalRoom = updatedRoom; // 最新のルーム情報を使用
            console.log("Updated room with latest user info:", updatedRoom);
          }
        } else {
          console.warn("Failed to fetch updated room details:", roomsResponse.data);
        }
      } catch (error) {
        console.warn("Error fetching updated room details:", error);
      }

      // ルームに参加成功後、PlayViewに遷移（最新のルーム情報を使用）
      router.push({
        name: "play",
        params: { roomId: finalRoom.roomId.toString() },
        query: {
          roomName: finalRoom.roomName,
          isOpened: finalRoom.isOpened.toString(),
        },
      });
    } else {
      console.error("Failed to join room:", response.data);
      alert("ルームへの参加に失敗しました");
    }
  } catch (error) {
    console.error("Error joining room:", error);
    alert("ルームへの参加中にエラーが発生しました");
  }
}

onMounted(() => {
  // 認証トークンをチェックしてから部屋一覧を取得
  if (checkAuthToken()) {
    fetchRooms();
  }
});

const onClick = () => {
  router.push("/help");
};
</script>

<style module>
.container {
  width: 360px;
  height: 100%;
  margin: 40px auto;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.logo {
  width: 300px;
  margin-bottom: 26px;
}

.header {
  font-size: 24px;
  font-weight: bold;
  text-align: center;
  margin-bottom: 20px;
}

.rooms {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(100px, 1fr));
  gap: 12px;
  margin: 10px;
  margin-bottom: 30px;
  width: 100%;
}

.button {
  width: 220px;
  padding: 12px 24px;
  font-size: 16px;
  font-weight: 500;
  background-color: #28a745;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.2s ease;
  text-align: center;
  text-align: center;
}

.button:hover {
  background-color: #218838 !important;
}

.empty {
  text-align: center;
  color: #666;
  margin: 20px;
  font-size: 16px;
}

.loading {
  text-align: center;
  color: #666;
  margin: 20px;
  font-size: 16px;
  font-style: italic;
}

.retryButton {
  padding: 8px 16px;
  margin-top: 10px;
  margin-right: 8px;
  font-size: 14px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.retryButton:hover:not(:disabled) {
  background-color: #0056b3;
}

.retryButton:disabled {
  background-color: #6c757d;
  cursor: not-allowed;
}

.backButton {
  padding: 8px 16px;
  margin-top: 10px;
  margin-left: 8px;
  font-size: 14px;
  background-color: #6c757d;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.backButton:hover {
  background-color: #545b62;
}
</style>
