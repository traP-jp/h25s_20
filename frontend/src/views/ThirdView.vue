<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from "vue";
import { useWebSocket } from "@/lib/websocket";
import { getWsUrl } from "@/config/app";

// 接続設定
const wsUrl = getWsUrl("username=debug_user");

// WebSocket接続の初期化
const {
  isConnected,
  isConnecting,
  connectionError,
  messages,
  connect,
  manualConnect,
  disconnect,
  clearMessages,
  destroy,
  getReconnectAttempts,
  getMaxReconnectAttempts,
} = useWebSocket(wsUrl);

// 再接続情報の取得用
const reconnectAttempts = ref(0);
const maxReconnectAttempts = ref(5);

// 再接続情報を定期的に更新
const updateReconnectInfo = () => {
  reconnectAttempts.value = getReconnectAttempts();
  maxReconnectAttempts.value = getMaxReconnectAttempts();
};

// 定期的に再接続情報を更新
setInterval(updateReconnectInfo, 1000);

// ライフサイクル
onMounted(() => {
  connect();
});

onBeforeUnmount(() => {
  destroy();
});
</script>

<template>
  <div class="websocket-test">
    <h1>WebSocket接続テスト</h1>

    <!-- 接続状態表示 -->
    <div class="connection-status">
      <div v-if="isConnecting" class="status connecting">🔄 接続中...</div>
      <div v-else-if="isConnected" class="status connected">✅ 接続済み</div>
      <div v-else class="status disconnected">❌ 未接続</div>

      <div v-if="connectionError" class="error">
        {{ connectionError }}
      </div>

      <div v-if="reconnectAttempts > 0" class="reconnect-info">
        再接続試行: {{ reconnectAttempts }}/{{ maxReconnectAttempts }}
      </div>
    </div>

    <!-- 操作ボタン -->
    <div class="controls">
      <button @click="manualConnect" :disabled="isConnecting" class="btn btn-connect">
        {{ isConnecting ? "接続中..." : "接続" }}
      </button>

      <button @click="disconnect" :disabled="!isConnected && !isConnecting" class="btn btn-disconnect">切断</button>

      <button @click="clearMessages" class="btn btn-clear">ログクリア</button>
    </div>

    <!-- メッセージ表示 -->
    <div class="messages">
      <h3>受信メッセージ ({{ messages.length }}件)</h3>
      <div class="message-list">
        <div v-for="(msg, idx) in messages.slice().reverse()" :key="idx" class="message">
          {{ msg }}
        </div>
        <div v-if="messages.length === 0" class="no-messages">メッセージはありません</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.websocket-test {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
}

.connection-status {
  margin-bottom: 20px;
  padding: 15px;
  border-radius: 8px;
  background-color: #f5f5f5;
}

.status {
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 8px;
}

.status.connected {
  color: #28a745;
}

.status.connecting {
  color: #ffc107;
}

.status.disconnected {
  color: #dc3545;
}

.error {
  color: #dc3545;
  font-size: 14px;
  margin-top: 8px;
}

.reconnect-info {
  color: #6c757d;
  font-size: 14px;
  margin-top: 8px;
}

.controls {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: background-color 0.2s;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-connect {
  background-color: #007bff;
  color: white;
}

.btn-connect:hover:not(:disabled) {
  background-color: #0056b3;
}

.btn-disconnect {
  background-color: #dc3545;
  color: white;
}

.btn-disconnect:hover:not(:disabled) {
  background-color: #c82333;
}

.btn-clear {
  background-color: #6c757d;
  color: white;
}

.btn-clear:hover {
  background-color: #545b62;
}

.messages {
  border: 1px solid #dee2e6;
  border-radius: 8px;
  overflow: hidden;
}

.messages h3 {
  margin: 0;
  padding: 15px;
  background-color: #f8f9fa;
  border-bottom: 1px solid #dee2e6;
  font-size: 16px;
}

.message-list {
  max-height: 400px;
  overflow-y: auto;
  padding: 10px;
}

.message {
  padding: 8px 0;
  border-bottom: 1px solid #eee;
  font-family: "Monaco", "Menlo", "Ubuntu Mono", monospace;
  font-size: 13px;
  line-height: 1.4;
  word-break: break-all;
}

.message:last-child {
  border-bottom: none;
}

.no-messages {
  text-align: center;
  color: #6c757d;
  font-style: italic;
  padding: 20px;
}

/* レスポンシブ対応 */
@media (max-width: 600px) {
  .websocket-test {
    padding: 10px;
  }

  .controls {
    flex-direction: column;
  }

  .btn {
    width: 100%;
  }
}
</style>
