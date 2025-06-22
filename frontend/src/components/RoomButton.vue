<template>
  <button :disabled="!room.isOpened" @click="emit('click')" :class="$style.button">
    <div :class="$style.name">{{ room.roomName }}</div>
    <div :class="$style.icons">
      <UserIcon
        :class="$style.icon"
        v-for="player in room.users.slice(0, 3)"
        :key="player.id"
        :id="player.id"
        :size="20"
      />
      <div :class="$style.surplus" v-if="room.users.length > 3">+{{ room.users.length - 3 }}</div>
    </div>
  </button>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from "vue";
import UserIcon from "@/components/UserIcon.vue";

import type { Room } from "@/lib/types.ts";

defineProps<{
  room: Room;
}>();

const emit = defineEmits<{
  click: [];
}>();
</script>

<style module>
.name {
  font-size: 24px;
}

.icons {
  margin: 4px 0;
  display: flex;
  align-items: center;
}

.icon {
  outline: white solid 2px;
}

.surplus {
  margin-left: 6px;
}

.button {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: left;
  margin: 4px;
}
</style>
