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

export async function fetchPosts(): Promise<Post[]> {
  const response = await fetch(`${API_BASE_URL}/feed`);
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

export async function publishPost(data: PublishPostRequest): Promise<PublishPostResponse> {
  const token = getAuthToken();
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_BASE_URL}/publish`, {
    method: 'POST',
    headers,
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
    headers: {
      'Content-Type': 'application/json',
    },
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
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });

  if (!response.ok) {
    throw await handleErrorResponse(response);
  }

  return response.json();
}

