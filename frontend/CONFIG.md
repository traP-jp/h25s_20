# フロントエンド設定管理

このプロジェクトでは、バックエンドAPIのURLを一元管理するための設定システムを導入しています。

## 設定ファイル

- `src/config/app.ts`: メインの設定ファイル
- `.env.example`: 環境変数の設定例

## 環境による自動切り替え

### 開発環境 (development)
- API Base URL: `http://localhost:8080/api`
- WebSocket URL: `ws://localhost:8080/api/ws`

### 本番環境 (production)
- API Base URL: `https://10ten.trap.show/api`
- WebSocket URL: `wss://10ten.trap.show/api/ws`

## カスタム設定

### 環境変数ファイルの読み込み順序

Viteは以下の順序で環境変数ファイルを読み込みます（後から読み込まれたものが優先）：

1. `.env` - 全ての環境で読み込まれる基本設定
2. `.env.local` - 全ての環境で読み込まれるローカル設定（Gitに含めない）
3. `.env.[mode]` - 特定のモード（development/production）でのみ読み込み
4. `.env.[mode].local` - 特定のモードでのローカル設定（Gitに含めない）

### 推奨される使用方法

```bash
# .env - チーム共通のデフォルト設定（Git管理対象）
# VITE_API_BASE_URL=https://api.example.com

# .env.local - 個人のローカル開発設定（Git管理対象外）
VITE_API_BASE_URL=http://localhost:8080/api
VITE_WS_BASE_URL=ws://localhost:8080/api/ws

# .env.development - 開発環境専用設定（Git管理対象）
# VITE_API_BASE_URL=https://dev-api.example.com

# .env.production - 本番環境専用設定（Git管理対象）
# VITE_API_BASE_URL=https://api.example.com
```

## 使用方法

### API設定の取得
```typescript
import { getConfig, getApiUrl, getWsUrl } from '@/config/app';

// 設定全体を取得
const config = getConfig();
console.log(config.api.baseUrl); // https://10ten.trap.show/api

// エンドポイント付きのURLを生成
const healthUrl = getApiUrl('/health'); // https://10ten.trap.show/api/health

// WebSocket URLを生成（パラメータ付き）
const wsUrl = getWsUrl('username=test'); // wss://10ten.trap.show/api/ws?username=test
```

### APIClientの使用
```typescript
import { apiClient } from '@/api';

// デフォルト設定で使用（自動的に環境に応じたURLが設定される）
const response = await apiClient.checkHealth();

// カスタムURLで使用
const customClient = new ApiClient('http://localhost:3000/api');
```

## 変更箇所

以下のファイルが更新されました：

1. `src/api.ts` - APIClientのコンストラクタが設定を自動取得
2. `src/store.ts` - WebSocket URLが設定から取得
3. `src/views/ThirdView.vue` - WebSocket URLが設定から取得
4. `src/views/ApiTestView.vue` - デフォルトURLが設定から取得

## 環境変数の優先順位

1. 環境変数 (`VITE_API_BASE_URL`, `VITE_WS_BASE_URL`)
2. 本番環境判定 (`import.meta.env.PROD`)
3. 開発環境 (デフォルト)

これにより、デプロイ時や開発時にURLを一括で変更することが可能になります。

### 実際の使用例

**開発時のローカル設定** (`.env.local`):
```bash
# 個人のローカル開発サーバー設定
VITE_API_BASE_URL=http://localhost:8080/api
VITE_WS_BASE_URL=ws://localhost:8080/api/ws
```

**本番デプロイ時** (`.env.production`):
```bash
# 本番環境設定
VITE_API_BASE_URL=https://10ten.trap.show/api
VITE_WS_BASE_URL=wss://10ten.trap.show/api/ws
```

### 設定の確認方法

開発中に現在の設定を確認するには：
```typescript
console.log('現在の設定:', import.meta.env);
console.log('API URL:', import.meta.env.VITE_API_BASE_URL);
```
