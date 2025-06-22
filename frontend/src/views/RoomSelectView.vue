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
import { useRouter } from "vue-router";
import { roomData } from "@/lib/sample-data";
import RoomButton from "@/components/RoomButton.vue";
import type { Room } from "@/lib/types";

const router = useRouter();

function handleRoomClick(room: Room) {
  console.log("Room clicked:", room);
}

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
  text-align: center; text-align: center;
}

.button:hover {
  background-color: #218838; background-color: #218838;
}
</style>
