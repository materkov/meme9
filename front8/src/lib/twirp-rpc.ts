/**
 * RPC implementation for twirp-ts generated clients
 * Works with JSON format (Twirp's default)
 */

const API_BASE_URL = 'http://localhost:8080';
//const API_BASE_URL = 'https://meme2.mmaks.me';

// Auth token storage key
const AUTH_TOKEN_KEY = 'auth_token';

// Helper to get auth token from localStorage (works in browser only)
function getAuthToken(): string | null {
  if (typeof window === 'undefined') return null;
  return localStorage.getItem(AUTH_TOKEN_KEY);
}

/**
 * RPC interface that twirp-ts generated clients expect
 */
export interface TwirpRpc {
  request(
    service: string,
    method: string,
    contentType: "application/json" | "application/protobuf",
    data: object | Uint8Array
  ): Promise<object | Uint8Array>;
}

/**
 * JSON-based RPC implementation for Twirp
 */
export class TwirpRpcImpl implements TwirpRpc {
  private baseURL: string;
  private headers: HeadersInit;

  constructor(authToken?: string | null) {
    this.baseURL = API_BASE_URL;
    this.headers = {
      'Content-Type': 'application/json',
    };
    const token = authToken !== undefined ? authToken : (typeof window !== 'undefined' ? getAuthToken() : null);
    if (token) {
      this.headers['Authorization'] = token;
    }
  }

  async request(
    service: string,
    method: string,
    contentType: "application/json" | "application/protobuf",
    data: object | Uint8Array
  ): Promise<object | Uint8Array> {
    const url = `${this.baseURL}/twirp/${service}/${method}`;
    
    const body = contentType === 'application/json' 
      ? JSON.stringify(data)
      : data;

    const response = await fetch(url, {
      method: 'POST',
      headers: this.headers,
      body: body as BodyInit,
      cache: 'no-store',
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ 
        code: 'unknown',
        msg: `HTTP ${response.status}: ${response.statusText}` 
      }));
      throw new Error(errorData.msg || errorData.error || `Request failed: ${response.statusText}`);
    }

    if (contentType === 'application/json') {
      return response.json();
    } else {
      const arrayBuffer = await response.arrayBuffer();
      return new Uint8Array(arrayBuffer);
    }
  }
}
