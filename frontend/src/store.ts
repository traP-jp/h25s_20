import { defineStore } from "pinia";
import { ref } from "vue";
import { useWebSocket, type WebSocketEvent } from "@/lib/websocket";
import { type ResultPlayer, type StartPlayer, type Room } from "@/lib/types";

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

// ゲーム開始前のプレイヤー管理ストア
export const useRoomPlayersStore = defineStore("roomPlayersStore", () => {
  const players = ref<StartPlayer[]>([]);

  const updatePlayers = (roomPlayers: Array<{ user_name: string; is_ready: boolean }>) => {
    players.value = roomPlayers.map((player) => ({
      name: player.user_name,
      gold: 0, // デフォルト値（実績データが無いため）
      silver: 0,
      bronze: 0,
      isReady: player.is_ready,
    }));
  };

  const setPlayerReady = (username: string, isReady: boolean) => {
    const player = players.value.find((p) => p.name === username);
    if (player) {
      player.isReady = isReady;
    }
  };

  const clearPlayers = () => {
    players.value = [];
  };

  return {
    players,
    updatePlayers,
    setPlayerReady,
    clearPlayers,
  };
});

// ゲーム結果管理ストア
export const useGameResultStore = defineStore("gameResultStore", () => {
  const players = ref<ResultPlayer[]>([]);

  const updatePlayers = (finalScores: Array<{ user_name: string; score: number }>) => {
    // スコア順でソートしてランクを計算
    const sortedScores = [...finalScores].sort((a, b) => b.score - a.score);

    players.value = sortedScores.map((player, index) => ({
      name: player.user_name,
      score: player.score,
      rank: index + 1,
    }));
  };

  const clearPlayers = () => {
    players.value = [];
  };

  return {
    players,
    updatePlayers,
    clearPlayers,
  };
});

// ルーム情報管理ストア
export const useCurrentRoomStore = defineStore("currentRoomStore", () => {
  const currentRoom = ref<Room | null>(null);

  const setCurrentRoom = (room: Room) => {
    currentRoom.value = room;
    // localStorageにも保存してPWAでの状態保持を確実にする
    localStorage.setItem("currentRoom", JSON.stringify(room));
  };

  const getCurrentRoom = () => {
    if (!currentRoom.value) {
      // メモリにない場合はlocalStorageから復元
      const stored = localStorage.getItem("currentRoom");
      if (stored) {
        try {
          currentRoom.value = JSON.parse(stored);
        } catch (error) {
          console.error("Failed to parse stored room data:", error);
        }
      }
    }
    return currentRoom.value;
  };

  const clearCurrentRoom = () => {
    currentRoom.value = null;
    localStorage.removeItem("currentRoom");
  };

  return {
    currentRoom,
    setCurrentRoom,
    getCurrentRoom,
    clearCurrentRoom,
  };
});

// WebSocket接続のグローバル管理
export const useWebSocketStore = defineStore("webSocketStore", () => {
  const wsManager = ref<ReturnType<typeof useWebSocket> | null>(null);
  const isConnected = ref(false);
  const currentUsername = ref<string>("");

  // 初期化時にlocalStorageからユーザー名を復元
  const initializeFromStorage = () => {
    const storedUsername = localStorage.getItem("username") || sessionStorage.getItem("username");
    if (storedUsername && !currentUsername.value) {
      currentUsername.value = storedUsername;
      console.log("Restored username from storage:", storedUsername);
      return true;
    }
    return false;
  };

  // WebSocket接続の自動復元
  const autoRestoreConnection = () => {
    if (initializeFromStorage() && !wsManager.value) {
      console.log("Auto-restoring WebSocket connection for:", currentUsername.value);
      initializeWebSocket(currentUsername.value);
      return true;
    }
    return false;
  };

  // WebSocket接続を初期化
  const initializeWebSocket = (username: string, onMessage?: (event: WebSocketEvent) => void) => {
    if (wsManager.value) {
      // 既存の接続がある場合は切断
      wsManager.value.destroy();
    }

    currentUsername.value = username;
    // localStorageとsessionStorageの両方にユーザー名を保存
    localStorage.setItem("username", username);
    sessionStorage.setItem("username", username);
    
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
    initializeFromStorage,
    autoRestoreConnection,
    initializeWebSocket,
    disconnectWebSocket,
    sendMessage,
    getWebSocketManager,
    getConnectionStatus,
  };
});
