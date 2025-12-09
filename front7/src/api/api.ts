//const API_BASE_URL = 'http://localhost:8080';
//const API_BASE_URL = 'https://meme.mmaks.me';
const API_BASE_URL = window.API_BASE_URL;

export interface Post {
  id: string;
  text: string;
  user_id: string;
  username: string;
  createdAt: string;
}

export interface ErrorResponse {
  error: string;
  error_details: string;
}

export class ApiError extends Error {
  errorCode: string;
  errorDetails: string;

  constructor(errorCode: string, errorDetails: string) {
    super(`${errorCode}: ${errorDetails}`);
    this.name = 'ApiError';
    this.errorCode = errorCode;
    this.errorDetails = errorDetails;
  }
}

async function handleErrorResponse(response: Response): Promise<ApiError> {
  const errorData: ErrorResponse = await response.json();
  return new ApiError(errorData.error, errorData.error_details);
}

export type FeedType = 'global' | 'subscriptions';

export interface FeedRequest {
  type: FeedType;
}

export async function fetchPosts(feedType: FeedType = 'global'): Promise<Post[]> {
  const response = await fetch(`${API_BASE_URL}/feed`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify({ type: feedType }),
  });
  if (!response.ok) {
    throw await handleErrorResponse(response);
  }
  return response.json();
}

export interface FetchUserPostsRequest {
  user_id: string;
}

export async function fetchUserPosts(userID: string): Promise<Post[]> {
  const response = await fetch(`${API_BASE_URL}/userPosts`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify({ user_id: userID }),
  });
  if (!response.ok) {
    throw await handleErrorResponse(response);
  }
  return response.json();
}

export interface PublishPostRequest {
  text: string;
}

export interface PublishPostResponse {
  id: string;
}

function getAuthToken(): string | null {
  return localStorage.getItem('auth_token');
}

function getHeaders(): Record<string, string> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };
  const token = getAuthToken();
  
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }
  
  return headers;
}

export async function publishPost(data: PublishPostRequest): Promise<PublishPostResponse> {
  const response = await fetch(`${API_BASE_URL}/publish`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify(data),
  });

  if (!response.ok) {
    throw await handleErrorResponse(response);
  }

  return response.json();
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user_id: string;
  username: string;
}

export async function login(data: LoginRequest): Promise<LoginResponse> {
  const response = await fetch(`${API_BASE_URL}/login`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify(data),
  });

  if (!response.ok) {
    throw await handleErrorResponse(response);
  }

  return response.json();
}

export interface RegisterRequest {
  username: string;
  password: string;
}

export async function register(data: RegisterRequest): Promise<LoginResponse> {
  const response = await fetch(`${API_BASE_URL}/register`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify(data),
  });

  if (!response.ok) {
    throw await handleErrorResponse(response);
  }

  return response.json();
}

export interface SubscribeRequest {
  user_id: string;
}

export interface SubscribeResponse {
  subscribed: boolean;
}

export async function subscribe(userID: string): Promise<SubscribeResponse> {
  const response = await fetch(`${API_BASE_URL}/subscribe`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify({ user_id: userID }),
  });

  if (!response.ok) {
    throw await handleErrorResponse(response);
  }

  return response.json();
}

export async function unsubscribe(userID: string): Promise<SubscribeResponse> {
  const response = await fetch(`${API_BASE_URL}/unsubscribe`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify({ user_id: userID }),
  });

  if (!response.ok) {
    throw await handleErrorResponse(response);
  }

  return response.json();
}

export async function getSubscriptionStatus(userID: string): Promise<SubscribeResponse> {
  const response = await fetch(`${API_BASE_URL}/subscriptionStatus`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify({ user_id: userID }),
  });

  if (!response.ok) {
    throw await handleErrorResponse(response);
  }

  return response.json();
}

