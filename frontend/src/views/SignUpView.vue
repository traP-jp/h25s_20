<template>
  <div class="container">

    <div :class="$style.logo">
      <img src="/logo.svg" alt="Logo">
    </div>

    <h2>名前を入れてはじめよう！</h2>
    <div :class="$style.user_info">
      <UserIcon :id="username" :size="50" />
      <input v-model="username" required type="text" placeholder="名前を入力" />
    </div>
    <button :class="$style.button" @click="onClick">ゲームをはじめる</button>

  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import UserIcon from "@/components/UserIcon.vue";

const username = ref('');
const router = useRouter();

const onClick = () => {
  router.push("/");
}
</script>

<style module>

.container {
  height: 100vh;
}

.logo img {
  width: 300px;
}

.user_info {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: row;
  gap: 20px;
}

.icon img {
  width: 60px;
  height: 60px;
  clip-path: circle(50%);
}

</style>
