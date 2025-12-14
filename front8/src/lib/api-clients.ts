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

// Helper to get auth token from localStorage (works in browser only)
function getAuthTokenInternal(): string | null {
  if (typeof window === 'undefined') return null;
  return localStorage.getItem(AUTH_TOKEN_KEY);
}

/**
 * RPC implementation that twirp-ts generated clients expect
 */
class TwirpRpcImpl implements TwirpRpc {
  private baseURL: string;
  private headers: HeadersInit;

  constructor(authToken?: string | null) {
    this.baseURL = API_BASE_URL;
    this.headers = {
      'Content-Type': 'application/json',
    };
    const token = authToken !== undefined ? authToken : (typeof window !== 'undefined' ? getAuthTokenInternal() : null);
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
}

// Feed Service Client - pre-configured instance
export const FeedClient = new FeedClientJSONClass(new TwirpRpcImpl());

// Posts Service Client - pre-configured instance
export const PostsClient = new PostsClientJSONClass(new TwirpRpcImpl());

// Users Service Client - pre-configured instance
export const UsersClient = new UsersClientJSONClass(new TwirpRpcImpl());

// Subscriptions Service Client - pre-configured instance
export const SubscriptionsClient = new SubscriptionsClientJSONClass(new TwirpRpcImpl());

// Helper function to create FeedClient with custom token (for server-side)
export function createFeedClientWithToken(token: string | null): FeedClientJSONClass {
  return new FeedClientJSONClass(new TwirpRpcImpl(token));
}
