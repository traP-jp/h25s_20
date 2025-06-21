# WebSocket実装ドキュメント

## 概要

このドキュメントでは、ゲームアプリケーションにおけるWebSocket実装について詳説します。本実装は、リアルタイムな通知と双方向通信を可能にし、ルーム機能、ゲーム状態の同期、プレイヤー間の通信をサポートします。

## アーキテクチャ概要

### システム構成

```
Frontend (Vue.js)
       ↓ WebSocket Connection
Backend (Go)
  ├── Handler Layer    (HTTP/WebSocket接続処理)
  ├── Use Case Layer   (ビジネスロジック)
  ├── Domain Layer     (ドメインモデル・インターフェース)
  └── Infrastructure   (WebSocket実装・Room管理)
```

### 設計原則

- **クリーンアーキテクチャ**: 依存関係の逆転原理に基づく層化構造
- **イベント駆動**: サーバー主導の一方向通知システム
- **並行処理**: Goroutineによる高効率な接続管理
- **フェイルセーフ**: 接続断絶時の自動クリーンアップ

## 実装詳細

### 1. ドメインレイヤー (Domain Layer)

#### WebSocketConnection インターフェース

```go
type WebSocketConnection interface {
    ReadMessage(ctx context.Context) (messageType int, data []byte, err error)
    WriteMessage(ctx context.Context, messageType int, data []byte) error
    Close() error
    GetUserID() string
    GetRoomID() string
    SetRoomID(roomID string)
}
```

#### WebSocketManager インターフェース

```go
type WebSocketManager interface {
    UpgradeConnection(w http.ResponseWriter, r *http.Request, userID string) (WebSocketConnection, error)
    AddConnection(roomID string, conn WebSocketConnection) error
    RemoveConnection(roomID string, userID string) error
    BroadcastToRoom(roomID string, messageType int, data []byte) error
    SendToUser(userID string, messageType int, data []byte) error
    GetRoomConnections(roomID string) []WebSocketConnection
    GetUserConnection(userID string) WebSocketConnection
}
```

#### メッセージ構造

```go
type WebSocketMessage struct {
    Type      string      `json:"type"`
    RoomID    string      `json:"roomId,omitempty"`
    UserID    string      `json:"userId,omitempty"`
    Data      interface{} `json:"data,omitempty"`
    Timestamp int64       `json:"timestamp"`
}
```

#### メッセージタイプ定数

| メッセージタイプ | 説明 |
|----------------|------|
| `room_update` | ルーム状態の更新 |
| `game_update` | ゲーム状態の更新 |
| `board_update` | ボード状態の更新 |
| `game_start` | ゲーム開始通知 |
| `game_end` | ゲーム終了通知 |
| `user_joined` | ユーザー参加通知 |
| `user_left` | ユーザー退出通知 |
| `notification` | 一般通知 |
| `error` | エラー通知 |
| `ping` | 接続確認（クライアント → サーバー） |
| `pong` | 接続確認応答（サーバー → クライアント） |

### 2. インフラストラクチャレイヤー (Infrastructure Layer)

#### WebSocket接続管理

**Connection構造体**:
```go
type connection struct {
    conn   *websocket.Conn  // 実際のWebSocket接続
    userID string           // ユーザーID
    roomID string           // 所属ルームID
    mu     sync.RWMutex     // 並行アクセス制御
}
```

**Manager構造体**:
```go
type manager struct {
    rooms map[string]map[string]domain.WebSocketConnection  // ルーム別接続管理
    users map[string]domain.WebSocketConnection             // ユーザー別接続管理
    mu    sync.RWMutex                                     // 並行アクセス制御
}
```

#### 接続管理機能

1. **接続のアップグレード**
   - HTTPリクエストをWebSocket接続にアップグレード
   - ユーザー認証とバリデーション
   - 既存接続の置き換え処理

2. **ルーム管理**
   - ユーザーのルーム参加・退出
   - ルームごとの接続追跡
   - 空ルームの自動削除

3. **メッセージ配信**
   - ルーム内全員への同報配信
   - 特定ユーザーへのメッセージ送信
   - 失敗した接続の自動クリーンアップ

### 3. ユースケースレイヤー (Use Case Layer)

#### WebSocketUsecase の主要機能

```go
type WebSocketUsecase struct {
    wsManager   domain.WebSocketManager
    roomManager domain.RoomManager
}
```

**主要メソッド**:

1. **HandleConnection**: 新しいWebSocket接続の処理
   - ウェルカムメッセージ送信
   - 接続維持とPing/Pong処理
   - エラー回復とログ出力

2. **JoinRoom**: ユーザーのルーム参加処理
   - プレイヤー作成
   - ルームマネージャーとの連携
   - 参加通知の配信

3. **LeaveRoom**: ユーザーのルーム退出処理
   - ルームマネージャーからの削除
   - WebSocket接続の切断
   - 退出通知の配信

4. **通知メソッド群**:
   - `NotifyRoomUpdate`: ルーム状態変更通知
   - `NotifyBoardUpdate`: ボード更新通知
   - `NotifyGameStart`: ゲーム開始通知
   - `NotifyGameEnd`: ゲーム終了通知

### 4. ハンドラーレイヤー (Handler Layer)

#### WebSocketHandler

```go
type WebSocketHandler struct {
    wsManager   domain.WebSocketManager
    roomManager domain.RoomManager
    wsUsecase   *usecase.WebSocketUsecase
}
```

**主要機能**:

1. **HandleWebSocket**: WebSocket接続エンドポイント
   - ユーザー認証（複数方式対応）
   - 接続のアップグレード
   - 接続処理の非同期実行
   - タイムアウト管理（30分）

2. **GetWebSocketStats**: 接続統計情報の取得
   - ルーム別接続数
   - アクティブユーザー一覧
   - デバッグ情報

## API仕様

### WebSocket接続エンドポイント

**URL**: `ws://localhost:8080/api/ws`

**認証方式**:
1. クエリパラメータ: `?userId=USER_ID`
2. HTTPヘッダー: `X-User-ID: USER_ID`
3. Authorizationヘッダー: `Authorization: Bearer TOKEN` (未実装)

**接続例**:
```javascript
// クエリパラメータ方式
const ws = new WebSocket('ws://localhost:8080/api/ws?userId=user123');

// ヘッダー方式（実装はブラウザ依存）
const ws = new WebSocket('ws://localhost:8080/api/ws', {
    headers: {
        'X-User-ID': 'user123'
    }
});
```

### 統計情報エンドポイント

**URL**: `GET /api/ws/stats`

**クエリパラメータ**:
- `roomId`: 特定ルームの統計情報を取得

**レスポンス例**:
```json
{
    "roomId": "1",
    "connectionCount": 3,
    "users": ["user1", "user2", "user3"]
}
```

## メッセージフォーマット

### サーバー → クライアント

#### ユーザー参加通知
```json
{
    "type": "user_joined",
    "roomId": "1",
    "userId": "user123",
    "data": {
        "action": "joined",
        "player": {
            "id": "user123",
            "name": "Player Name",
            "isReady": false,
            "score": 0
        },
        "room": {
            "id": 1,
            "name": "Game Room",
            "players": [...],
            "isOpened": true
        }
    },
    "timestamp": 1703123456
}
```

#### ゲーム開始通知
```json
{
    "type": "game_start",
    "roomId": "1",
    "data": {
        "room": {...},
        "update": {
            "status": "started"
        }
    },
    "timestamp": 1703123456
}
```

#### ボード更新通知
```json
{
    "type": "board_update",
    "roomId": "1",
    "data": {
        "room": {...},
        "update": {
            "version": 2,
            "board": [1, 2, 3, 4, 5, 6, 7, 8, 9, 0]
        }
    },
    "timestamp": 1703123456
}
```

#### エラー通知
```json
{
    "type": "error",
    "userId": "user123",
    "data": {
        "error": "エラーメッセージ"
    },
    "timestamp": 1703123456
}
```

### クライアント → サーバー

#### Ping（接続確認）
```json
{
    "type": "ping",
    "userId": "user123",
    "data": {},
    "timestamp": 1703123456
}
```

## 設定と制約

### 開発環境設定

```go
options := &websocket.AcceptOptions{
    InsecureSkipVerify: true,          // 開発環境のみ
    OriginPatterns:     []string{"*"}, // 全オリジン許可（開発環境のみ）
}
```

### タイムアウト設定

- **接続タイムアウト**: 30分
- **メッセージ送信タイムアウト**: 10秒（ブロードキャスト）、5秒（個別送信）
- **読み取りタイムアウト**: コンテキストベース

### 並行処理制限

- **最大接続数**: 制限なし（実装により調整可能）
- **ルーム数**: 制限なし
- **メッセージキュー**: インメモリ（永続化なし）

## エラーハンドリング

### 接続エラー処理

1. **接続失敗**
   - ログ出力とエラーレスポンス
   - 既存接続の適切なクリーンアップ

2. **メッセージ送信失敗**
   - 失敗した接続の自動削除
   - エラーログ出力
   - 他の接続への影響最小化

3. **パニック回復**
   - Goroutineレベルでのパニック回復
   - エラーログ出力
   - 接続の適切な終了

### ログ設定

```go
// 接続確立時
log.Info().Str("userID", userID).Msg("WebSocket connection established")

// エラー時
log.Error().Err(err).Str("userID", userID).Msg("Failed to send message to user")

// 警告時
log.Warn().Str("roomID", roomID).Int("failedCount", len(failedUsers)).Msg("Some connections failed during broadcast")
```

## 運用・監視

### 接続状態監視

1. **統計情報API**の活用
2. **ログ監視**による異常検知
3. **リソース使用量**の追跡

### スケーリング考慮事項

1. **メモリ使用量**: 接続数に比例
2. **CPU使用量**: メッセージ配信頻度に依存
3. **ネットワーク帯域**: メッセージサイズと頻度に依存

### セキュリティ考慮事項

1. **認証**: 現在は簡易実装、JWT導入予定
2. **認可**: ルーム参加権限の実装が必要
3. **レート制限**: メッセージ送信頻度の制限検討
4. **入力検証**: メッセージ内容の検証強化

## フロントエンド統合

### Vue.js での使用例

```javascript
// WebSocket接続の確立
const connectWebSocket = (userId) => {
    const ws = new WebSocket(`ws://localhost:8080/api/ws?userId=${userId}`);
    
    ws.onopen = () => {
        console.log('WebSocket connection established');
    };
    
    ws.onmessage = (event) => {
        const message = JSON.parse(event.data);
        handleWebSocketMessage(message);
    };
    
    ws.onclose = () => {
        console.log('WebSocket connection closed');
    };
    
    ws.onerror = (error) => {
        console.error('WebSocket error:', error);
    };
    
    return ws;
};

// メッセージハンドラー
const handleWebSocketMessage = (message) => {
    switch (message.type) {
        case 'user_joined':
            updateRoomUsers(message.data.room);
            break;
        case 'game_start':
            startGame(message.data);
            break;
        case 'board_update':
            updateBoard(message.data.update);
            break;
        case 'error':
            showError(message.data.error);
            break;
    }
};
```

### 再接続処理

```javascript
class WebSocketManager {
    constructor(userId) {
        this.userId = userId;
        this.ws = null;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
    }
    
    connect() {
        this.ws = connectWebSocket(this.userId);
        this.ws.onclose = () => this.handleDisconnect();
    }
    
    handleDisconnect() {
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
            setTimeout(() => {
                this.reconnectAttempts++;
                this.connect();
            }, 1000 * Math.pow(2, this.reconnectAttempts));
        }
    }
}
```

## 将来の拡張計画

### 短期的改善

1. **JWT認証**の実装
2. **レート制限**の導入
3. **メッセージ永続化**（Redis等）
4. **接続プール**の最適化

### 長期的拡張

1. **水平スケーリング**対応
2. **メッセージブローカー**の導入（RabbitMQ/Apache Kafka）
3. **リアルタイム分析**機能
4. **カスタムプロトコル**の実装

## トラブルシューティング

### よくある問題

1. **接続が頻繁に切断される**
   - ネットワーク環境の確認
   - タイムアウト設定の調整
   - Ping/Pong機能の活用

2. **メッセージが届かない**
   - 接続状態の確認
   - ログ出力の確認
   - 統計情報APIでの診断

3. **メモリリークが発生する**
   - 接続のクリーンアップ実装の確認
   - ルーム管理の適切性確認

### デバッグ方法

1. **ログレベルの調整**
2. **統計情報APIの活用**
3. **メモリプロファイリング**の実行

## まとめ

本WebSocket実装は、リアルタイムゲームアプリケーションに必要な基本機能を提供し、スケーラブルで保守性の高いアーキテクチャを採用しています。継続的な改善と拡張により、より堅牢で高性能なシステムへと発展させることができます。