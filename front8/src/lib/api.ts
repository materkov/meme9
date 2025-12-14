// Import proto-generated types
import type {
  FeedRequest,
  FeedResponse,
  FeedPostResponse,
} from '@/schema/feed';
import type {
  GetPostRequest,
  GetPostResponse,
  GetByUsersRequest,
  GetByUsersResponse,
  UserPostResponse,
  PublishRequest,
  PublishResponse,
} from '@/schema/posts';
import type {
  GetUserRequest,
  GetUserResponse,
} from '@/schema/users';
import type {
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  VerifyTokenRequest,
  VerifyTokenResponse,
} from '@/schema/auth';
import type {
  SubscribeRequest,
  SubscribeResponse,
  GetFollowingRequest,
  GetFollowingResponse,
  IsSubscribedRequest,
  IsSubscribedResponse,
} from '@/schema/subscriptions';

const API_BASE_URL = 'http://localhost:8080';
//const API_BASE_URL = 'https://meme2.mmaks.me';

// Auth token storage key
const AUTH_TOKEN_KEY = 'auth_token';

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

// Helper to get headers with auth if available (client-side)
export function getHeaders(): HeadersInit {
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
  };
  const token = getAuthToken();
  if (token) {
    headers['Authorization'] = token;
  }
  return headers;
}

// Helper to get headers with auth token (server-side)
export function getServerHeaders(token: string | null): HeadersInit {
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
  };
  if (token) {
    headers['Authorization'] = token;
  }
  return headers;
}


// Feed API
export async function getFeed(feedType: string = 'all', authToken?: string | null): Promise<FeedPostResponse[]> {
  const headers = authToken !== undefined 
    ? getServerHeaders(authToken)
    : getHeaders();
  
  const requestBody: FeedRequest = { type: feedType };
  
  const response = await fetch(`${API_BASE_URL}/twirp/meme.feed.Feed/GetFeed`, {
    method: 'POST',
    headers,
    body: JSON.stringify(requestBody),
    cache: 'no-store',
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || `Failed to fetch feed: ${response.statusText}`);
  }

  const data: FeedResponse = await response.json();
  return data.posts || [];
}

// Posts API
export async function getPost(postId: string, authToken?: string | null): Promise<GetPostResponse> {
  const headers = authToken !== undefined 
    ? getServerHeaders(authToken)
    : getHeaders();
  
  const requestBody: GetPostRequest = { postId };
  
  const response = await fetch(`${API_BASE_URL}/twirp/meme.posts.Posts/Get`, {
    method: 'POST',
    headers,
    body: JSON.stringify(requestBody),
    cache: 'no-store',
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || `Failed to fetch post: ${response.statusText}`);
  }

  const data: GetPostResponse = await response.json();
  return data;
}

export async function getUser(userId: string, authToken?: string | null): Promise<GetUserResponse> {
  const headers = authToken !== undefined 
    ? getServerHeaders(authToken)
    : getHeaders();
  
  const requestBody: GetUserRequest = { userId };
  
  const response = await fetch(`${API_BASE_URL}/twirp/meme.users.Users/Get`, {
    method: 'POST',
    headers,
    body: JSON.stringify(requestBody),
    cache: 'no-store',
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || `Failed to fetch user: ${response.statusText}`);
  }

  const data: GetUserResponse = await response.json();
  return data;
}

// Auth API functions
export async function login(request: LoginRequest): Promise<LoginResponse> {
  const response = await fetch(`${API_BASE_URL}/twirp/meme.auth.Auth/Login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(request),
    cache: 'no-store',
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || `Failed to login: ${response.statusText}`);
  }

  const data: LoginResponse = await response.json();
  return data;
}

export async function register(request: RegisterRequest): Promise<LoginResponse> {
  const response = await fetch(`${API_BASE_URL}/twirp/meme.auth.Auth/Register`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(request),
    cache: 'no-store',
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || `Failed to register: ${response.statusText}`);
  }

  const data: LoginResponse = await response.json();
  return data;
}

export async function verifyToken(token: string): Promise<VerifyTokenResponse> {
  const requestBody: VerifyTokenRequest = { token };
  
  const response = await fetch(`${API_BASE_URL}/twirp/meme.auth.Auth/VerifyToken`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(requestBody),
    cache: 'no-store',
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || `Failed to verify token: ${response.statusText}`);
  }

  const data: VerifyTokenResponse = await response.json();
  return data;
}

// Publish API function
export async function publishPost(request: PublishRequest, authToken?: string | null): Promise<PublishResponse> {
  const headers = authToken !== undefined 
    ? getServerHeaders(authToken)
    : getHeaders();
  
  const response = await fetch(`${API_BASE_URL}/twirp/meme.posts.Posts/Publish`, {
    method: 'POST',
    headers,
    body: JSON.stringify(request),
    cache: 'no-store',
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || `Failed to publish post: ${response.statusText}`);
  }

  const data: PublishResponse = await response.json();
  return data;
}

// UserPosts API function
export async function getUserPosts(userId: string, authToken?: string | null): Promise<UserPostResponse[]> {
  const headers = authToken !== undefined 
    ? getServerHeaders(authToken)
    : getHeaders();
  
  const requestBody: GetByUsersRequest = { userId };
  
  const response = await fetch(`${API_BASE_URL}/twirp/meme.posts.Posts/GetByUsers`, {
    method: 'POST',
    headers,
    body: JSON.stringify(requestBody),
    cache: 'no-store',
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || `Failed to fetch user posts: ${response.statusText}`);
  }

  const data: GetByUsersResponse = await response.json();
  return data.posts || [];
}

// Subscribe API function
export async function subscribe(userId: string): Promise<SubscribeResponse> {
  const requestBody: SubscribeRequest = { userId };
  
  const response = await fetch(`${API_BASE_URL}/twirp/meme.subscriptions.Subscriptions/Subscribe`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify(requestBody),
    cache: 'no-store',
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || `Failed to subscribe: ${response.statusText}`);
  }

  const data: SubscribeResponse = await response.json();
  return data;
}

// Unsubscribe API function
export async function unsubscribe(userId: string): Promise<SubscribeResponse> {
  const requestBody: SubscribeRequest = { userId };
  
  const response = await fetch(`${API_BASE_URL}/twirp/meme.subscriptions.Subscriptions/Unsubscribe`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify(requestBody),
    cache: 'no-store',
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || `Failed to unsubscribe: ${response.statusText}`);
  }

  const data: SubscribeResponse = await response.json();
  return data;
}

// GetSubscriptionStatus API function
export async function getSubscriptionStatus(userId: string, authToken?: string | null): Promise<SubscribeResponse> {
  const headers = authToken !== undefined 
    ? getServerHeaders(authToken)
    : getHeaders();
  
  const requestBody: SubscribeRequest = { userId };
  
  const response = await fetch(`${API_BASE_URL}/twirp/meme.subscriptions.Subscriptions/GetStatus`, {
    method: 'POST',
    headers,
    body: JSON.stringify(requestBody),
    cache: 'no-store',
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || `Failed to get subscription status: ${response.statusText}`);
  }

  const data: SubscribeResponse = await response.json();
  return data;
}
