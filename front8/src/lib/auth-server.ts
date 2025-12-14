import { cookies } from 'next/headers';
import { AuthClient, UsersClient } from './api-clients';
import type { VerifyTokenResponse } from '@/schema/auth';
import type { GetUserResponse as User } from '@/schema/users';

// Server-side auth token key (stored in cookies)
const AUTH_TOKEN_COOKIE = 'auth_token';

/**
 * Get auth token from cookies (server-side)
 */
export async function getServerAuthToken(): Promise<string | null> {
  const cookieStore = await cookies();
  return cookieStore.get(AUTH_TOKEN_COOKIE)?.value || null;
}

/**
 * Verify token on server side and return user ID
 */
export async function verifyServerToken(token: string): Promise<VerifyTokenResponse | null> {
  try {
    const response = await AuthClient.VerifyToken({ token });
    return response;
  } catch (error) {
    return null;
  }
}

/**
 * Get authenticated user ID from server-side token
 * Returns null if token is invalid or missing
 */
export async function getServerUserId(): Promise<string | null> {
  const token = await getServerAuthToken();
  if (!token) {
    return null;
  }

  const verifyResponse = await verifyServerToken(token);
  return verifyResponse?.userId || null;
}

/**
 * Get authenticated user info from server-side token
 * Returns null if token is invalid or missing
 */
export async function getServerUser(): Promise<User | null> {
  const userId = await getServerUserId();
  if (!userId) {
    return null;
  }

  try {
    // Standard client automatically reads token from cookies on server
    const user = await UsersClient.Get({ userId });
    return user;
  } catch (error) {
    return null;
  }
}

/**
 * Check if user is authenticated on server side
 */
export async function isServerAuthenticated(): Promise<boolean> {
  const userId = await getServerUserId();
  return userId !== null;
}
