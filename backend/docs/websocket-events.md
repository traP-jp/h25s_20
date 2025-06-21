# WebSocketイベント統一実装

このドキュメントでは、統一されたWebSocketイベント送信システムの実装について説明します。

## 概要

従来のWebSocketイベント発行は、`map[string]interface{}`を使った非構造化データと文字列ベースのイベント名により、以下の問題を抱えていました：

- **型安全性の欠如**: コンパイル時の型チェックができない
- **イベント名の分散管理**: ハードコードされた文字列が各所に散らばる
- **形式の不統一**: 一部で異なる構造体を使用

## 新しい統一実装

### 1. イベント名の定数化

全てのWebSocketイベント名を定数として管理：

```go
// backend/internal/infrastructure/websocket/events.go
const (
    // 接続関連
    EventConnection = "connection"
    
    // プレイヤー関連
    EventPlayerJoined   = "player_joined"
    EventPlayerReady    = "player_ready" 
    EventPlayerCanceled = "player_canceled"
    EventPlayerLeft     = "player_left"
    
    // ゲーム関連
    EventGameStarted    = "game_started"
    EventGameStart      = "game_start"
    EventCountdownStart = "countdown_start"
    EventCountdown      = "countdown"
    EventBoardUpdated   = "board_updated"
    EventResultClosed   = "result_closed"
)
```

### 2. 型安全なイベント構造体

全てのイベントコンテンツに対して型安全な構造体を定義：

```go
// 基本イベントコンテンツ
type BaseEventContent struct {
    UserID    int    `json:"user_id,omitempty"`
    UserName  string `json:"user_name,omitempty"`
    RoomID    int    `json:"room_id,omitempty"`
    Message   string `json:"message,omitempty"`
    Timestamp int64  `json:"timestamp,omitempty"`
}

// プレイヤーアクション用
type PlayerEventContent struct {
    BaseEventContent
}

// ボード更新用
type BoardUpdateEventContent struct {
    BaseEventContent
    Board     BoardData `json:"board"`
    GainScore int       `json:"gain_score"`
}

// ボードデータ構造体
type BoardData struct {
    Content []int `json:"content"`
    Version int   `json:"version"`
    Size    int   `json:"size"`
}
```

### 3. ヘルパー関数

イベント作成を簡単にするヘルパー関数群：

```go
func NewPlayerEvent(eventType string, userID int, userName string, roomID int) WebSocketEvent
func NewBoardUpdateEvent(userID int, userName string, roomID int, board BoardData, gainScore int) WebSocketEvent
func NewGameStartBoardEvent(roomID int, message string, board BoardData) WebSocketEvent
```

### 4. 統一された送信メソッド

WebSocketハンドラーに統一されたイベント送信メソッドを追加：

```go
// プレイヤーイベント送信
func (h *WebSocketHandler) SendPlayerEventToRoom(roomID int, eventType string, userID int, userName string)

// ボード更新イベント送信
func (h *WebSocketHandler) SendBoardUpdateEventTyped(roomID int, userID int, userName string, board BoardData, gainScore int)

// ゲーム開始ボードイベント送信
func (h *WebSocketHandler) SendGameStartBoardEventToRoom(roomID int, message string, board BoardData)
```

## 利用例

### 旧実装（非推奨）

```go
h.WebSocketHandler.BroadcastToRoom(roomId, "player_joined", map[string]interface{}{
    "user_id":   player.ID,
    "user_name": player.UserName,
    "room_id":   roomId,
})
```

### 新実装（推奨）

```go
h.WebSocketHandler.SendPlayerEventToRoom(
    roomId,
    wsManager.EventPlayerJoined,
    player.ID,
    player.UserName,
)
```

## メリット

1. **型安全性**: コンパイル時の型チェックにより、実行時エラーを防止
2. **一覧性**: 全てのイベント名が定数として一箇所で管理される
3. **統一性**: 全てのイベントが同じ構造とパターンで発行される
4. **保守性**: イベント名やフィールドの変更が容易
5. **フロントエンド連携**: 型定義を共有することでフロントエンドとの連携が向上

## 移行状況

### 完了済み

- ✅ `player_joined` - プレイヤー参加通知
- ✅ `player_ready` - プレイヤー準備完了通知
- ✅ `player_canceled` - プレイヤー準備キャンセル通知
- ✅ `player_left` - プレイヤー退出通知
- ✅ `game_started` - ゲーム開始通知
- ✅ `result_closed` - 結果表示終了通知
- ✅ `connection` - 接続確立通知
- ✅ `countdown_start` - カウントダウン開始通知
- ✅ `countdown` - カウントダウン中通知
- ✅ `board_updated` - ボード更新通知
- ✅ `game_start` - ゲーム実開始（ボード付き）通知

### 後方互換性

旧来の `NotificationMessage` や `StandardEventContent` は非推奨として残されており、段階的な移行が可能です。

## フロントエンド連携のための型定義

TypeScriptでの型定義例：

```typescript
// WebSocketイベントの基本構造
interface WebSocketEvent {
  event: string;
  content: EventContent;
}

// 基本イベントコンテンツ
interface BaseEventContent {
  user_id?: number;
  user_name?: string;
  room_id?: number;
  message?: string;
  timestamp?: number;
}

// プレイヤーイベント
interface PlayerEventContent extends BaseEventContent {}

// ボード更新イベント
interface BoardUpdateEventContent extends BaseEventContent {
  board: BoardData;
  gain_score: number;
}

interface BoardData {
  content: number[];
  version: number;
  size: number;
}
```

これにより、フロントエンドでも型安全なWebSocketイベント処理が可能になります。 