<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from "vue";

// WebSocketイベントの型定義
interface WebSocketEvent {
  event: string;
  content: EventContent;
}

interface BaseEventContent {
  user_id?: number;
  user_name?: string;
  room_id?: number;
  message?: string;
  timestamp?: number;
}

interface ConnectionEventContent extends BaseEventContent {
  client_id: string;
}

interface PlayerEventContent extends BaseEventContent {}

interface BoardData {
  content: number[];
  version: number;
  size: number;
}

interface BoardUpdateEventContent extends BaseEventContent {
  board: BoardData;
  gain_score: number;
}

interface CountdownEventContent extends BaseEventContent {
  count?: number;
  countdown?: number;
}

type EventContent = 
  | ConnectionEventContent 
  | PlayerEventContent 
  | BoardUpdateEventContent 
  | CountdownEventContent 
  | BaseEventContent;

// WebSocketイベント名の定数
const WS_EVENTS = {
  CONNECTION: 'connection',
  PLAYER_JOINED: 'player_joined',
  PLAYER_READY: 'player_ready',
  PLAYER_CANCELED: 'player_canceled',
  PLAYER_LEFT: 'player_left',
  GAME_STARTED: 'game_started',
  GAME_START: 'game_start',
  COUNTDOWN_START: 'countdown_start',
  COUNTDOWN: 'countdown',
  BOARD_UPDATED: 'board_updated',
  RESULT_CLOSED: 'result_closed',
  GAME_ENDED: 'game_ended'
} as const;

// リアクティブな状態
const isConnected = ref(false);
const isConnecting = ref(false);
const ws = ref<WebSocket | null>(null);
const messages = ref<string[]>([]);
const connectionError = ref<string | null>(null);
const reconnectAttempts = ref(0);
const maxReconnectAttempts = 5;

// 接続先の選択
const wsEndpoint = ref<'local' | 'docker'>('local');
const customUsername = ref('debug_user');

// WebSocket接続関数
const connectWebSocket = () => {
  if (isConnecting.value || (ws.value && ws.value.readyState === WebSocket.OPEN)) {
    return;
  }

  isConnecting.value = true;
  connectionError.value = null;

  // デバッグユーザーでの接続（本番では適切なusernameを使用）
  const wsUrl = "ws://localhost:8080/ws?username=debug_user";
  
  try {
    ws.value = new WebSocket(wsUrl);

    // 接続開始のタイムアウト設定
    const connectTimeout = setTimeout(() => {
      if (ws.value && ws.value.readyState === WebSocket.CONNECTING) {
        ws.value.close();
        connectionError.value = "接続タイムアウト";
        isConnecting.value = false;
      }
    }, 10000); // 10秒でタイムアウト

    ws.value.onopen = () => {
      clearTimeout(connectTimeout);
      console.log("WebSocket接続が確立されました");
      isConnected.value = true;
      isConnecting.value = false;
      reconnectAttempts.value = 0;
      connectionError.value = null;
      
      addMessage("✅ WebSocket接続が確立されました");
    };

    ws.value.onmessage = (event) => {
      try {
        const wsEvent: WebSocketEvent = JSON.parse(event.data);
        handleWebSocketEvent(wsEvent);
      } catch (error) {
        console.error("WebSocketメッセージの解析に失敗:", error);
        addMessage(`❌ メッセージ解析エラー: ${event.data}`);
      }
    };

    ws.value.onerror = (event) => {
      clearTimeout(connectTimeout);
      console.error("WebSocketエラー:", event);
      connectionError.value = "WebSocket接続エラーが発生しました";
      isConnecting.value = false;
      
      addMessage("❌ WebSocket接続エラー");
    };

    ws.value.onclose = (event) => {
      clearTimeout(connectTimeout);
      console.log(`WebSocket接続が閉じられました - Code: ${event.code}, Reason: ${event.reason}`);
      isConnected.value = false;
      isConnecting.value = false;
      
      if (event.code !== 1000) { // 正常終了以外の場合
        connectionError.value = `接続が閉じられました (Code: ${event.code})`;
        addMessage(`⚠️ 接続が閉じられました (Code: ${event.code})`);
        
        // 自動再接続を試行
        attemptReconnect();
      } else {
        addMessage("🔌 WebSocket接続が正常に閉じられました");
      }
    };

  } catch (error) {
    console.error("WebSocket接続の作成に失敗:", error);
    connectionError.value = "WebSocket接続の作成に失敗しました";
    isConnecting.value = false;
    addMessage("❌ WebSocket接続の作成に失敗");
  }
};

// WebSocketイベントハンドラー
const handleWebSocketEvent = (wsEvent: WebSocketEvent) => {
  console.log("受信イベント:", wsEvent);
  
  switch (wsEvent.event) {
    case WS_EVENTS.CONNECTION:
      const connectionContent = wsEvent.content as ConnectionEventContent;
      addMessage(`🔗 接続確立: ClientID: ${connectionContent.client_id}, UserID: ${connectionContent.user_id}`);
      break;
      
    case WS_EVENTS.PLAYER_JOINED:
      const joinedContent = wsEvent.content as PlayerEventContent;
      addMessage(`👤 プレイヤー参加: ${joinedContent.user_name} (ID: ${joinedContent.user_id}) がルーム ${joinedContent.room_id} に参加`);
      break;
      
    case WS_EVENTS.PLAYER_READY:
      const readyContent = wsEvent.content as PlayerEventContent;
      addMessage(`✅ プレイヤー準備完了: ${readyContent.user_name} が準備完了`);
      break;
      
    case WS_EVENTS.PLAYER_CANCELED:
      const canceledContent = wsEvent.content as PlayerEventContent;
      addMessage(`❌ プレイヤー準備キャンセル: ${canceledContent.user_name} が準備をキャンセル`);
      break;
      
    case WS_EVENTS.PLAYER_LEFT:
      const leftContent = wsEvent.content as PlayerEventContent;
      addMessage(`👋 プレイヤー退出: ${leftContent.user_name} がルームから退出`);
      break;
      
    case WS_EVENTS.GAME_STARTED:
      const gameStartedContent = wsEvent.content as BaseEventContent;
      addMessage(`🎮 ゲーム開始: ${gameStartedContent.message}`);
      break;
      
    case WS_EVENTS.COUNTDOWN_START:
      const countdownStartContent = wsEvent.content as CountdownEventContent;
      addMessage(`⏰ カウントダウン開始: ${countdownStartContent.countdown}秒`);
      break;
      
    case WS_EVENTS.COUNTDOWN:
      const countdownContent = wsEvent.content as CountdownEventContent;
      addMessage(`⏱️ カウントダウン: ${countdownContent.count}`);
      break;
      
    case WS_EVENTS.BOARD_UPDATED:
      const boardContent = wsEvent.content as BoardUpdateEventContent;
      addMessage(`📋 ボード更新: ${boardContent.user_name} がスコア ${boardContent.gain_score} 獲得 (Version: ${boardContent.board.version})`);
      break;
      
    case WS_EVENTS.GAME_ENDED:
      const gameEndedContent = wsEvent.content as BaseEventContent;
      addMessage(`🏁 ゲーム終了: ${gameEndedContent.message}`);
      break;
      
    default:
      addMessage(`📨 未知のイベント: ${wsEvent.event} - ${JSON.stringify(wsEvent.content)}`);
  }
};

// メッセージ追加関数
const addMessage = (message: string) => {
  const timestamp = new Date().toLocaleTimeString();
  messages.value.push(`[${timestamp}] ${message}`);
  
  // メッセージ数を制限（パフォーマンス対策）
  if (messages.value.length > 100) {
    messages.value = messages.value.slice(-50);
  }
};

// 再接続試行
const attemptReconnect = () => {
  if (reconnectAttempts.value >= maxReconnectAttempts) {
    addMessage("❌ 最大再接続試行回数に達しました");
    return;
  }
  
  reconnectAttempts.value++;
  const delay = Math.min(1000 * Math.pow(2, reconnectAttempts.value - 1), 30000); // 指数バックオフ（最大30秒）
  
  addMessage(`🔄 ${delay/1000}秒後に再接続を試行します... (${reconnectAttempts.value}/${maxReconnectAttempts})`);
  
  setTimeout(() => {
    connectWebSocket();
  }, delay);
};

// 手動接続関数
const manualConnect = () => {
  if (ws.value) {
    ws.value.close();
  }
  reconnectAttempts.value = 0;
  connectWebSocket();
};

// 切断関数
const disconnect = () => {
  if (ws.value) {
    ws.value.close(1000, "ユーザーによる切断");
  }
  reconnectAttempts.value = maxReconnectAttempts; // 自動再接続を停止
};

// メッセージクリア
const clearMessages = () => {
  messages.value = [];
};

// ライフサイクル
onMounted(() => {
  connectWebSocket();
});

onBeforeUnmount(() => {
  if (ws.value) {
    ws.value.close(1000, "コンポーネントアンマウント");
  }
});
</script>

<template>
  <div class="websocket-test">
    <h1>WebSocket接続テスト</h1>
    
    <!-- 接続状態表示 -->
    <div class="connection-status">
      <div v-if="isConnecting" class="status connecting">
        🔄 接続中...
      </div>
      <div v-else-if="isConnected" class="status connected">
        ✅ 接続済み
      </div>
      <div v-else class="status disconnected">
        ❌ 未接続
      </div>
      
      <div v-if="connectionError" class="error">
        {{ connectionError }}
      </div>
      
      <div v-if="reconnectAttempts > 0" class="reconnect-info">
        再接続試行: {{ reconnectAttempts }}/{{ maxReconnectAttempts }}
      </div>
    </div>
    
    <!-- 操作ボタン -->
    <div class="controls">
      <button 
        @click="manualConnect" 
        :disabled="isConnecting"
        class="btn btn-connect"
      >
        {{ isConnecting ? '接続中...' : '接続' }}
      </button>
      
      <button 
        @click="disconnect" 
        :disabled="!isConnected && !isConnecting"
        class="btn btn-disconnect"
      >
        切断
      </button>
      
      <button 
        @click="clearMessages"
        class="btn btn-clear"
      >
        ログクリア
      </button>
    </div>
    
    <!-- メッセージ表示 -->
    <div class="messages">
      <h3>受信メッセージ ({{ messages.length }}件)</h3>
      <div class="message-list">
        <div 
          v-for="(msg, idx) in messages.slice().reverse()" 
          :key="idx" 
          class="message"
        >
          {{ msg }}
        </div>
        <div v-if="messages.length === 0" class="no-messages">
          メッセージはありません
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.websocket-test {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
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
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
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
