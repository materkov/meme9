/**
 * Auto-generated API clients using twirp-ts
 */

import {
  AuthClientJSON as AuthClientJSONClass,
} from '@/schema/auth.twirp';
import {
  FeedClientJSON as FeedClientJSONClass,
} from '@/schema/feed.twirp';
import {
  PostsClientJSON as PostsClientJSONClass,
} from '@/schema/posts.twirp';
import {
  UsersClientJSON as UsersClientJSONClass,
} from '@/schema/users.twirp';
import {
  SubscriptionsClientJSON as SubscriptionsClientJSONClass,
} from '@/schema/subscriptions.twirp';
// RPC implementation for twirp-ts generated clients
interface TwirpRpc {
  request(
    service: string,
    method: string,
    contentType: "application/json" | "application/protobuf",
    data: object | Uint8Array
  ): Promise<object | Uint8Array>;
}

const API_BASE_URL = 'http://localhost:8080';
//const API_BASE_URL = 'https://meme2.mmaks.me';

// Auth token storage key
const AUTH_TOKEN_KEY = 'auth_token';

// Helper to get auth token - works on both client and server
async function getAuthTokenInternal(): Promise<string | null> {
  // Client-side: read from localStorage
  if (typeof window !== 'undefined') {
    return localStorage.getItem(AUTH_TOKEN_KEY);
  }
  
  // Server-side: read from cookies
  try {
    const { cookies } = await import('next/headers');
    const cookieStore = await cookies();
    return cookieStore.get(AUTH_TOKEN_KEY)?.value || null;
  } catch {
    // If cookies() is not available (e.g., in client component), return null
    return null;
  }
}

/**
 * RPC implementation that twirp-ts generated clients expect
 */
class TwirpRpcImpl implements TwirpRpc {
  private baseURL: string;
  private headers: Record<string, string>;
  private getToken?: () => Promise<string | null> | string | null;

  constructor(authToken?: string | null | (() => Promise<string | null> | string | null)) {
    this.baseURL = API_BASE_URL;
    this.headers = {
      'Content-Type': 'application/json',
    };
    
    // If authToken is a function, store it for async token retrieval
    if (typeof authToken === 'function') {
      this.getToken = authToken;
    } else if (authToken !== undefined) {
      // Use provided token directly
      if (authToken) {
        this.headers['Authorization'] = authToken;
      }
    } else {
      // No token provided - use getAuthTokenInternal which works on both client and server
      this.getToken = getAuthTokenInternal;
    }
  }

  private async getAuthHeader(): Promise<string | null> {
    if (this.getToken) {
      const token = await this.getToken();
      return token || null;
    }
    return this.headers['Authorization'] || null;
  }

  async request(
    service: string,
    method: string,
    contentType: "application/json" | "application/protobuf",
    data: object | Uint8Array
  ): Promise<object | Uint8Array> {
    const url = `${this.baseURL}/twirp/${service}/${method}`;
    
    // Get token if we have a token getter function
    const headers: Record<string, string> = { ...this.headers } as Record<string, string>;
    if (this.getToken) {
      const token = await this.getAuthHeader();
      if (token) {
        headers['Authorization'] = token;
      }
    }
    
    const body = contentType === 'application/json' 
      ? JSON.stringify(data)
      : data;

    const response = await fetch(url, {
      method: 'POST',
      headers,
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

export { TwirpRpcImpl };

// Auth Service Client - pre-configured instance
export const AuthClient = new AuthClientJSONClass(new TwirpRpcImpl());

// Helper to get auth token from localStorage (works in browser only)
export function getAuthToken(): string | null {
  if (typeof window === 'undefined') return null;
  return localStorage.getItem(AUTH_TOKEN_KEY);
}

// Helper to set auth token
export function setAuthToken(token: string): void {
  if (typeof window === 'undefined') return;
  localStorage.setItem(AUTH_TOKEN_KEY, token);
}

// Helper to remove auth token
export function removeAuthToken(): void {
  if (typeof window === 'undefined') return;
  localStorage.removeItem(AUTH_TOKEN_KEY);
  localStorage.removeItem('auth_username');
  localStorage.removeItem('auth_user_id');
}

// Feed Service Client - pre-configured instance
export const FeedClient = new FeedClientJSONClass(new TwirpRpcImpl());

// Posts Service Client - pre-configured instance
export const PostsClient = new PostsClientJSONClass(new TwirpRpcImpl());

// Users Service Client - pre-configured instance
export const UsersClient = new UsersClientJSONClass(new TwirpRpcImpl());

// Subscriptions Service Client - pre-configured instance
export const SubscriptionsClient = new SubscriptionsClientJSONClass(new TwirpRpcImpl());

