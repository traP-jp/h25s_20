<template>
  <div :class="$style.container">
    <div :class="$style.header">Select a Room</div>
    <div :class="$style.rooms" v-if="roomData.length > 0">
      <RoomButton v-for="room in roomData" :key="room.roomId" :room="room" @click="handleRoomClick(room)" />
    </div>
    <div :class="$style.empty" v-else>No rooms available</div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { apiClient } from "@/api";
import RoomButton from "@/components/RoomButton.vue";
import type { Room } from "@/lib/types";

const roomData = ref<Room[]>([]);

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

function handleRoomClick(room: Room) {
  console.log("Room clicked:", room);
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
}

.header {
  font-size: 24px;
  font-weight: bold;
  text-align: center;
  margin-bottom: 16px;
}

.rooms {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(100px, 1fr));
  margin: 10px;
}

.empty {
  text-align: center;
  color: #666;
  margin: 20px;
  font-size: 16px;
}
</style>
