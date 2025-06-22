<template>
  <div :class="$style.container">
    <img :class="$style.logo" src="/logo.svg" alt="Logo" />
    <div :class="$style.header">部屋を選んで入室</div>
    <div :class="$style.rooms">
      <RoomButton v-for="room in roomData" :key="room.id" :room="room" @click="handleRoomClick(room)" />
    </div>

    <button :class="$style.button" @click="onClick">
      プレイ方法を確認
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import { apiClient } from "@/api";
import RoomButton from "@/components/RoomButton.vue";
import type { Room } from "@/lib/types";

const roomData = ref<Room[]>([]);
const router = useRouter();

async function fetchRooms() {
  try {
    const response = await apiClient.getRooms();
    if (response.success) {
      roomData.value = response.data;
      console.log(response.data);
    } else {
      console.error("Failed to fetch rooms:", response.data);
    }
  } catch (error) {
    console.error("Error fetching rooms:", error);
  }
}

async function handleRoomClick(room: Room) {
  try {
    console.log("Joining room:", room);
    const response = await apiClient.performRoomAction(room.roomId, { action: "JOIN" });

    if (response.success) {
      console.log("Successfully joined room:", room.roomId);
      // ルームに参加成功後、PlayViewに遷移（ルーム情報をクエリパラメータとして渡す）
      router.push({
        name: "play",
        params: { roomId: room.roomId.toString() },
        query: {
          roomName: room.roomName,
          isOpened: room.isOpened.toString(),
        },
        state: { room },
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
  fetchRooms();
});
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
  text-align: center; text-align: center;
}

.button:hover {
  background-color: #218838; background-color: #218838;
}

.empty {
  text-align: center;
  color: #666;
  margin: 20px;
  font-size: 16px;
}
</style>
