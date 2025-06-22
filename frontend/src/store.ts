import { defineStore } from "pinia";
import { ref } from "vue";
import { useWebSocket, type WebSocketEvent } from "@/lib/websocket";

export const useNotificationStore = defineStore("notificationStore", () => {
  const notifications = ref<string[]>([]);

  const addNotification = (message: string) => {
    notifications.value.push(message);
  };

  return {
    notifications,
    addNotification,
  };
});

// WebSocket接続のグローバル管理
export const useWebSocketStore = defineStore("webSocketStore", () => {
  const wsManager = ref<ReturnType<typeof useWebSocket> | null>(null);
  const isConnected = ref(false);
  const currentUsername = ref<string>("");

  // WebSocket接続を初期化
  const initializeWebSocket = (username: string, onMessage?: (event: WebSocketEvent) => void) => {
    if (wsManager.value) {
      // 既存の接続がある場合は切断
      wsManager.value.destroy();
    }

    currentUsername.value = username;
    const wsUrl = `wss://10ten.trap.show/api/ws?username=${encodeURIComponent(username)}`;

    wsManager.value = useWebSocket(wsUrl, (event: WebSocketEvent) => {
      console.log("Global WebSocket received:", event);
      if (onMessage) {
        onMessage(event);
      }
    });

    // 接続開始
    if (wsManager.value) {
      wsManager.value.connect();
    }

    return wsManager.value;
  };

  // WebSocket接続を切断
  const disconnectWebSocket = () => {
    if (wsManager.value) {
      wsManager.value.destroy();
      wsManager.value = null;
      isConnected.value = false;
    }
  };

  // メッセージ送信
  const sendMessage = (event: any) => {
    if (wsManager.value) {
      wsManager.value.send(event);
    } else {
      console.warn("WebSocket is not connected");
    }
  };

  // 現在の接続マネージャーを取得
  const getWebSocketManager = () => {
    return wsManager.value;
  };

  // 現在の接続状態を取得
  const getConnectionStatus = () => {
    if (!wsManager.value) return false;
    // 型エラー回避のため any でキャスト
    return (wsManager.value as any).isConnected?.value || false;
  };

  return {
    wsManager,
    isConnected,
    currentUsername,
    initializeWebSocket,
    disconnectWebSocket,
    sendMessage,
    getWebSocketManager,
    getConnectionStatus,
  };
});
