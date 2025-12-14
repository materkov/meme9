const COOKIE_AUTH_TOKEN = 'auth_token';

const LS_AUTH_TOKEN = 'auth_token';
const LS_AUTH_USERNAME = 'auth_username';
const LS_AUTH_USER_ID = 'auth_user_id';

export function setAuthTokenCookie(token: string, username: string, userId: string) {
  localStorage.setItem(LS_AUTH_TOKEN, token);
  localStorage.setItem(LS_AUTH_USERNAME, username);
  localStorage.setItem(LS_AUTH_USER_ID, userId);
  
  document.cookie = `${COOKIE_AUTH_TOKEN}=${token}; path=/; max-age=31536000; SameSite=Lax`;
}

export function removeAuthTokenCookie() {
  localStorage.removeItem(LS_AUTH_TOKEN);
  localStorage.removeItem(LS_AUTH_USERNAME);
  localStorage.removeItem(LS_AUTH_USER_ID);
  
  document.cookie = `${COOKIE_AUTH_TOKEN}=; path=/; max-age=0`;
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
  } catch {
    // If cookies() fails (e.g., not in a server context), return empty
    return '';
  }
}

export function getAuthUsername(): string {
  // Client-side
  if (typeof window !== 'undefined') {
    return localStorage.getItem(LS_AUTH_USERNAME) || '';
  }

  // Server-side
  return ''
}

export function getAuthUserId(): string {
  // Client-side
  if (typeof window !== 'undefined') {
    return localStorage.getItem(LS_AUTH_USER_ID) || '';
  }

  // Server-side
  return ''
}
