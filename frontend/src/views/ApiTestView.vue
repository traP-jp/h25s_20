<template>
  <div class="api-test">
    <div class="header">
      <h1>API Test Console</h1>
      
      <!-- 現在の設定情報表示 -->
      <div class="config-status">
        <h3>🔧 現在の設定状況</h3>
        <div class="status-grid">
          <div class="status-item">
            <span class="label">環境:</span>
            <span class="value" :class="environmentClass">{{ currentEnvironment }}</span>
          </div>
          <div class="status-item">
            <span class="label">API URL:</span>
            <span class="value">{{ config.api.baseUrl }}</span>
          </div>
          <div class="status-item">
            <span class="label">WebSocket URL:</span>
            <span class="value">{{ config.api.wsBaseUrl }}</span>
          </div>
          <div class="status-item">
            <span class="label">設定ソース:</span>
            <span class="value" :class="configSourceClass">{{ configSource }}</span>
          </div>
          <div class="status-item">
            <span class="label">実際の接続先:</span>
            <span class="value" :class="connectionTargetClass">{{ baseUrl || config.api.baseUrl }}</span>
          </div>
          <div class="status-item">
            <span class="label">接続状態:</span>
            <span class="value" :class="connectionStatusClass">{{ connectionStatus }}</span>
            <button @click="testConnection" :disabled="testing" class="test-btn">
              {{ testing ? '確認中...' : '接続確認' }}
            </button>
          </div>
        </div>
        <div class="env-vars">
          <h4>📋 環境変数</h4>
          <div class="env-list">
            <div v-for="(value, key) in envVars" :key="key" class="env-item">
              <span class="env-key">{{ key }}:</span>
              <span class="env-value">{{ value || '(未設定)' }}</span>
            </div>
          </div>
        </div>
      </div>

      <div class="server-config">
        <label>
          Base URL:
          <input v-model="baseUrl" type="text" :placeholder="config.api.baseUrl" />
          <small>空の場合は上記の設定を使用</small>
        </label>
        <label>
          JWT Token:
          <input v-model="authToken" type="text" placeholder="eyJhbGciOiJIUzI1NiIs..." />
        </label>
      </div>
    </div>

    <div class="api-sections">
      <!-- Health Check -->
      <div class="api-section">
        <h2>Health Check</h2>
        <button @click="testHealth" :disabled="isLoading">GET /health</button>
        <div v-if="responses.health" class="response" :class="responses.health.success ? 'success' : 'error'">
          <pre>{{ responses.health.data }}</pre>
        </div>
      </div>

      <!-- User Registration/Login -->
      <div class="api-section">
        <h2>User Management</h2>
        <div class="form-group">
          <input v-model="userData.username" type="text" placeholder="Username" />
          <input v-model="userData.password" type="password" placeholder="Password" />
          <button @click="testCreateUser" :disabled="isLoading">POST /users</button>
        </div>
        <div v-if="responses.users" class="response" :class="responses.users.success ? 'success' : 'error'">
          <pre>{{ responses.users.data }}</pre>
        </div>
      </div>

      <!-- Rooms -->
      <div class="api-section">
        <h2>Rooms</h2>
        <button @click="testGetRooms" :disabled="isLoading">GET /rooms</button>
        <div v-if="responses.rooms" class="response" :class="responses.rooms.success ? 'success' : 'error'">
          <pre>{{ responses.rooms.data }}</pre>
        </div>
      </div>

      <!-- Room Actions -->
      <div class="api-section">
        <h2>Room Actions</h2>
        <div class="form-group">
          <input v-model="roomAction.roomId" type="number" placeholder="Room ID" />
          <select v-model="roomAction.action">
            <option value="">Select Action</option>
            <option value="JOIN">JOIN</option>
            <option value="READY">READY</option>
            <option value="CANCEL">CANCEL</option>
            <option value="START">START</option>
            <option value="ABORT">ABORT</option>
            <option value="CLOSE_RESULT">CLOSE_RESULT</option>
          </select>
          <button @click="testRoomAction" :disabled="isLoading">POST /rooms/:id/actions</button>
        </div>
        <div
          v-if="responses.roomActions"
          class="response"
          :class="responses.roomActions.success ? 'success' : 'error'"
        >
          <pre>{{ responses.roomActions.data }}</pre>
        </div>
      </div>

      <!-- Formula Submission -->
      <div class="api-section">
        <h2>Formula Submission</h2>
        <div class="form-group">
          <input v-model="formula.roomId" type="number" placeholder="Room ID" />
          <input v-model="formula.version" type="number" placeholder="Version" />
          <input v-model="formula.formula" type="text" placeholder="Formula (e.g., 1+2*3-4)" />
          <button @click="testSubmitFormula" :disabled="isLoading">POST /rooms/:id/formulas</button>
        </div>
        <div v-if="responses.formulas" class="response" :class="responses.formulas.success ? 'success' : 'error'">
          <pre>{{ responses.formulas.data }}</pre>
        </div>
      </div>

      <!-- Room Results -->
      <div class="api-section">
        <h2>Room Results</h2>
        <div class="form-group">
          <input v-model="resultRoomId" type="number" placeholder="Room ID" />
          <button @click="testGetResults" :disabled="isLoading">GET /rooms/:id/result</button>
        </div>
        <div v-if="responses.results" class="response" :class="responses.results.success ? 'success' : 'error'">
          <pre>{{ responses.results.data }}</pre>
        </div>
      </div>
    </div>

    <!-- Loading indicator -->
    <div v-if="isLoading" class="loading">Testing API...</div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from "vue";
import { apiClient, type ApiResponse } from "@/api";
import { getConfig } from "@/config/app";

const isLoading = ref(false);
const testing = ref(false);
const connectionStatus = ref('未確認');
const config = getConfig();
const baseUrl = ref("");
const authToken = ref("");

// デバッグ情報の計算プロパティ
const currentEnvironment = computed(() => {
  // 環境変数で明示的に設定されているかチェック
  const hasCustomUrl = import.meta.env.VITE_API_BASE_URL;
  const isProd = import.meta.env.PROD;
  
  if (hasCustomUrl) {
    // 環境変数で設定されている場合、URLの内容で判定
    const customUrl = import.meta.env.VITE_API_BASE_URL;
    if (customUrl.includes('localhost') || customUrl.includes('127.0.0.1')) {
      return '開発環境 (env設定)';
    } else if (customUrl.includes('10ten.trap.show')) {
      return '本番環境 (env設定)';
    } else {
      return 'カスタム環境';
    }
  }
  
  return isProd ? '本番環境' : '開発環境';
});

const environmentClass = computed(() => {
  const env = currentEnvironment.value;
  if (env.includes('本番')) return 'env-production';
  if (env.includes('カスタム')) return 'env-custom';
  return 'env-development';
});

const configSource = computed(() => {
  if (import.meta.env.VITE_API_BASE_URL) {
    return '環境変数 (.env)';
  }
  if (import.meta.env.PROD) {
    return '自動設定 (本番)';
  }
  return '自動設定 (開発)';
});

const configSourceClass = computed(() => {
  const source = configSource.value;
  if (source.includes('環境変数')) return 'source-env';
  return 'source-file';
});

const connectionTargetClass = computed(() => {
  const target = baseUrl.value || config.api.baseUrl;
  if (target.includes('localhost') || target.includes('127.0.0.1')) {
    return 'target-local';
  }
  if (target.includes('10ten.trap.show')) {
    return 'target-production';
  }
  return 'target-other';
});

const envVars = computed(() => {
  return {
    'VITE_API_BASE_URL': import.meta.env.VITE_API_BASE_URL,
    'VITE_WS_BASE_URL': import.meta.env.VITE_WS_BASE_URL,
    'MODE': import.meta.env.MODE,
    'PROD': import.meta.env.PROD,
    'DEV': import.meta.env.DEV,
  };
});

const connectionStatusClass = computed(() => {
  const status = connectionStatus.value;
  if (status === '接続成功') return 'status-success';
  if (status === '接続失敗') return 'status-error';
  return 'status-unknown';
});

// 接続テスト関数
const testConnection = async () => {
  testing.value = true;
  connectionStatus.value = '確認中...';
  
  try {
    const testUrl = baseUrl.value || config.api.baseUrl;
    const response = await fetch(`${testUrl}/health`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });
    
    if (response.ok) {
      connectionStatus.value = '接続成功';
    } else {
      connectionStatus.value = '接続失敗';
    }
  } catch (error) {
    connectionStatus.value = '接続失敗';
  } finally {
    testing.value = false;
  }
};

// Initialize API client with reactive values
const updateApiClient = () => {
  apiClient.setBaseUrl(baseUrl.value);
  apiClient.setAuthToken(authToken.value);
};

const userData = reactive({
  username: "",
  password: "",
});

const roomAction = reactive({
  roomId: null as number | null,
  action: "" as "" | "JOIN" | "READY" | "CANCEL" | "START" | "ABORT" | "CLOSE_RESULT",
});

const formula = reactive({
  roomId: null as number | null,
  version: 1,
  formula: "",
});

const resultRoomId = ref<number | null>(null);

const responses = reactive({
  health: null as ApiResponse | null,
  users: null as ApiResponse | null,
  rooms: null as ApiResponse | null,
  roomActions: null as ApiResponse | null,
  formulas: null as ApiResponse | null,
  results: null as ApiResponse | null,
});

const testHealth = async () => {
  isLoading.value = true;
  updateApiClient();
  responses.health = await apiClient.checkHealth();
  isLoading.value = false;
};

const testCreateUser = async () => {
  if (!userData.username || !userData.password) {
    alert("Please enter username and password");
    return;
  }
  isLoading.value = true;
  updateApiClient();
  responses.users = await apiClient.createUser({
    username: userData.username,
    password: userData.password,
  });

  // If successful, update the auth token
  if (responses.users.success && responses.users.data?.token) {
    authToken.value = responses.users.data.token;
  }
  isLoading.value = false;
};

const testGetRooms = async () => {
  isLoading.value = true;
  updateApiClient();
  responses.rooms = await apiClient.getRooms();
  isLoading.value = false;
};

const testRoomAction = async () => {
  if (!roomAction.roomId || !roomAction.action) {
    alert("Please enter room ID and select an action");
    return;
  }
  isLoading.value = true;
  updateApiClient();
  responses.roomActions = await apiClient.performRoomAction(roomAction.roomId, { action: roomAction.action });
  isLoading.value = false;
};

const testSubmitFormula = async () => {
  if (!formula.roomId || !formula.formula) {
    alert("Please enter room ID and formula");
    return;
  }
  isLoading.value = true;
  updateApiClient();
  responses.formulas = await apiClient.submitFormula(formula.roomId, {
    version: formula.version,
    formula: formula.formula,
  });
  isLoading.value = false;
};

const testGetResults = async () => {
  if (!resultRoomId.value) {
    alert("Please enter room ID");
    return;
  }
  isLoading.value = true;
  updateApiClient();
  responses.results = await apiClient.getRoomResults(resultRoomId.value);
  isLoading.value = false;
};
</script>

<style scoped>
.api-test {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
}

.header {
  text-align: center;
  margin-bottom: 30px;
}

.header h1 {
  color: #2c3e50;
  margin-bottom: 20px;
}

.server-config {
  display: flex;
  gap: 20px;
  justify-content: center;
  flex-wrap: wrap;
  margin-bottom: 20px;
}

.server-config label {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  font-weight: 600;
  color: #555;
}

.server-config input {
  margin-top: 5px;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  width: 300px;
  font-size: 14px;
}

.server-config small {
  margin-top: 2px;
  font-size: 12px;
  color: #666;
  font-weight: normal;
}

/* デバッグ情報のスタイル */
.config-status {
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 20px;
  margin: 20px 0;
  text-align: left;
}

.config-status h3 {
  margin: 0 0 15px 0;
  color: #495057;
  font-size: 16px;
}

.status-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 12px;
  margin-bottom: 20px;
}

.status-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: white;
  border: 1px solid #dee2e6;
  border-radius: 4px;
  gap: 8px;
}

.status-item .label {
  font-weight: 600;
  color: #6c757d;
  font-size: 14px;
  flex-shrink: 0;
}

.status-item .value {
  font-family: Monaco, 'Courier New', monospace;
  font-size: 13px;
  font-weight: 500;
  padding: 2px 6px;
  border-radius: 3px;
  word-break: break-all;
  flex: 1;
}

.test-btn {
  padding: 4px 8px;
  background: #17a2b8;
  color: white;
  border: none;
  border-radius: 3px;
  cursor: pointer;
  font-size: 12px;
  flex-shrink: 0;
  transition: background-color 0.2s;
}

.test-btn:hover {
  background: #138496;
}

.test-btn:disabled {
  background: #6c757d;
  cursor: not-allowed;
}

/* 環境別の色分け */
.env-development {
  background: #d4edda;
  color: #155724;
}

.env-production {
  background: #f8d7da;
  color: #721c24;
}

.env-custom {
  background: #fff3cd;
  color: #856404;
}

/* 設定ソース別の色分け */
.source-env {
  background: #cce5ff;
  color: #004085;
}

.source-file {
  background: #e2e3e5;
  color: #383d41;
}

/* 接続先別の色分け */
.target-local {
  background: #d1ecf1;
  color: #0c5460;
}

.target-production {
  background: #f8d7da;
  color: #721c24;
}

.target-other {
  background: #fff3cd;
  color: #856404;
}

/* 接続状態別の色分け */
.status-success {
  background: #d4edda;
  color: #155724;
}

.status-error {
  background: #f8d7da;
  color: #721c24;
}

.status-unknown {
  background: #e2e3e5;
  color: #383d41;
}

/* 環境変数セクション */
.env-vars {
  border-top: 1px solid #dee2e6;
  padding-top: 15px;
}

.env-vars h4 {
  margin: 0 0 10px 0;
  color: #495057;
  font-size: 14px;
}

.env-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 8px;
}

.env-item {
  display: flex;
  justify-content: space-between;
  padding: 6px 10px;
  background: white;
  border: 1px solid #dee2e6;
  border-radius: 3px;
  font-size: 12px;
}

.env-key {
  font-weight: 600;
  color: #6c757d;
}

.env-value {
  font-family: Monaco, 'Courier New', monospace;
  color: #495057;
  max-width: 60%;
  word-break: break-all;
}

.api-sections {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 20px;
}

.api-section {
  border: 1px solid #e1e8ed;
  border-radius: 8px;
  padding: 20px;
  background: #f8f9fa;
}

.api-section h2 {
  margin: 0 0 15px 0;
  color: #2c3e50;
  font-size: 18px;
}

.form-group {
  display: flex;
  gap: 10px;
  margin-bottom: 15px;
  flex-wrap: wrap;
}

.form-group input,
.form-group select {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  flex: 1;
  min-width: 100px;
}

.form-group button,
.api-section > button {
  padding: 8px 16px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.2s;
}

.form-group button:hover,
.api-section > button:hover {
  background: #0056b3;
}

.form-group button:disabled,
.api-section > button:disabled {
  background: #6c757d;
  cursor: not-allowed;
}

.response {
  margin-top: 15px;
  padding: 15px;
  border-radius: 4px;
  font-family: "Monaco", "Courier New", monospace;
  font-size: 12px;
  white-space: pre-wrap;
  word-break: break-all;
  max-height: 300px;
  overflow-y: auto;
}

.response.success {
  background: #d4edda;
  border: 1px solid #c3e6cb;
  color: #155724;
}

.response.error {
  background: #f8d7da;
  border: 1px solid #f5c6cb;
  color: #721c24;
}

.loading {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: rgba(0, 0, 0, 0.8);
  color: white;
  padding: 20px 40px;
  border-radius: 8px;
  font-size: 16px;
  z-index: 1000;
}

@media (max-width: 768px) {
  .api-sections {
    grid-template-columns: 1fr;
  }

  .server-config {
    flex-direction: column;
    align-items: center;
  }

  .server-config input {
    width: 100%;
    max-width: 300px;
  }
}
</style>
