<template>
  <div :class="$style.container">
    <img :class="$style.logo" src="/logo.svg" alt="Logo" />
    <div :class="$style.userInfo">
      <UserIcon :id="username" :size="50" />
      <div :class="$style.input">
        <div v-if="false">{{ username }}</div>
        <input
          v-else
          v-model="username"
          required
          type="text"
          placeholder="名前を入力"
          @keydown.enter="onEnter"
          @compositionstart="onCompositionStart"
          @compositionend="onCompositionEnd"
        />
      </div>
    </div>
    <button :disabled="!isValid" :class="$style.button" @click="onClick">
      ゲームをはじめる
    </button>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed } from "vue";
import { useRouter } from "vue-router";
import UserIcon from "@/components/UserIcon.vue";

const username = ref("");
const router = useRouter();
const isComposing = ref(false);

const onClick = () => {
  router.push("/rooms");
};

const onEnter = () => {
  if (!isComposing.value) {
    onClick();
  }
};

const onCompositionStart = () => {
  isComposing.value = true;
};

const onCompositionEnd = () => {
  isComposing.value = false;
};

const isValid = computed(() => {
  return (
    username.value.trim().length >= 1 && username.value.trim().length <= 32
  );
});
</script>

// ...existing code...

<style module>
.container {
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.guide {
  text-align: center;
}

.logo {
  width: 300px;
}

.userInfo {
  margin: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: row;
  gap: 20px;
}

.input {
  width: 200px;
  font-weight: 500;
  border-bottom: 1px solid #bbb;
}

.icon img {
  width: 60px;
  height: 60px;
  clip-path: circle(50%);
}
</style>
