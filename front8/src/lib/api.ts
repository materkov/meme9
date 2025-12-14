const API_BASE_URL = 'http://localhost:8080';
//const API_BASE_URL = 'https://meme.mmaks.me';

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

export interface FeedPost {
  id: string;
  text: string;
  user_id: string;
  username: string;
  created_at: string;
}

export interface FeedRequest {
  type?: string; // "all" or "subscriptions"
}

export interface FeedResponse {
  posts: FeedPost[];
}

export async function getFeed(feedType: string = 'all', authToken?: string | null): Promise<FeedPost[]> {
  const headers = authToken !== undefined 
    ? getServerHeaders(authToken)
    : getHeaders();
  
  const response = await fetch(`${API_BASE_URL}/twirp/meme.json_api.JsonAPI/GetFeed`, {
    method: 'POST',
    headers,
    body: JSON.stringify({ type: feedType }),
    // Enable caching for server-side rendering
    cache: 'no-store', // Use 'no-store' for always fresh data, or 'force-cache' for caching
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || `Failed to fetch feed: ${response.statusText}`);
  }

  const data: FeedResponse = await response.json();
  return data.posts || [];
}

export interface Post {
  id: string;
  text: string;
  user_id: string;
  created_at: string;
}

export interface GetPostRequest {
  post_id: string;
}

export interface GetPostResponse {
  id: string;
  text: string;
  user_id: string;
  created_at: string;
}

export async function getPost(postId: string, authToken?: string | null): Promise<Post> {
  const headers = authToken !== undefined 
    ? getServerHeaders(authToken)
    : getHeaders();
  
  const response = await fetch(`${API_BASE_URL}/twirp/meme.json_api.JsonAPI/GetPost`, {
    method: 'POST',
    headers,
    body: JSON.stringify({ post_id: postId }),
    cache: 'no-store',
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || `Failed to fetch post: ${response.statusText}`);
  }

  const data: GetPostResponse = await response.json();
  return data;
}

export interface User {
  id: string;
  username: string;
}

export interface GetUserRequest {
  user_id: string;
}

export interface GetUserResponse {
  id: string;
  username: string;
}

export async function getUser(userId: string, authToken?: string | null): Promise<User> {
  const headers = authToken !== undefined 
    ? getServerHeaders(authToken)
    : getHeaders();
  
  const response = await fetch(`${API_BASE_URL}/twirp/meme.json_api.JsonAPI/GetUser`, {
    method: 'POST',
    headers,
    body: JSON.stringify({ user_id: userId }),
    cache: 'no-store',
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || `Failed to fetch user: ${response.statusText}`);
  }

  const data: GetUserResponse = await response.json();
  return data;
}

// Auth interfaces
export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user_id: string;
  username: string;
}

export interface RegisterRequest {
  username: string;
  password: string;
}

export interface VerifyTokenRequest {
  token: string;
}

export interface VerifyTokenResponse {
  user_id: string;
}

// Auth API functions
export async function login(request: LoginRequest): Promise<LoginResponse> {
  const response = await fetch(`${API_BASE_URL}/twirp/meme.json_api.JsonAPI/Login`, {
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
  const response = await fetch(`${API_BASE_URL}/twirp/meme.json_api.JsonAPI/Register`, {
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
  const response = await fetch(`${API_BASE_URL}/twirp/meme.json_api.JsonAPI/VerifyToken`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ token }),
    cache: 'no-store',
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || `Failed to verify token: ${response.statusText}`);
  }

  const data: VerifyTokenResponse = await response.json();
  return data;
}

// Publish interfaces
export interface PublishRequest {
  text: string;
}

export interface PublishResponse {
  id: string;
}

// Publish API function
export async function publishPost(request: PublishRequest, authToken?: string | null): Promise<PublishResponse> {
  const headers = authToken !== undefined 
    ? getServerHeaders(authToken)
    : getHeaders();
  
  const response = await fetch(`${API_BASE_URL}/twirp/meme.json_api.JsonAPI/Publish`, {
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

// UserPosts interfaces
export interface UserPost {
  id: string;
  text: string;
  user_id: string;
  username: string;
  created_at: string;
}

export interface UserPostsRequest {
  user_id: string;
}

export interface UserPostsResponse {
  posts: UserPost[];
}

// UserPosts API function
export async function getUserPosts(userId: string, authToken?: string | null): Promise<UserPost[]> {
  const headers = authToken !== undefined 
    ? getServerHeaders(authToken)
    : getHeaders();
  
  const response = await fetch(`${API_BASE_URL}/twirp/meme.json_api.JsonAPI/GetUserPosts`, {
    method: 'POST',
    headers,
    body: JSON.stringify({ user_id: userId }),
    cache: 'no-store',
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || `Failed to fetch user posts: ${response.statusText}`);
  }

  const data: UserPostsResponse = await response.json();
  return data.posts || [];
}

// Subscribe interfaces
export interface SubscribeRequest {
  user_id: string;
}

export interface SubscribeResponse {
  subscribed: boolean;
}

// Subscribe API function
export async function subscribe(userId: string): Promise<SubscribeResponse> {
  const response = await fetch(`${API_BASE_URL}/twirp/meme.json_api.JsonAPI/Subscribe`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify({ user_id: userId }),
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
  const response = await fetch(`${API_BASE_URL}/twirp/meme.json_api.JsonAPI/Unsubscribe`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify({ user_id: userId }),
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
  
  const response = await fetch(`${API_BASE_URL}/twirp/meme.json_api.JsonAPI/GetSubscriptionStatus`, {
    method: 'POST',
    headers,
    body: JSON.stringify({ user_id: userId }),
    cache: 'no-store',
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || `Failed to get subscription status: ${response.statusText}`);
  }

  const data: SubscribeResponse = await response.json();
  return data;
}
