<template>
  <div :class="$style.container">
    <img :class="$style.logo" src="/logo.svg" alt="Logo" />
    <div :class="$style.content">
      <div :class="$style.errorCode">404</div>
      <div :class="$style.message">ページが見つかりません</div>
      <div :class="$style.description">
        お探しのページは存在しないか、削除された可能性があります。
      </div>
    </div>
    <div :class="$style.actions">
      <button :class="$style.button" @click="goHome">
        ホームに戻る
      </button>
      <button :class="$style.button" @click="goBack">
        前のページに戻る
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";

const router = useRouter();

const goHome = () => {
  router.push("/");
};

const goBack = () => {
  router.back();
};
</script>

<style module>
.container {
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 20px;
}

.logo {
  width: 200px;
  margin-bottom: 40px;
}

.content {
  margin-bottom: 40px;
}

.errorCode {
  font-size: 72px;
  font-weight: bold;
  color: #006aeb;
  margin-bottom: 16px;
  line-height: 1;
}

.message {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 16px;
  color: #1a1a1a;
}

.description {
  font-size: 16px;
  color: #666;
  line-height: 1.5;
  margin-bottom: 20px;
}

.actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
  width: 200px;
}

.button {
  width: 100%;
  padding: 12px 24px;
  font-size: 16px;
}
</style> 