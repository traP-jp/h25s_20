<template>
  <div :class="$style.container">
    <div :class="$style.header">Select a Room</div>
    <div :class="$style.rooms">
      <RoomButton v-for="room in roomData" :key="room.id" :room="room" @click="handleRoomClick(room)" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { roomData } from "@/lib/sample-data";
import RoomButton from "@/components/RoomButton.vue";
import type { Room } from "@/lib/types";

function handleRoomClick(room: Room) {
  console.log("Room clicked:", room);
}
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
</style>
