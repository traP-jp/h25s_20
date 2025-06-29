import axios, { type AxiosResponse, type AxiosError } from "axios";

// API response type
export interface ApiResponse<T = any> {
  success: boolean;
  status: number;
  data: T;
}

// API request types
export interface UserData {
  username: string;
  password: string;
}

export interface RoomAction {
  action: "JOIN" | "READY" | "CANCEL" | "START" | "ABORT" | "CLOSE_RESULT";
}

export interface FormulaSubmission {
  version: number;
  formula: string;
}

// API configuration
export class ApiClient {
  private baseUrl: string;
  private authToken: string;

  constructor(baseUrl: string = "https://10ten.trap.show/api", authToken: string = "") {
    this.baseUrl = baseUrl;
    this.authToken = authToken;
  }

  setBaseUrl(url: string) {
    this.baseUrl = url;
  }

  setAuthToken(token: string) {
    this.authToken = token;
  }

  private async makeRequest<T = any>(
    method: "GET" | "POST",
    endpoint: string,
    data?: any,
    needsAuth: boolean = false
  ): Promise<ApiResponse<T>> {
    try {
      const config: any = {
        method,
        url: `${this.baseUrl}${endpoint}`,
        headers: {},
      };

      if (needsAuth && this.authToken) {
        config.headers.Authorization = `Bearer ${this.authToken}`;
      }

      if (data) {
        config.data = data;
        config.headers["Content-Type"] = "application/json";
      }

      const response: AxiosResponse = await axios(config);
      return {
        success: true,
        status: response.status,
        data: response.data,
      };
    } catch (error) {
      const axiosError = error as AxiosError;
      return {
        success: false,
        status: axiosError.response?.status || 0,
        data: (axiosError.response?.data || axiosError.message) as T,
      };
    }
  }

  // Health check
  async checkHealth(): Promise<ApiResponse> {
    return this.makeRequest("GET", "/health");
  }

  // User management
  async createUser(userData: UserData): Promise<ApiResponse> {
    try {
      const config: any = {
        method: "POST",
        url: `${this.baseUrl}/users`,
        data: userData,
        headers: {
          "Content-Type": "application/json",
        },
      };

      const response: AxiosResponse = await axios(config);

      // レスポンスボディからtokenを取得して設定
      if (response.data && response.data.token) {
        this.setAuthToken(response.data.token);
      }

      return {
        success: true,
        status: response.status,
        data: response.data,
      };
    } catch (error) {
      const axiosError = error as AxiosError;
      return {
        success: false,
        status: axiosError.response?.status || 0,
        data: (axiosError.response?.data || axiosError.message) as any,
      };
    }
  }

  // Rooms
  async getRooms(): Promise<ApiResponse> {
    return this.makeRequest("GET", "/rooms", undefined, true);
  }

  // Room actions
  async performRoomAction(roomId: number, action: RoomAction): Promise<ApiResponse> {
    return this.makeRequest("POST", `/rooms/${roomId}/actions`, action, true);
  }

  // Formula submission
  async submitFormula(roomId: number, formula: FormulaSubmission): Promise<ApiResponse> {
    return this.makeRequest("POST", `/rooms/${roomId}/formulas`, formula, true);
  }

  // Room results
  async getRoomResults(roomId: number): Promise<ApiResponse> {
    return this.makeRequest("GET", `/rooms/${roomId}/result`, undefined, true);
  }
}

// Default export instance
export const apiClient = new ApiClient();
