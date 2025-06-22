# 4x4数字パズルゲーム - ゲームフロー仕様書

## 概要

本ドキュメントは、WebSocketとREST APIを組み合わせたマルチプレイヤー4x4数字パズルゲームのゲームフロー仕様を定義する。

## ゲームの基本ルール

- 4x4の盤面（1-9の数字）
- プレイヤーは数式（四則演算）で10を作成
- 使用できる数字は4つで、盤面の特定の領域から取得
- マッチング対象：行（4つ）、列（4つ）、対角線（2つ）、2×2ブロック（4つ）
- 各マッチでスコア+10点
- 制限時間：120秒

## ゲームフロー

### 1. 部屋選択フェーズ
- プレイヤーは利用可能な部屋一覧を取得
- 参加したい部屋を選択して`JOIN`アクション送信
- 部屋が満員の場合は失敗

### 2. 待機フェーズ（StateWaitingForPlayers）
- プレイヤーは`READY`または`CANCEL`アクションを送信可能
- 全員が`READY`状態になると`StateAllReady`に遷移
- 他プレイヤーの準備状況をリアルタイム表示

### 3. ゲーム開始準備（StateAllReady）
- 列の先頭プレイヤーのみが`START`アクション送信可能
- `START`送信により`StateCountdown`に遷移

### 4. カウントダウンフェーズ（StateCountdown）
- 3秒間のカウントダウン実行
- 自動的に`StateGameInProgress`に遷移
- 盤面データの配信開始

### 5. ゲームプレイフェーズ（StateGameInProgress）
- プレイヤーは`SUBMIT_FORMULA`アクションで数式送信
- バージョン管理による楽観的ロッキング
- マッチした領域は自動更新（新しいランダム数字）
- 120秒でゲーム終了

### 6. 結果表示フェーズ（StateGameEnded）
- 最終スコアを全プレイヤーに配信
- プレイヤーは`CLOSE_RESULT`アクションで結果画面を閉じる
- 全員が閉じると自動的に`StateWaitingForPlayers`にリセット

### 7. 中断・退出処理
- `ABORT`アクションで強制リセット
- プレイヤー切断時の自動処理

## API仕様

### REST API エンドポイント

#### 部屋一覧取得
```
GET /api/rooms
Response: {
  "rooms": [
    {
      "id": 1,
      "name": "部屋1",
      "player_count": 2,
      "max_players": 4,
      "state": "WaitingForPlayers"
    }
  ]
}
```

#### 部屋詳細取得
```
GET /api/rooms/{id}
Response: {
  "id": 1,
  "name": "部屋1",
  "players": [...],
  "state": "WaitingForPlayers",
  "current_board": {...}
}
```

### WebSocket イベント

#### プレイヤー → サーバー

| アクション | フェーズ | パラメータ | 説明 |
|-----------|---------|-----------|------|
| `JOIN` | 部屋選択 | `room_id` | 部屋への参加 |
| `READY` | 待機中 | なし | 準備完了 |
| `CANCEL` | 待機中 | なし | 準備キャンセル |
| `START` | 全員準備完了 | なし | ゲーム開始 |
| `SUBMIT_FORMULA` | ゲーム中 | `expression`, `version` | 数式送信 |
| `CLOSE_RESULT` | 結果表示 | なし | 結果画面を閉じる |
| `ABORT` | 任意 | なし | ゲーム中断 |

#### サーバー → プレイヤー

| イベント | タイミング | データ | 説明 |
|---------|-----------|--------|------|
| `ROOM_STATE_CHANGED` | 状態変更時 | `state`, `players` | 部屋状態更新 |
| `BOARD_UPDATE` | 盤面変更時 | `board`, `version` | 盤面更新 |
| `COUNTDOWN` | カウントダウン中 | `count` | カウントダウン表示 |
| `GAME_STARTED` | ゲーム開始時 | `board`, `start_time` | ゲーム開始通知 |
| `FORMULA_RESULT` | 数式送信後 | `success`, `message`, `score` | 数式結果 |
| `GAME_ENDED` | ゲーム終了時 | `final_scores` | 最終結果 |
| `PLAYER_ACTION` | プレイヤー行動時 | `player_id`, `action` | 他プレイヤーの行動 |

## 数式計算システム

### 新しい安全な実装（Go版）

セキュリティを強化した新しい数式計算システムを実装：

```go
type FormulaCalculator struct{}

func (fc *FormulaCalculator) EvaluateFormula(expression string) (float64, error) {
    // 1. 入力サニタイズ
    // 2. 数字チェック（1-9が4つ）
    // 3. 逆ポーランド記法に変換
    // 4. 安全な計算実行
    // 5. 結果検証
}
```

**主な特徴：**
- 逆ポーランド記法による安全計算
- 厳密な入力検証（1-9の数字4つのみ）
- ゼロ除算・不正演算子の適切な処理
- フロントエンドとの実装統一

### 不可能な数字組み合わせ

10を作成することが不可能な4つの数字の組み合わせ：
```
1111, 1112, 1113, 1122, 1159, 1169, 1177, 1178, 1179, 1188,
1399, 1444, 1499, 1666, 1667, 1677, 1699, 1777, 2257, 3444,
3669, 3779, 3999, 4444, 4459, 4477, 4558, 4899, 4999, 5668,
5788, 5799, 5899, 6666, 6667, 6677, 6777, 6778, 6888, 6899,
6999, 7777, 7788, 7789, 7799, 7888, 7999, 8899
```

## バージョン管理アルゴリズム

### 楽観的ロッキング

```go
func AttemptMoveWithVersion(gb *GameBoard, expression string, submittedVersion int) (bool, string, int) {
    // 1. 数式マッチング検証
    matches, found := FindAllMatchingLinesWithSets(gb, expression)
    
    // 2. バージョン衝突チェック
    hasConflict, conflictMsg := gb.CheckConflictWithPositions(submittedVersion, matches)
    
    // 3. 数式計算・検証
    calculator := NewFormulaCalculator()
    result, err := calculator.EvaluateFormula(expression)
    
    // 4. 盤面更新
    gb.UpdateLinesWithPositions(matches)
}
```

### 衝突検出ロジック

```go
func (gb *GameBoard) CheckConflictWithPositions(submittedVersion int, formulaMatches []Matches) (bool, string) {
    // 提出バージョン以降の変更マス位置を収集
    changedPositions := make(map[string]bool)
    for version := submittedVersion + 1; version <= gb.Version; version++ {
        // 変更履歴から衝突チェック
    }
    
    // 数式使用マス位置との重複検出
    for _, match := range formulaMatches {
        for _, pos := range match.Positions {
            if changedPositions[key] {
                return true, "衝突エラー"
            }
        }
    }
}
```

## 特定された問題点と対応策

### 重大な問題

#### 1. 認証システムの不備
**問題：** 固定のmockPlayerを使用  
**影響：** セキュリティリスク、不正アクセス可能  
**修正案：** JWT認証、セッション管理の実装

#### 2. ゲーム終了条件の不明確性
**問題：** 120秒タイマーのみで終了条件が不十分  
**影響：** ゲームバランスの問題  
**修正案：** 
- 全マス消去での即座終了
- スコア上限到達での終了
- 早期終了条件の追加

#### 3. スコア計算の単純性
**問題：** マッチ数×10点のみ  
**影響：** 戦略性の欠如  
**修正案：**
- 連続マッチボーナス
- 時間ボーナス
- 難易度別スコア

### 警告レベルの問題

#### 4. WebSocketタイムアウト設定
**問題：** 30秒は短すぎる可能性  
**修正案：** 環境変数での設定可能化

#### 5. エラーハンドリングの不統一
**問題：** 日本語・英語メッセージが混在  
**修正案：** 国際化対응（i18n）システム

#### 6. フロントエンド連携の不備
**問題：** sample-dataを使用している部分  
**修正案：** 実際のWebSocket通信への移行

## 改善提案

### 1. パフォーマンス向上
- 盤面更新の差分送信
- WebSocketメッセージの圧縮
- データベースクエリの最適化

### 2. ユーザビリティ改善
- リアルタイム盤面同期の視覚的フィードバック
- 操作ヒント機能
- 観戦モード

### 3. セキュリティ強化
- ✅ **完了**: 新しい数式計算システム実装
- 入力値検証の強化
- レート制限の実装

### 4. 拡張性向上
- 部屋設定のカスタマイズ機能
- 異なる盤面サイズ対応
- トーナメントモード

## 実装状況

### ✅ 完了
- ゲームフロー仕様の明確化
- 新しい数式計算システム（Go実装）
- バージョン管理システムの詳細仕様
- 問題点の特定と分析

### 🔄 進行中
- 認証システムの実装検討
- フロントエンド連携の改善
- エラーハンドリングの統一

### 📋 計画中
- パフォーマンス改善
- セキュリティ強化
- 機能拡張

---

**最終更新日：** 2024年12月現在  
**バージョン：** 1.0  
**ステータス：** 実装完了（Go実装版） 