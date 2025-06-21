<template>
  <div :class="$style.container">
    <div :class="$style.header">Select a Room</div>
    <div :class="$style.rooms">
      <RoomButton
        v-for="room in roomData"
        :key="room.id"
        :room="room"
        @click="handleRoomClick(room)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { roomData } from "@/lib/sample-data";
import RoomButton from "@/components/RoomButton.vue";

function handleRoomClick(room) {
  console.log("Room clicked:", room);
}
</script>

<style module>
.container {
  width: 500px;
  height: 100vh;
  margin: 0 auto;
  padding: 20px;
  border: 1px solid red;
  display: flex;
  flex-direction: column;
  justify-content: center;
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
</style>
