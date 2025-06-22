<template>
  <div class="api-test">
    <div class="header">
      <h1>API Test Console</h1>
      <div class="server-config">
        <label>
          Base URL:
          <input v-model="baseUrl" type="text" placeholder="https://10ten.trap.show/api" />
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
        <button @click="testHealth" :disabled="loading">GET /health</button>
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
          <button @click="testCreateUser" :disabled="loading">POST /users</button>
        </div>
        <div v-if="responses.users" class="response" :class="responses.users.success ? 'success' : 'error'">
          <pre>{{ responses.users.data }}</pre>
        </div>
      </div>

      <!-- Rooms -->
      <div class="api-section">
        <h2>Rooms</h2>
        <button @click="testGetRooms" :disabled="loading">GET /rooms</button>
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
          <button @click="testRoomAction" :disabled="loading">POST /rooms/:id/actions</button>
        </div>
        <div v-if="responses.roomActions" class="response" :class="responses.roomActions.success ? 'success' : 'error'">
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
          <button @click="testSubmitFormula" :disabled="loading">POST /rooms/:id/formulas</button>
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
          <button @click="testGetResults" :disabled="loading">GET /rooms/:id/result</button>
        </div>
        <div v-if="responses.results" class="response" :class="responses.results.success ? 'success' : 'error'">
          <pre>{{ responses.results.data }}</pre>
        </div>
      </div>
    </div>

    <!-- Loading indicator -->
    <div v-if="loading" class="loading">Testing API...</div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import axios, { type AxiosResponse, type AxiosError } from 'axios'

const loading = ref(false)
const baseUrl = ref('https://10ten.trap.show/api')
const authToken = ref('')

const userData = reactive({
  username: '',
  password: ''
})

const roomAction = reactive({
  roomId: null as number | null,
  action: ''
})

const formula = reactive({
  roomId: null as number | null,
  version: 1,
  formula: ''
})

const resultRoomId = ref<number | null>(null)

const responses = reactive({
  health: null as any,
  users: null as any,
  rooms: null as any,
  roomActions: null as any,
  formulas: null as any,
  results: null as any
})

const makeRequest = async (
  method: 'GET' | 'POST',
  endpoint: string,
  data?: any,
  needsAuth: boolean = false
) => {
  loading.value = true
  try {
    const config: any = {
      method,
      url: `${baseUrl.value}${endpoint}`,
      headers: {}
    }

    if (needsAuth && authToken.value) {
      config.headers.Authorization = `Bearer ${authToken.value}`
    }

    if (data) {
      config.data = data
      config.headers['Content-Type'] = 'application/json'
    }

    const response: AxiosResponse = await axios(config)
    return {
      success: true,
      status: response.status,
      data: response.data
    }
  } catch (error) {
    const axiosError = error as AxiosError
    return {
      success: false,
      status: axiosError.response?.status || 0,
      data: axiosError.response?.data || axiosError.message
    }
  } finally {
    loading.value = false
  }
}

const testHealth = async () => {
  responses.health = await makeRequest('GET', '/health')
}

const testCreateUser = async () => {
  if (!userData.username || !userData.password) {
    alert('Please enter username and password')
    return
  }
  responses.users = await makeRequest('POST', '/users', {
    username: userData.username,
    password: userData.password
  })
  
  // If successful, update the auth token
  if (responses.users.success && responses.users.data?.token) {
    authToken.value = responses.users.data.token
  }
}

const testGetRooms = async () => {
  responses.rooms = await makeRequest('GET', '/rooms', undefined, true)
}

const testRoomAction = async () => {
  if (!roomAction.roomId || !roomAction.action) {
    alert('Please enter room ID and select an action')
    return
  }
  responses.roomActions = await makeRequest(
    'POST',
    `/rooms/${roomAction.roomId}/actions`,
    { action: roomAction.action },
    true
  )
}

const testSubmitFormula = async () => {
  if (!formula.roomId || !formula.formula) {
    alert('Please enter room ID and formula')
    return
  }
  responses.formulas = await makeRequest(
    'POST',
    `/rooms/${formula.roomId}/formulas`,
    {
      version: formula.version,
      formula: formula.formula
    },
    true
  )
}

const testGetResults = async () => {
  if (!resultRoomId.value) {
    alert('Please enter room ID')
    return
  }
  responses.results = await makeRequest('GET', `/rooms/${resultRoomId.value}/result`, undefined, true)
}
</script>

<style scoped>
.api-test {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
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
  font-family: 'Monaco', 'Courier New', monospace;
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