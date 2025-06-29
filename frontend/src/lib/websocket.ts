import { ref } from "vue";

// WebSocketイベントの型定義
export interface WebSocketEvent {
  event: string;
  content: EventContent;
}

export interface BaseEventContent {
  user_id?: number;
  user_name?: string;
  room_id?: number;
  message?: string;
  timestamp?: number;
}

export interface ConnectionEventContent extends BaseEventContent {
  client_id: string;
}

export interface PlayerEventContent extends BaseEventContent {}

export interface PlayerJoinedEventContent extends BaseEventContent {
  room: RoomInfo;
}

export interface PlayerLeftEventContent extends BaseEventContent {
  room: RoomInfo;
}

export interface RoomInfo {
  id: number;
  name: string;
  state: string;
  is_opened: boolean;
  players: PlayerInfo[];
}

export interface PlayerInfo {
  id: number;
  user_name: string;
  is_ready: boolean;
  has_closed_result: boolean;
  score: number;
}

export interface BoardData {
  content: number[];
  version: number;
  size: number;
}

export interface BoardUpdateEventContent extends BaseEventContent {
  board: BoardData;
  gain_score: number;
}

export interface CountdownEventContent extends BaseEventContent {
  count?: number;
  countdown?: number;
}

export interface GameEndEventContent extends BaseEventContent {
  final_scores?: Array<{
    user_id: number;
    user_name: string;
    score: number;
  }>;
}

export interface RoomStateEventContent extends BaseEventContent {
  state?: string;
  players?: Array<{
    user_id: number;
    user_name: string;
    is_ready: boolean;
    // 実績データは現在のAPIでは提供されていないため、デフォルト値を使用
  }>;
}

export interface RoomClosedEventContent extends BaseEventContent {}

export type EventContent =
  | ConnectionEventContent
  | PlayerEventContent
  | PlayerJoinedEventContent
  | PlayerLeftEventContent
  | BoardUpdateEventContent
  | CountdownEventContent
  | GameEndEventContent
  | RoomStateEventContent
  | RoomClosedEventContent
  | BaseEventContent;

// WebSocketイベント名の定数
export const WS_EVENTS = {
  CONNECTION: "connection",
  PLAYER_JOINED: "player_joined",
  PLAYER_READY: "player_ready",
  PLAYER_CANCELED: "player_canceled",
  PLAYER_ALL_READY: "player_all_ready",
  PLAYER_LEFT: "player_left",
  ROOM_STATE_CHANGED: "room_state_changed",
  ROOM_CLOSED: "room_closed",
  GAME_STARTED: "game_started",
  GAME_START: "game_start",
  COUNTDOWN_START: "countdown_start",
  COUNTDOWN: "countdown",
  BOARD_UPDATED: "board_updated",
  RESULT_CLOSED: "result_closed",
  GAME_ENDED: "game_ended",
} as const;

// WebSocket接続管理クラス
export class WebSocketManager {
  private ws: WebSocket | null = null;
  private reconnectAttempts = 0;
  private readonly maxReconnectAttempts = 5;
  private connectTimeout: ReturnType<typeof setTimeout> | null = null;

  // リアクティブな状態
  public isConnected = ref(false);
  public isConnecting = ref(false);
  public connectionError = ref<string | null>(null);
  public messages = ref<string[]>([]);

  constructor(wsUrl: string, onMessage?: (event: WebSocketEvent) => void) {
    this.wsUrl = wsUrl;
    this.onMessage = onMessage;
  }

  private wsUrl: string;
  private onMessage?: (event: WebSocketEvent) => void;

  // WebSocket接続関数
  connect(): void {
    if (this.isConnecting.value || (this.ws && this.ws.readyState === WebSocket.OPEN)) {
      return;
    }

    this.isConnecting.value = true;
    this.connectionError.value = null;

    try {
      this.ws = new WebSocket(this.wsUrl);

      // 接続開始のタイムアウト設定
      this.connectTimeout = setTimeout(() => {
        if (this.ws && this.ws.readyState === WebSocket.CONNECTING) {
          this.ws.close();
          this.connectionError.value = "接続タイムアウト";
          this.isConnecting.value = false;
        }
      }, 10000); // 10秒でタイムアウト

      this.ws.onopen = () => {
        if (this.connectTimeout) {
          clearTimeout(this.connectTimeout);
          this.connectTimeout = null;
        }
        console.log("WebSocket接続が確立されました");
        this.isConnected.value = true;
        this.isConnecting.value = false;
        this.reconnectAttempts = 0;
        this.connectionError.value = null;

        this.addMessage("✅ WebSocket接続が確立されました");
      };

      this.ws.onmessage = (event) => {
        try {
          const wsEvent: WebSocketEvent = JSON.parse(event.data);
          this.handleWebSocketEvent(wsEvent);

          // 外部ハンドラーがある場合は呼び出し
          if (this.onMessage) {
            this.onMessage(wsEvent);
          }
        } catch (error) {
          console.error("WebSocketメッセージの解析に失敗:", error);
          this.addMessage(`❌ メッセージ解析エラー: ${event.data}`);
        }
      };

      this.ws.onerror = (event) => {
        if (this.connectTimeout) {
          clearTimeout(this.connectTimeout);
          this.connectTimeout = null;
        }
        console.error("WebSocketエラー:", event);
        this.connectionError.value = "WebSocket接続エラーが発生しました";
        this.isConnecting.value = false;

        this.addMessage("❌ WebSocket接続エラー");
      };

      this.ws.onclose = (event) => {
        if (this.connectTimeout) {
          clearTimeout(this.connectTimeout);
          this.connectTimeout = null;
        }
        console.log(`WebSocket接続が閉じられました - Code: ${event.code}, Reason: ${event.reason}`);
        this.isConnected.value = false;
        this.isConnecting.value = false;

        if (event.code !== 1000) {
          // 正常終了以外の場合
          this.connectionError.value = `接続が閉じられました (Code: ${event.code})`;
          this.addMessage(`⚠️ 接続が閉じられました (Code: ${event.code})`);

          // 自動再接続を試行
          this.attemptReconnect();
        } else {
          this.addMessage("🔌 WebSocket接続が正常に閉じられました");
        }
      };
    } catch (error) {
      console.error("WebSocket接続の作成に失敗:", error);
      this.connectionError.value = "WebSocket接続の作成に失敗しました";
      this.isConnecting.value = false;
      this.addMessage("❌ WebSocket接続の作成に失敗");
    }
  }

  // WebSocketイベントハンドラー
  private handleWebSocketEvent(wsEvent: WebSocketEvent): void {
    console.log("受信イベント:", wsEvent);

    switch (wsEvent.event) {
      case WS_EVENTS.CONNECTION:
        const connectionContent = wsEvent.content as ConnectionEventContent;
        this.addMessage(`🔗 接続確立: ClientID: ${connectionContent.client_id}, UserID: ${connectionContent.user_id}`);
        break;

      case WS_EVENTS.PLAYER_JOINED:
        const joinedContent = wsEvent.content as PlayerJoinedEventContent;
        this.addMessage(
          `👤 プレイヤー参加: ${joinedContent.user_name} (ID: ${joinedContent.user_id}) がルーム ${joinedContent.room_id} に参加\n` +
            `🏠 ルーム状態: ${joinedContent.room.state}, プレイヤー数: ${joinedContent.room.players.length}`
        );
        break;

      case WS_EVENTS.PLAYER_READY:
        const readyContent = wsEvent.content as PlayerEventContent;
        this.addMessage(`✅ プレイヤー準備完了: ${readyContent.user_name} が準備完了`);
        break;

      case WS_EVENTS.PLAYER_CANCELED:
        const canceledContent = wsEvent.content as PlayerEventContent;
        this.addMessage(`❌ プレイヤー準備キャンセル: ${canceledContent.user_name} が準備をキャンセル`);
        break;

      case WS_EVENTS.PLAYER_ALL_READY:
        const allReadyContent = wsEvent.content as PlayerEventContent;
        this.addMessage(`🎉 全プレイヤー準備完了: ${allReadyContent.message || "All players are ready!"}`);
        break;

      case WS_EVENTS.PLAYER_LEFT:
        const leftContent = wsEvent.content as PlayerLeftEventContent;
        this.addMessage(
          `👋 プレイヤー退出: ${leftContent.user_name} がルームから退出\n` +
            `🏠 ルーム状態: ${leftContent.room.state}, プレイヤー数: ${leftContent.room.players.length}`
        );
        break;

      case WS_EVENTS.ROOM_STATE_CHANGED:
        const roomStateContent = wsEvent.content as RoomStateEventContent;
        this.addMessage(
          `🏠 部屋状態変更: ${roomStateContent.state}, プレイヤー数: ${roomStateContent.players?.length || 0}`
        );
        break;

      case WS_EVENTS.ROOM_CLOSED:
        const roomClosedContent = wsEvent.content as RoomClosedEventContent;
        this.addMessage(`🔒 ルームクローズ: ${roomClosedContent.message || "Room has been closed"}`);
        break;

      case WS_EVENTS.GAME_STARTED:
        const gameStartedContent = wsEvent.content as BaseEventContent;
        this.addMessage(`🎮 ゲーム開始: ${gameStartedContent.message}`);
        break;

      case WS_EVENTS.COUNTDOWN_START:
        const countdownStartContent = wsEvent.content as CountdownEventContent;
        this.addMessage(`⏰ カウントダウン開始: ${countdownStartContent.countdown}秒`);
        break;

      case WS_EVENTS.COUNTDOWN:
        const countdownContent = wsEvent.content as CountdownEventContent;
        this.addMessage(`⏱️ カウントダウン: ${countdownContent.count}`);
        break;

      case WS_EVENTS.BOARD_UPDATED:
        const boardContent = wsEvent.content as BoardUpdateEventContent;
        this.addMessage(
          `📋 ボード更新: ${boardContent.user_name} がスコア ${boardContent.gain_score} 獲得 (Version: ${boardContent.board.version})`
        );
        break;

      case WS_EVENTS.GAME_ENDED:
        const gameEndedContent = wsEvent.content as GameEndEventContent;
        this.addMessage(`🏁 ゲーム終了: ${gameEndedContent.message}`);
        if (gameEndedContent.final_scores) {
          this.addMessage(
            `📊 最終スコア: ${gameEndedContent.final_scores.map((s) => `${s.user_name}: ${s.score}点`).join(", ")}`
          );
        }
        break;

      default:
        this.addMessage(`📨 未知のイベント: ${wsEvent.event} - ${JSON.stringify(wsEvent.content)}`);
    }
  }

  // メッセージ追加関数
  private addMessage(message: string): void {
    const timestamp = new Date().toLocaleTimeString();
    this.messages.value.push(`[${timestamp}] ${message}`);

    // メッセージ数を制限（パフォーマンス対策）
    if (this.messages.value.length > 100) {
      this.messages.value = this.messages.value.slice(-50);
    }
  }

  // 再接続試行
  private attemptReconnect(): void {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      this.addMessage("❌ 最大再接続試行回数に達しました");
      return;
    }

    this.reconnectAttempts++;
    const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts - 1), 30000); // 指数バックオフ（最大30秒）

    this.addMessage(
      `🔄 ${delay / 1000}秒後に再接続を試行します... (${this.reconnectAttempts}/${this.maxReconnectAttempts})`
    );

    setTimeout(() => {
      this.connect();
    }, delay);
  }

  // 手動接続関数
  manualConnect(): void {
    if (this.ws) {
      this.ws.close();
    }
    this.reconnectAttempts = 0;
    this.connect();
  }

  // 切断関数
  disconnect(): void {
    if (this.ws) {
      this.ws.close(1000, "ユーザーによる切断");
    }
    this.reconnectAttempts = this.maxReconnectAttempts; // 自動再接続を停止
  }

  // メッセージクリア
  clearMessages(): void {
    this.messages.value = [];
  }

  // WebSocketメッセージ送信
  send(data: any): void {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(data));
    } else {
      console.error("WebSocket is not connected");
      this.addMessage("❌ WebSocketが接続されていません");
    }
  }

  // リソースクリーンアップ
  destroy(): void {
    if (this.connectTimeout) {
      clearTimeout(this.connectTimeout);
      this.connectTimeout = null;
    }
    if (this.ws) {
      this.ws.close(1000, "コンポーネントアンマウント");
    }
  }

  // 再接続試行回数の取得
  getReconnectAttempts(): number {
    return this.reconnectAttempts;
  }

  // 最大再接続試行回数の取得
  getMaxReconnectAttempts(): number {
    return this.maxReconnectAttempts;
  }

  // イベントハンドラーの設定/追加
  setMessageHandler(handler: (event: WebSocketEvent) => void): void {
    this.onMessage = handler;
  }

  // 既存のイベントハンドラーと新しいハンドラーを組み合わせる
  addMessageHandler(handler: (event: WebSocketEvent) => void): void {
    const existingHandler = this.onMessage;
    this.onMessage = (event: WebSocketEvent) => {
      if (existingHandler) {
        existingHandler(event);
      }
      handler(event);
    };
  }
}

// WebSocket接続用のComposable関数
export function useWebSocket(wsUrl: string, onMessage?: (event: WebSocketEvent) => void) {
  const manager = new WebSocketManager(wsUrl, onMessage);

  return {
    // 状態
    isConnected: manager.isConnected,
    isConnecting: manager.isConnecting,
    connectionError: manager.connectionError,
    messages: manager.messages,

    // メソッド
    connect: () => manager.connect(),
    manualConnect: () => manager.manualConnect(),
    disconnect: () => manager.disconnect(),
    clearMessages: () => manager.clearMessages(),
    send: (data: any) => manager.send(data),
    destroy: () => manager.destroy(),
    getReconnectAttempts: () => manager.getReconnectAttempts(),
    getMaxReconnectAttempts: () => manager.getMaxReconnectAttempts(),
    setMessageHandler: (handler: (event: WebSocketEvent) => void) => manager.setMessageHandler(handler),
    addMessageHandler: (handler: (event: WebSocketEvent) => void) => manager.addMessageHandler(handler),
  };
}
