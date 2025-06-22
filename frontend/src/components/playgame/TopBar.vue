<template>
  <div :class="$style.container">
    <div :class="$style.content">
      <UserIcon :class="$style.myIcon" :id="`me`" :size="60" />
      <div :class="$style.myName">{{ "me" }}</div>
      <div :class="$style.right">
        <div :class="$style.roomName">{{ room.name }}</div>
        <img src="/logo.svg" alt="Logo" :class="$style.logo" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineModel } from "vue";
import type { Room } from "@/lib/types.ts";
import UserIcon from "@/components/UserIcon.vue";
const room = defineModel<Room>("room");
</script>

<style module>
.container {
  background-color: #ff88bb;
  position: fixed;
  top: 0;
  left: 0;
  z-index: 10;
  width: 100vw;
  height: 50px;
}

.content {
  position: relative;
  width: min(360px, 100%);
  height: 100%;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.myIcon {
  z-index: 1;
  border: 3px solid #ff88bb;
  position: absolute;
  left: 20px;
  top: 14px;
}

.myName {
  margin-top: 10px;
  margin-left: 95px;
  font-weight: bold;
  font-size: 20px;
}

.right {
  height: 100%;
  display: flex;
  align-items: center;
}

.roomName {
  background-color: #008800;
  color: white;
  padding: 0 10px;
  border-radius: 1em;
  font-weight: bold;
}

.logo {
  height: 24px;
  margin: 0 10px;
}
</style>
