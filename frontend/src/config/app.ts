// 環境設定
export interface AppConfig {
  api: {
    baseUrl: string;
    wsBaseUrl: string;
  };
}

// 環境別設定
const configs: Record<string, AppConfig> = {
  development: {
    api: {
      baseUrl: 'http://localhost:8080/api',
      wsBaseUrl: 'ws://localhost:8080/api/ws',
    },
  },
  production: {
    api: {
      baseUrl: 'https://10ten.trap.show/api',
      wsBaseUrl: 'wss://10ten.trap.show/api/ws',
    },
  },
};

// 現在の環境を判定
const getCurrentEnvironment = (): string => {
  // Viteの環境変数を使用
  if (import.meta.env.VITE_API_BASE_URL) {
    return 'custom';
  }
  
  // 本番環境の判定（ビルド時やドメインベース）
  if (import.meta.env.PROD) {
    return 'production';
  }
  
  return 'development';
};

// カスタム設定（環境変数から）
const getCustomConfig = (): AppConfig => {
  const baseUrl = import.meta.env.VITE_API_BASE_URL || configs.production.api.baseUrl;
  const wsBaseUrl = import.meta.env.VITE_WS_BASE_URL || 
    baseUrl.replace(/^https?:/, baseUrl.startsWith('https:') ? 'wss:' : 'ws:') + '/ws';
  
  return {
    api: {
      baseUrl,
      wsBaseUrl,
    },
  };
};

// 設定取得
export const getConfig = (): AppConfig => {
  const env = getCurrentEnvironment();
  
  if (env === 'custom') {
    return getCustomConfig();
  }
  
  return configs[env] || configs.development;
};

// 便利な関数
export const getApiUrl = (endpoint: string = ''): string => {
  const config = getConfig();
  return `${config.api.baseUrl}${endpoint}`;
};

export const getWsUrl = (params: string = ''): string => {
  const config = getConfig();
  return `${config.api.wsBaseUrl}${params ? `?${params}` : ''}`;
};

// デフォルトエクスポート
export default getConfig();
