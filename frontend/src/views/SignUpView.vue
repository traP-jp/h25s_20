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
    <div v-if="error" :class="$style.error">
      {{ error }}
    </div>
    <button
      :disabled="!isValid || isLoading"
      :class="[$style.button, { [$style.loading]: isLoading }]"
      @click="onClick"
    >
      {{ isLoading ? "作成中..." : "ゲームをはじめる" }}
    </button>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed } from "vue";
import { useRouter } from "vue-router";
import UserIcon from "@/components/UserIcon.vue";
import { apiClient } from "@/api";
import { useWebSocketStore } from "@/store";

const username = ref("");
const router = useRouter();
const isComposing = ref(false);
const isLoading = ref(false);
const error = ref<string | null>(null);

// WebSocketストアを取得
const webSocketStore = useWebSocketStore();

const onClick = async () => {
  if (!isValid.value || isLoading.value) return;

  isLoading.value = true;
  error.value = null;

  try {
    // ユーザー作成API呼び出し
    const response = await apiClient.createUser({
      username: username.value.trim(),
      password: "", // パスワードが不要な場合は空文字、必要な場合は適切な値を設定
    });

    if (response.success) {
      // ユーザー作成成功時の処理
      // JWTトークンがレスポンスに含まれる場合は保存
      if (response.data?.token) {
        apiClient.setAuthToken(response.data.token);
        // セッションストレージに保存（タブごとに独立）
        sessionStorage.setItem("authToken", response.data.token);
        sessionStorage.setItem("username", username.value.trim());
      }

      // WebSocket接続を初期化
      console.log("WebSocket接続を初期化します:", username.value.trim());
      webSocketStore.initializeWebSocket(username.value.trim());

      router.push("/rooms");
    } else {
      error.value = response.data?.message || "ユーザー作成に失敗しました";
    }
  } catch (err) {
    error.value = "ネットワークエラーが発生しました";
    console.error("User creation error:", err);
  } finally {
    isLoading.value = false;
  }
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
  return username.value.trim().length >= 1 && username.value.trim().length <= 32;
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

.error {
  color: #ff6b6b;
  font-size: 14px;
  margin-bottom: 10px;
  text-align: center;
}

.button {
  transition: opacity 0.2s ease;
}

.button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.loading {
  opacity: 0.8;
}
</style>
