const COOKIE_AUTH_TOKEN = 'auth_token';
const COOKIE_AUTH_USERNAME = 'auth_username';
const COOKIE_AUTH_USER_ID = 'auth_user_id';

const LS_AUTH_TOKEN = 'auth_token';
const LS_AUTH_USERNAME = 'auth_username';
const LS_AUTH_USER_ID = 'auth_user_id';

export function setAuthTokenCookie(token: string, username: string, userId: string) {
  localStorage.setItem(LS_AUTH_TOKEN, token);
  localStorage.setItem(LS_AUTH_USERNAME, username);
  localStorage.setItem(LS_AUTH_USER_ID, userId);
  
  // Set cookies with proper attributes for server-side access
  // SameSite=Lax allows the cookie to be sent with requests
  // path=/ makes it available across the entire site
  document.cookie = `${COOKIE_AUTH_TOKEN}=${token}; path=/; max-age=31536000; SameSite=Lax`;
  document.cookie = `${COOKIE_AUTH_USERNAME}=${username}; path=/; max-age=31536000; SameSite=Lax`;
  document.cookie = `${COOKIE_AUTH_USER_ID}=${userId}; path=/; max-age=31536000; SameSite=Lax`;
}

export function removeAuthTokenCookie() {
  localStorage.removeItem(LS_AUTH_TOKEN);
  localStorage.removeItem(LS_AUTH_USERNAME);
  localStorage.removeItem(LS_AUTH_USER_ID);
  
  document.cookie = `${COOKIE_AUTH_TOKEN}=; path=/; max-age=0`;
  document.cookie = `${COOKIE_AUTH_USERNAME}=; path=/; max-age=0`;
  document.cookie = `${COOKIE_AUTH_USER_ID}=; path=/; max-age=0`;
}

export async function getAuthToken(): Promise<string> {
  // Client-side
  if (typeof window !== 'undefined') {
    const token = localStorage.getItem(LS_AUTH_TOKEN);
    return token || '';
  }
  
  // Server-side
  try {
    const { cookies } = await import('next/headers');
    const cookieStore = await cookies();
    return cookieStore.get(COOKIE_AUTH_TOKEN)?.value || '';
  } catch (error) {
    // If cookies() fails (e.g., not in a server context), return empty
    // This can happen in middleware or edge runtime
    return '';
  }
}

export async function getAuthUsername(): Promise<string> {
  // Client-side
  if (typeof window !== 'undefined') {
    return localStorage.getItem(LS_AUTH_USERNAME) || '';
  }

  // Server-side
  try {
    const { cookies } = await import('next/headers');
    const cookieStore = await cookies();
    return cookieStore.get(COOKIE_AUTH_USERNAME)?.value || '';
  } catch {
    return '';
  }
}

export async function getAuthUserId(): Promise<string> {
  // Client-side
  if (typeof window !== 'undefined') {
    return localStorage.getItem(LS_AUTH_USER_ID) || '';
  }

  // Server-side
  try {
    const { cookies } = await import('next/headers');
    const cookieStore = await cookies();
    return cookieStore.get(COOKIE_AUTH_USER_ID)?.value || '';
  } catch {
    return '';
  }
}

// Client-side hooks for React components
// These read synchronously from localStorage (client-side only)
export function useAuthToken(): string {
  if (typeof window === 'undefined') return '';
  return localStorage.getItem(LS_AUTH_TOKEN) || '';
}

export function useAuthUsername(): string | null {
  if (typeof window === 'undefined') return null;
  return localStorage.getItem(LS_AUTH_USERNAME) || null;
}

export function useAuthUserId(): string | null {
  if (typeof window === 'undefined') return null;
  return localStorage.getItem(LS_AUTH_USER_ID) || null;
}

